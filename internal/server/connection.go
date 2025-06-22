package server

import (
	"bufio"
	"context"
	"net"
	"sync"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/game"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
	"github.com/google/uuid"
)

const MaxLastSeenTime = 120 * time.Second
const MaxEnqueuedInputMessages = 10

type Connection struct {
	ctx context.Context
	wg  *sync.WaitGroup

	server *Server

	identity uuid.UUID

	conn     *net.TCPConn
	lastSeen time.Time

	inputChannel chan []byte
}

func CreateConnection(server *Server, conn *net.TCPConn, ctx context.Context, wg *sync.WaitGroup) (c *Connection) {
	logging.Info("Connect: ", conn.RemoteAddr())

	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(1 * time.Second)

	c = &Connection{
		ctx:          ctx,
		wg:           wg,
		server:       server,
		identity:     uuid.New(),
		conn:         conn,
		inputChannel: make(chan []byte, MaxEnqueuedInputMessages),
		lastSeen:     time.Now(),
	}

	c.wg.Add(2)
	go c.listen()
	go c.checkAlive()

	server.game.EnqueueEvent(game.CreatePlayerJoinEvent(c.Id()))

	return
}

func (c *Connection) Id() uuid.UUID {
	return c.identity
}

func (c *Connection) listen() {
	defer c.wg.Done()

	logging.Info("Listening on connection ", c.conn.RemoteAddr())

	for {
		c.conn.SetReadDeadline(time.Time{})

		message, err := bufio.NewReader(c.conn).ReadBytes('\n')

		if err != nil {
			logging.Warn(err)
			break
		}

		if len(c.inputChannel) == MaxEnqueuedInputMessages {
			c.conn.Write([]byte("You have too many commands enqueued. Please wait until some are processed.\n"))
			continue
		}

		c.inputChannel <- message

		c.conn.Write([]byte(message))

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
	close(c.inputChannel)

	c.conn.Close()

	c.server.game.EnqueueEvent(game.CreatePlayerLeaveEvent(c.Id()))

	logging.Info("Disconnected: ", c.conn.RemoteAddr())
}

func (c *Connection) NextInput() (input []byte, err error) {
	select {
	case val := <-c.inputChannel:
		return val, nil
	default:
		return nil, newInputEmptyError()
	}
}

func (c *Connection) Write(output []byte) (err error) {
	_, err = c.conn.Write(output)
	return
}
