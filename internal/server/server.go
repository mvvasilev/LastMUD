package server

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	listener *net.TCPListener

	connections []*Connection

	stopChannel chan struct{}
}

func CreateServer(port string) (srv *Server, err error) {
	addr, err := net.ResolveTCPAddr("tcp", port)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ln, err := net.ListenTCP("tcp", addr)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Listening on port", port)

	srv = &Server{
		listener:    ln,
		connections: []*Connection{},
		stopChannel: make(chan struct{}),
	}

	return
}

func (srv *Server) Listen() {
	// Wait for 200 millis for a new connection, and loop if none found in that time
	srv.listener.SetDeadline(time.Now().Add(1 * time.Second))

	for {
		if srv.shouldStop() {
			break
		}

		conn, err := srv.listener.Accept()

		if err != nil {
			continue
		}

		c := CreateConnection(conn)
		srv.connections = append(srv.connections, c)

		go c.listen()
	}

	for _, v := range srv.connections {
		v.Close()
	}

	err := srv.listener.Close()

	if err != nil {
		fmt.Println(err)
	}
}

func (srv *Server) shouldStop() bool {
	select {
	case <-srv.stopChannel:
		return true
	default:
		return false
	}
}

func (srv *Server) Stop() {
	srv.stopChannel <- struct{}{}
}
