package server

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
)

type Connection struct {
	identity uuid.UUID

	conn net.Conn

	inputChannel chan string
	closeChannel chan struct{}
}

func CreateConnection(conn net.Conn) *Connection {
	return &Connection{
		identity:     uuid.New(),
		conn:         conn,
		inputChannel: make(chan string),
		closeChannel: make(chan struct{}),
	}
}

func (c *Connection) listen() (err error) {
	defer c.conn.Close()

	c.conn.SetReadDeadline(time.Time{})

	for {
		if c.shouldClose() {
			break
		}

		message, err := bufio.NewReader(c.conn).ReadString('\n')

		if err != nil {
			fmt.Println(err)
			return err
		}

		c.inputChannel <- message

		c.conn.Write([]byte(message))
	}

	return
}

func (c *Connection) shouldClose() bool {
	select {
	case <-c.closeChannel:
		return true
	default:
		return false
	}
}

func (c *Connection) Close() {
	c.closeChannel <- struct{}{}
}

func (c *Connection) NextInput() (next string, err error) {
	select {
	case val := <-c.inputChannel:
		return val, nil
	default:
		return "", newInputEmptyError()
	}
}
