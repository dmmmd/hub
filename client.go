package main

import (
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
	message, err := c.readMessage()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(message, "\n")

	cmd := lines[0]
	switch cmd {
	case MessageTypeIdentity:
		return newIdentityMessage(c.id), nil
	case MessageTypeList:
		return newListMessage(c.id), nil
	}

	if cmd != MessageTypeRelay {
		return nil, NewClientInvalidCommandError()
	}

	if len(lines) < 3 {
		return nil, NewClientInvalidCommandError()
	}

	// Receivers
	receivers, err := parseReceivers(lines[1])
	if err != nil {
		return nil, err
	}

	// Body
	body := strings.Join(lines[2:], "\n")
	return newRelayMessage(c.id, []int64(receivers), body), nil
}

func (c *Client) readMessage() (string, *ClientError) {
	var message string
	emptyLinesAmount := 0

	for emptyLinesAmount < 2 {
		line, err := c.readLine()
		if err != nil {
			return "", err
		}

		message += line + "\n"
		if "" == line {
			emptyLinesAmount++
		} else {
			emptyLinesAmount = 0
		}
	}

	return strings.TrimSpace(message), nil
}

func (c *Client) readLine() (string, *ClientError) {
	line, err := c.connection.Read()

	if err != nil {
		return "", NewClientConnectionError()
	}

	return strings.TrimSpace(line), nil
}

func parseReceivers(line string) ([]int64, *ClientError) {
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
