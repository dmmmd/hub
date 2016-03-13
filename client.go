package main

import (
	"net"
	"strconv"
	"strings"
)

type Client struct {
	id         int64
	outbox     chan string
	connection *net.TCPConn
}

func NewClient(id int64, connection *net.TCPConn) *Client {
	return &Client{id: id, connection: connection, outbox: make(chan string, 255)}
}

func (c *Client) Send(message string) {
	c.outbox <- message
	go func() {
		message := <-c.outbox
		c.connection.Write([]byte(message))
	}()
}

func (c *Client) NextMessage() *Message {
	cmd := c.readLine()

	switch cmd {
	case MessageTypeIdentity:
		return NewIdentityMessage(c.id)
	case MessageTypeList:
		return NewListMessage(c.id)
	}

	if cmd != MessageTypeRelay {
		// todo
	}

	receivers := parseReceivers(c.readLine())
	body := c.readLine()
	return NewRelayMessage(c.id, []int64(receivers), body)
}

func (c *Client) readLine() string {
	bytes := make([]byte, 1024)
	len, _ := c.connection.Read(bytes)

	//if err != nil {
	//	connection.Close()
	//	break
	//}

	return strings.TrimSpace(string(bytes[:len]))
}

func parseReceivers(line string) []int64 {
	var receivers []int64
	for _, word := range strings.Split(line, ",") {
		word = strings.TrimSpace(word)
		id, _ := strconv.ParseInt(word, 10, 64)
		receivers = append(receivers, id)
	}
	return receivers
}
