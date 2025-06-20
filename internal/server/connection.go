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

const MaxLastSeenTime = 120 * time.Second

type Connection struct {
	ctx context.Context
	wg  *sync.WaitGroup

	identity uuid.UUID

	conn     *net.TCPConn
	lastSeen time.Time

	inputChannel chan string
}

func CreateConnection(conn *net.TCPConn, ctx context.Context, wg *sync.WaitGroup) (c *Connection) {
	logging.Info("Connect: ", conn.RemoteAddr())

	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(1 * time.Second)

	c = &Connection{
		ctx:          ctx,
		wg:           wg,
		identity:     uuid.New(),
		conn:         conn,
		inputChannel: make(chan string),
		lastSeen:     time.Now(),
	}

	c.wg.Add(2)
	go c.listen()
	go c.checkAlive()

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
			c.Write("Server shutting down, bye bye!\r\n")
			break
		}

		if time.Since(c.lastSeen) > MaxLastSeenTime {
			c.Write("You have been away for too long, bye bye!\r\n")
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
	c.conn.Close()

	logging.Info("Disconnected: ", c.conn.RemoteAddr())
}

func (c *Connection) NextInput() (input string, err error) {
	select {
	case val := <-c.inputChannel:
		return val, nil
	default:
		return "", newInputEmptyError()
	}
}

func (c *Connection) Write(output string) (err error) {
	_, err = c.conn.Write([]byte(output))
	return
}
