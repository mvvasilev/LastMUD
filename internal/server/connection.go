package server

import (
	"bufio"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
	"code.haedhutner.dev/mvv/LastMUD/internal/term"
	"context"
	"github.com/google/uuid"
	"net"
	"sync"
	"time"
)

const CheckAlivePeriod = 50 * time.Millisecond

type Connection struct {
	ctx    context.Context
	wg     *sync.WaitGroup
	server *Server

	identity uuid.UUID

	term     *term.VirtualTerm
	conn     *net.TCPConn
	lastSeen time.Time

	stop context.CancelFunc
}

func CreateConnection(server *Server, conn *net.TCPConn, ctx context.Context, wg *sync.WaitGroup) (c *Connection, err error) {
	ctx, cancel := context.WithCancel(ctx)

	t, err := term.CreateVirtualTerm(
		ctx,
		wg,
		func(t time.Time) {
			_ = conn.SetReadDeadline(t)
		},
		bufio.NewReader(conn),
		bufio.NewWriter(conn),
	)

	if err != nil {
		cancel()
		return nil, err
	}

	c = &Connection{
		ctx:      ctx,
		wg:       wg,
		server:   server,
		identity: uuid.New(),
		term:     t,
		conn:     conn,
		lastSeen: time.Now(),
		stop:     cancel,
	}

	logging.Info("Connection from ", c.conn.RemoteAddr(), ": Assigned id ", c.Id().String())

	wg.Add(1)
	go c.checkAliveAndConsumeCommands()

	server.game().Connect(c.Id())

	return
}

func (c *Connection) Id() uuid.UUID {
	return c.identity
}

func (c *Connection) Write(output []byte) (err error) {
	if c.shouldClose() {
		return nil
	}

	err = c.term.Write(output)
	return
}

func (c *Connection) Close() {
	c.stop()
}

func (c *Connection) shouldClose() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
	}

	return false
}

func (c *Connection) checkAliveAndConsumeCommands() {
	defer c.wg.Done()
	defer c.closeConnection()

	for {
		if c.shouldClose() {
			break
		}

		_, err := c.conn.Write([]byte{0x00})

		if err != nil {
			break
		}

		cmd := c.term.NextCommand()

		if cmd != "" {
			c.server.game().SendCommand(c.Id(), cmd)
		}

		time.Sleep(CheckAlivePeriod)
	}
}

func (c *Connection) closeConnection() {
	c.term.Close()
	c.conn.Close()

	c.server.game().Disconnect(c.Id())

	logging.Info("Disconnect ", c.conn.RemoteAddr(), " with id ", c.Id().String())
}
