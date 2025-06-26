package server

import (
	"bufio"
	"context"
	"net"
	"sync"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
	"github.com/google/uuid"
)

const MaxLastSeenTime = 90 * time.Second

const CheckAlivePeriod = 50 * time.Millisecond

const DeleteBeforeAndMoveToStartOfLine = "\033[1K\r"

type Connection struct {
	ctx context.Context
	wg  *sync.WaitGroup

	server *Server

	identity uuid.UUID

	conn     *net.TCPConn
	lastSeen time.Time
}

func CreateConnection(server *Server, conn *net.TCPConn, ctx context.Context, wg *sync.WaitGroup) (c *Connection) {
	logging.Info("Connect: ", conn.RemoteAddr())

	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(1 * time.Second)

	c = &Connection{
		ctx:      ctx,
		wg:       wg,
		server:   server,
		identity: uuid.New(),
		conn:     conn,
		lastSeen: time.Now(),
	}

	c.wg.Add(2)
	go c.listen()
	go c.checkAlive()

	server.game().ConnectPlayer(c.Id())

	return
}

func (c *Connection) Id() uuid.UUID {
	return c.identity
}

func (c *Connection) Write(output []byte) (err error) {
	output = append([]byte(DeleteBeforeAndMoveToStartOfLine+"< "), output...)
	output = append(output, []byte("\n> ")...)
	_, err = c.conn.Write(output)
	return
}

func (c *Connection) listen() {
	defer c.wg.Done()

	logging.Info("Listening on connection ", c.conn.RemoteAddr())

	for {
		c.conn.SetReadDeadline(time.Time{})

		message, err := bufio.NewReader(c.conn).ReadString('\n')

		if err != nil {
			logging.Warn(err)
			break
		}

		c.server.game().SendPlayerCommand(c.Id(), message)

		c.lastSeen = time.Now()
	}
}

func (c *Connection) checkAlive() {
	defer c.wg.Done()
	defer c.closeConnection()

	for {
		if c.shouldClose() {
			c.Write([]byte("Server shutting down, bye bye!\r\n"))
			break
		}

		if time.Since(c.lastSeen) > MaxLastSeenTime {
			c.Write([]byte("You have been away for too long, bye bye!\r\n"))
			break
		}

		_, err := c.conn.Write([]byte{0x00})

		if err != nil {
			break
		}

		time.Sleep(CheckAlivePeriod)
	}
}

func (c *Connection) shouldClose() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}

func (c *Connection) closeConnection() {
	c.conn.Close()

	c.server.game().DisconnectPlayer(c.Id())

	logging.Info("Disconnected: ", c.conn.RemoteAddr())
}
