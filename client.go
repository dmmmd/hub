package main

import (
	"strconv"
	"strings"
)

const BoundaryPrefix string = "Boundary < "

type ClientInterface interface {
	Id() int64
	Send(message string)
	NextMessage() (MessageInterface, *ClientError)
}

type Client struct {
	id         int64
	outbox     chan string
	connection ConnectionInterface
}

func newClient(id int64, connection ConnectionInterface) *Client {
	return &Client{id: id, connection: connection, outbox: make(chan string, 255)}
}

func (c *Client) Id() int64 {
	return c.id
}

func (c *Client) Send(message string) {
	c.outbox <- message // To send them in order
	go func() {
		c.connection.Write(<-c.outbox)
	}()
}

func (c *Client) NextMessage() (MessageInterface, *ClientError) {
	cmd, err := c.readCommand()
	if err != nil {
		return nil, err
	}

	switch cmd {
	case MessageTypeIdentity:
		return newIdentityMessage(c.id), nil
	case MessageTypeList:
		return newListMessage(c.id), nil
	}

	// Now we know it's relay

	// Receivers
	receivers, err := c.readReceivers()
	if err != nil {
		return nil, err
	}

	// Body
	body, err := c.readBody()
	if err != nil {
		return nil, err
	}

	return newRelayMessage(c.id, receivers, body), nil
}

func (c *Client) readCommand() (string, *ClientError) {
	line, err := c.readLine()
	if err != nil {
		return "", err
	}

	switch strings.TrimSpace(line) {
	case MessageTypeIdentity:
		return MessageTypeIdentity, nil
	case MessageTypeList:
		return MessageTypeList, nil
	case MessageTypeRelay:
		return MessageTypeRelay, nil
	}

	return "", NewClientInvalidCommandError()
}

func (c *Client) readReceivers() ([]int64, *ClientError) {
	line, err := c.readLine()
	if err != nil {
		return []int64{}, err
	}

	return c.parseReceivers(line)
}

func (c *Client) readBody() (string, *ClientError) {
	var message string

	boundary, err := c.getBoundary()
	if err != nil {
		return "", err
	}

	for {
		line, err := c.readLine()
		if err != nil {
			return "", err
		}

		if boundary == line {
			return message, nil
		}

		message += line
	}
}

func (c *Client) getBoundary() (string, *ClientError) {
	line, err := c.readLine()
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(line, BoundaryPrefix) {
		return line[len(BoundaryPrefix):], nil
	} else {
		return "", NewClientInvalidCommandError()
	}
}

func (c *Client) readLine() (string, *ClientError) {
	line, err := c.connection.Read()

	if err != nil {
		return "", NewClientConnectionError()
	}

	return line, nil
}

func (c *Client) parseReceivers(line string) ([]int64, *ClientError) {
	var receivers []int64
	for _, word := range strings.Split(line, ",") {
		word = strings.TrimSpace(word)
		id, err := strconv.ParseInt(word, 10, 64)
		if err != nil {
			return make([]int64, 0), NewClientInvalidReceiversError()
		}

		receivers = append(receivers, id)
	}
	return receivers, nil
}
