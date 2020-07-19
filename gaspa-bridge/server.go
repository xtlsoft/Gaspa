package main

import (
	"bufio"
	"log"
	"net"
	"time"

	"github.com/satori/uuid"
)

type server struct {
	connections map[uuid.UUID]connection
	uuidMap     map[string]uuid.UUID
}

type connection struct {
	name string
	uuid uuid.UUID
	conn *net.TCPConn
}

func newServer() *server {
	srv := new(server)
	srv.connections = make(map[uuid.UUID]connection)
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
	writer := bufio.NewWriter(conn)
	bio := bufio.NewReadWriter(reader, writer)
	_ = bio
	header, err := reader.ReadString('!')
	if err != nil {
		log.Println(err)
	}

}
