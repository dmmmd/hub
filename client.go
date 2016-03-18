package main

import (
	"net"
	"strconv"
	"strings"
)

type ClientInterface interface {
	Id() int64
	Send(message string)
	NextMessage() (MessageInterface, *ClientError)
}

type Client struct {
	id         int64
	outbox     chan string
	connection *net.TCPConn
}

func NewClient(id int64, connection *net.TCPConn) *Client {
	return &Client{id: id, connection: connection, outbox: make(chan string, 255)}
}

func (c *Client) Id() int64 {
	return c.id
}

func (c *Client) Send(message string) {
	c.outbox <- message
	go func() {
		message := <-c.outbox
		c.connection.Write([]byte(message))
	}()
}

func (c *Client) NextMessage() (MessageInterface, *ClientError) {
	cmd, err := c.readLine()
	if err != nil {
		return nil, err
	}

	switch cmd {
	case MessageTypeIdentity:
		return NewIdentityMessage(c.id), nil
	case MessageTypeList:
		return NewListMessage(c.id), nil
	}

	if cmd != MessageTypeRelay {
		return nil, NewClientInvalidMessageError()
	}

	// Receivers
	line, err := c.readLine()
	if err != nil {
		return nil, err
	}
	receivers := parseReceivers(line)

	// Body
	body, err := c.readBody()
	if err != nil {
		return nil, err
	}

	return NewRelayMessage(c.id, []int64(receivers), body), nil
}

func (c *Client) readBody() (string, *ClientError) {
	var body string
	emptyLinesAmount := 0

	for emptyLinesAmount < 2 {
		line, err := c.readLine()
		if err != nil {
			return "", err
		}

		body += line + "\n"
		if "" == line {
			emptyLinesAmount++
		} else {
			emptyLinesAmount = 0
		}
	}

	return strings.TrimSpace(body), nil
}

func (c *Client) readLine() (string, *ClientError) {
	bytes := make([]byte, 1048576)
	len, err := c.connection.Read(bytes)

	if err != nil {
		c.connection.Close()
		return "", NewClientConnectionLostError()
	}

	return strings.TrimSpace(string(bytes[:len])), nil
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
