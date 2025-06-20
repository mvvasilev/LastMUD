package server

import (
	"context"
	"net"
	"sync"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/game"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
)

type Server struct {
	ctx context.Context
	wg  *sync.WaitGroup

	listener *net.TCPListener

	connections []*Connection

	game *game.LastMUDGame
}

func CreateServer(ctx context.Context, wg *sync.WaitGroup, port string) (srv *Server, err error) {

	logging.Info("	 _              _   __  __ _   _ ____     ")
	logging.Info("	| |    __ _ ___| |_|  \\/  | | | |  _ \\    ")
	logging.Info("	| |   / _` / __| __| |\\/| | | | | | | |   ")
	logging.Info("	| |__| (_| \\__ \\ |_| |  | | |_| | |_| |   ")
	logging.Info("	|_____\\__,_|___/\\__|_|  |_|\\___/|____/    ")
	logging.Info("										      ")

	addr, err := net.ResolveTCPAddr("tcp", port)

	if err != nil {
		logging.Error(err)
		return nil, err
	}

	ln, err := net.ListenTCP("tcp", addr)

	if err != nil {
		logging.Error(err)
		return nil, err
	}

	logging.Info("Starting server, listening on port ", port)

	srv = &Server{
		ctx:         ctx,
		wg:          wg,
		listener:    ln,
		connections: []*Connection{},
	}

	srv.game = game.CreateGame(ctx, srv.wg)

	srv.wg.Add(1)
	go srv.listen()

	return
}

func (srv *Server) listen() {
	defer srv.wg.Done()
	defer srv.shutdown()

	for {
		srv.listener.SetDeadline(time.Now().Add(1 * time.Second))

		if srv.shouldStop() {
			break
		}

		conn, err := srv.listener.Accept()

		if err != nil {
			continue
		}

		tcpConn, ok := conn.(*net.TCPConn)

		if !ok {
			logging.Warn("Attempted invalid connection from", conn.RemoteAddr())
			continue
		}

		c := CreateConnection(tcpConn, srv.ctx, srv.wg)
		srv.connections = append(srv.connections, c)
	}
}

func (srv *Server) shutdown() {
	logging.Info("Stopping Server...")

	err := srv.listener.Close()

	if err != nil {
		logging.Error(err)
	}
}

func (srv *Server) shouldStop() bool {
	select {
	case <-srv.ctx.Done():
		return true
	default:
		return false
	}
}
