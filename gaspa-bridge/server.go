package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/satori/uuid"
	gaspa "github.com/xtlsoft/Gaspa"
	"github.com/xtlsoft/Gaspa/gaspa-bridge/protocol"
)

type server struct {
	connections map[uuid.UUID]protocol.RegisteredConnection
	uuidMap     map[string]uuid.UUID
}

func newServer() *server {
	srv := new(server)
	srv.connections = make(map[uuid.UUID]protocol.RegisteredConnection)
	srv.uuidMap = make(map[string]uuid.UUID)
	return srv
}

func (srv *server) serve() error {
	addr, err := net.ResolveTCPAddr("tcp", config.listen)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("gaspa-bridge version %s\r\n", gaspa.Version)
	fmt.Printf("Serving on %s\r\n", config.listen)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		if config.keepAlive != 0 {
			conn.SetKeepAlive(true)
			conn.SetKeepAlivePeriod(
				time.Duration(config.keepAlive) * time.Second,
			)
		}
		srv.handleConnection(conn)
	}
}

func (srv *server) handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	method, err := reader.ReadByte()
	if err != nil {
		log.Println(err)
		return
	}
	switch method {
	case protocol.DispatcherServiceConnect:
		// a
	case protocol.DispatcherServiceMeta:
		srv.doMeta(conn, reader)
	case protocol.DispatcherServiceMetaQuery:
		// c
	case protocol.DispatcherServiceRegister:
		srv.doRegister(conn, reader)
	case protocol.DispatcherServiceJoin:
		// e
	}
}

func (srv *server) doRegister(conn *net.TCPConn, reader *bufio.Reader) {
	meta, err := reader.ReadBytes('!')
	if err != nil {
		log.Println(err)
		return
	}
	var uuid uuid.UUID
	copy(uuid[:], meta[:16])
	name := string(meta[16:])
	srv.uuidMap[name] = uuid
	srv.connections[uuid] = protocol.RegisteredConnection{
		Name: name,
		UUID: uuid,
		Conn: conn,
	}
	conn.Write([]byte(`A!`))
	for {
		reader.ReadByte()
	}
}

func (srv *server) doMeta(conn *net.TCPConn, reader *bufio.Reader) {
	json, _ := json.Marshal(&protocol.MetaResult{
		Status: "success",
		Result: srv.connections,
	})
	conn.Write(json)
}
