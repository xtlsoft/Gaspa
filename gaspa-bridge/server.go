package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/satori/uuid"
	gaspa "github.com/xtlsoft/Gaspa"
	"github.com/xtlsoft/Gaspa/gaspa-bridge/protocol"
)

type server struct {
	lock        *sync.Mutex
	connections map[uuid.UUID]protocol.RegisteredConnection
	uuidMap     map[string]uuid.UUID
}

func newServer() *server {
	srv := new(server)
	srv.connections = make(map[uuid.UUID]protocol.RegisteredConnection)
	srv.uuidMap = make(map[string]uuid.UUID)
	srv.lock = &sync.Mutex{}
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
		srv.doMetaQuery(conn, reader)
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
	name := string(meta[16 : len(meta)-1])
	srv.lock.Lock()
	srv.uuidMap[name] = uuid
	srv.connections[uuid] = protocol.RegisteredConnection{
		Name: name,
		UUID: uuid,
		Conn: conn,
	}
	srv.lock.Unlock()
	conn.Write([]byte(`A!`))
	for {
		_, err := reader.ReadByte()
		if err != nil {
			srv.lock.Lock()
			delete(srv.connections, uuid)
			delete(srv.uuidMap, name)
			srv.lock.Unlock()
			return
		}
	}
}

func (srv *server) doMeta(conn *net.TCPConn, reader *bufio.Reader) {
	json, _ := json.Marshal(&protocol.MetaResult{
		Status: "success",
		Result: srv.connections,
	})
	conn.Write(json)
}

func (srv *server) doMetaQuery(conn *net.TCPConn, reader *bufio.Reader) {
	typ, err := reader.ReadByte()
	if err != nil {
		log.Println(err)
		return
	}
	var uuid uuid.UUID
	name, err := reader.ReadBytes('!')
	if err != nil {
		log.Println(err)
		return
	}
	found := true
	switch typ {
	case protocol.MetaQueryTypeName:
		var ok bool
		uuid, ok = srv.uuidMap[string(name[:len(name)-1])]
		if !ok {
			found = false
		}
	case protocol.MetaQueryTypeUUID:
		if len(name) != 17 {
			log.Printf("invalid UUID %s\r\n", string(name[:len(name)-1]))
		}
		copy(uuid[:], name[:16])
	}
	s, ok := srv.connections[uuid]
	if !ok {
		found = false
	}
	if !found {
		conn.Write([]byte{protocol.MetaQueryNotFound})
		return
	}
	conn.Write([]byte{protocol.MetaQueryFound})
	conn.Write(s.UUID.Bytes())
	conn.Write([]byte(s.Name))
	conn.Write([]byte{'!'})
}

func (srv *server) doConnect(conn *net.TCPConn, reader *bufio.Reader) {
	// FIXME: #1 All the idents parsing fails when uuid contains '!'
	ident, err := reader.ReadBytes('!')
	if err != nil {
		log.Println(err)
		return
	}
	_ = ident
}
