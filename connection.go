package main

import "net"

type ConnectionInterface interface {
	Write(string) error
	Read() (string, error)
	Close()
}

type Connection struct {
	conn *net.TCPConn
}

func newConnection(c *net.TCPConn) *Connection {
	return &Connection{conn: c}
}

func (c *Connection) Write(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

func (c *Connection) Read() (string, error) {
	bytes := make([]byte, 1100000)
	len, err := c.conn.Read(bytes)

	if err != nil {
		return "", err
	}

	return string(bytes[:len]), nil
}

func (c *Connection) Close() {
	c.conn.Close()
}
