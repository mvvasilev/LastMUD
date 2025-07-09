package server

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game"
	"code.haedhutner.dev/mvv/LastMUD/shared/log"
	"context"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	ctx context.Context
	wg  *sync.WaitGroup

	listener *net.TCPListener

	connections map[uuid.UUID]*Connection

	lastmudgame *game.Game

	stop context.CancelFunc
}

func CreateServer(ctx context.Context, wg *sync.WaitGroup, port string) (srv *Server, err error) {
	ctx, cancel := context.WithCancel(ctx)

	log.Info(" _              _   __  __ _   _ ____")
	log.Info("| |    __ _ ___| |_|  \\/  | | | |  _ \\")
	log.Info("| |   / _` / __| __| |\\/| | | | | | | |")
	log.Info("| |__| (_| \\__ \\ |_| |  | | |_| | |_| |")
	log.Info("|_____\\__,_|___/\\__|_|  |_|\\___/|____/")
	log.Info("")

	addr, err := net.ResolveTCPAddr("tcp", port)

	if err != nil {
		log.Error(err)
		cancel()

		return nil, err
	}

	ln, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Error(err)
		cancel()

		return nil, err
	}

	log.Info("Starting server, listening on port ", port)

	srv = &Server{
		ctx:         ctx,
		wg:          wg,
		listener:    ln,
		connections: map[uuid.UUID]*Connection{},
		stop:        cancel,
	}

	srv.lastmudgame = game.CreateGame(ctx, srv.wg)

	srv.wg.Add(2)
	go srv.listen()
	go srv.consumeGameOutput()

	return
}

func (srv *Server) game() *game.Game {
	return srv.lastmudgame
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
			log.Warn("Attempted invalid connection from", conn.RemoteAddr())

			continue
		}

		c, err := CreateConnection(srv, tcpConn, srv.ctx, srv.wg)

		if err != nil {
			log.Error("Unable to create connection: ", err)
			_ = tcpConn.Close()
		} else {
			srv.connections[c.Id()] = c
		}
	}
}

func (srv *Server) consumeGameOutput() {
	defer srv.wg.Done()

	for {
		if srv.shouldStop() {
			break
		}

		output := srv.lastmudgame.ConsumeNextOutput()

		if output == nil {
			time.Sleep(20 * time.Millisecond)
			continue
		}

		conn, ok := srv.connections[output.Id()]

		if ok && output.Contents() != nil {
			err := conn.Write(output.Contents())

			if err != nil {
				log.Error("Error writing to connection ", output.Id(), ": ", err)
			}
		}

		if ok && output.ShouldCloseConnection() {
			conn.Close()
			delete(srv.connections, output.Id())
		}
	}
}

func (srv *Server) shutdown() {
	log.Info("Stopping Server...")

	err := srv.listener.Close()

	if err != nil {
		log.Error(err)
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
