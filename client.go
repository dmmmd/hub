package main

import (
	"strconv"
	"strings"
)

const BoundaryPrefix string = "Boundary: "

type ClientInterface interface {
	Id() uint64
	Send(message *string)
	NextMessage() (MessageInterface, *ClientError)
}

type Client struct {
	id         uint64
	outbox     chan *string
	connection ConnectionInterface
}

func newClient(id uint64, connection ConnectionInterface) *Client {
	return &Client{id: id, connection: connection, outbox: make(chan *string, 255)}
}

func (c *Client) Id() uint64 {
	return c.id
}

func (c *Client) Send(message *string) {
	c.outbox <- message // To send them in order
	go func() {
		message := <-c.outbox
		c.connection.Write(*message)
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
	case MessageTypeRelay:
		return c.buildRelayMessage()
	}

	return nil, newClientInvalidCommandError()
}

func (c *Client) buildRelayMessage() (MessageInterface, *ClientError) {
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

	return strings.TrimSpace(line), nil
}

func (c *Client) readReceivers() ([]uint64, *ClientError) {
	line, err := c.readLine()
	if err != nil {
		return []uint64{}, err
	}

	return c.parseReceivers(line)
}

func (c *Client) readBody() (*string, *ClientError) {
	var message string

	boundary, err := c.getBoundary()
	if err != nil {
		return new(string), err
	}

	for {
		line, err := c.readLine()
		if err != nil {
			return new(string), err
		}

		if boundary == line {
			return &message, nil
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
		return "", newClientInvalidCommandError()
	}
}

func (c *Client) readLine() (string, *ClientError) {
	line, err := c.connection.Read()

	if err != nil {
		return "", newClientConnectionError()
	}

	return line, nil
}

func (c *Client) parseReceivers(line string) ([]uint64, *ClientError) {
	var receivers []uint64
	for _, word := range strings.Split(line, ",") {
		word = strings.TrimSpace(word)
		id, err := strconv.ParseUint(word, 10, 64)
		if err != nil {
			return make([]uint64, 0), newClientInvalidReceiversError()
		}

		receivers = append(receivers, id)
	}
	return receivers, nil
}
