package main

import (
	"strconv"
	"strings"
)

type ClientInterface interface {
	Id() uint64
	Send(message *string)
	NextMessage() (MessageInterface, *ClientError)
}

type ClientFactoryInterface interface {
	Create(connection ConnectionInterface) ClientInterface
}

type ClientFactory struct {
	sequence IdSequenceInterface
}

func newClientFactory(s IdSequenceInterface) *ClientFactory {
	return &ClientFactory{sequence: s}
}

func (f *ClientFactory) Create(connection ConnectionInterface) ClientInterface {
	return newClient(f.sequence.NextId(), connection)
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
	command, err := c.readCommand()
	if err != nil {
		return nil, err
	}

	cmd := strings.SplitN(command, "\n", 3)

	switch cmd[0] {
	case MessageTypeIdentity:
		return newIdentityMessage(c.id), nil
	case MessageTypeList:
		return newListMessage(c.id), nil
	case MessageTypeRelay:
		if len(cmd) != 3 {
			return nil, newClientInvalidCommandError()
		}
		return c.buildRelayMessage(cmd[1], &cmd[2])
	}

	return nil, newClientInvalidCommandError()
}

func (c *Client) buildRelayMessage(receivers string, body *string) (MessageInterface, *ClientError) {
	receiverIds, err := c.parseReceivers(receivers)
	if err != nil {
		return nil, err
	}

	return newRelayMessage(c.id, receiverIds, body), nil
}

func (c *Client) readCommand() (string, *ClientError) {
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
