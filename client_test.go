package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestClientImplementsClientInterface(t *testing.T) {
	var _ ClientInterface = (*Client)(nil)
}

func TestIdReturnsId(t *testing.T) {
	id := uint64(42)
	conn := new(ClientMockedConnection)

	c := newClient(id, conn)
	assert.Equal(t, id, c.Id())
}

func TestNextMessageReturnsIdentityMessage(t *testing.T) {
	conn := new(ClientMockedConnection)

	conn.addExpectedLine(MessageTypeIdentity)

	c := newClient(42, conn)
	message, _ := c.NextMessage()

	assert.Equal(t, MessageTypeIdentity, message.Command())
}

func TestNextMessageReturnsListMessage(t *testing.T) {
	conn := new(ClientMockedConnection)

	conn.addExpectedLine(MessageTypeList)

	c := newClient(42, conn)
	message, _ := c.NextMessage()

	assert.Equal(t, MessageTypeList, message.Command())
}

func TestNextMessageReturnsRelayMessage(t *testing.T) {
	boundary := "testBoundary"

	conn := new(ClientMockedConnection)

	conn.addExpectedLine(MessageTypeRelay)
	conn.addExpectedLine("100500,42,56")
	conn.addExpectedLine(BoundaryPrefix + boundary)
	conn.addExpectedLine("foobar 1")
	conn.addExpectedLine("foobar 2")
	conn.addExpectedLine("")
	conn.addExpectedLine("foobar 3")
	conn.addExpectedLine(boundary)

	c := newClient(42, conn)
	message, _ := c.NextMessage()

	assert := assert.New(t)

	assert.Equal(MessageTypeRelay, message.Command())
	assert.Equal("foobar 1\nfoobar 2\n\nfoobar 3\n", *message.Body())

	receivers := message.Receivers()
	assert.Len(receivers, 3)
	assert.Contains(receivers, uint64(42))
	assert.Contains(receivers, uint64(56))
	assert.Contains(receivers, uint64(100500))
}

func TestNextMessageReturnsErrorOnInvalidCommand(t *testing.T) {
	conn := new(ClientMockedConnection)

	conn.addExpectedLine("testInvalidCommand")
	conn.addExpectedLine("100500,42,56")
	conn.addExpectedLine(BoundaryPrefix + "TEST")
	conn.addExpectedLine("foobar")
	conn.addExpectedLine("TEST")

	c := newClient(42, conn)
	message, err := c.NextMessage()

	assert := assert.New(t)

	assert.Nil(message)
	assert.True(err.InvalidCommand())
	assert.False(err.InvalidReceivers())
	assert.False(err.ConnectionError())
}

func TestNextMessageReturnsErrorOnInvalidReceivers(t *testing.T) {
	conn := new(ClientMockedConnection)

	conn.addExpectedLine(MessageTypeRelay)
	conn.addExpectedLine("100500,4foo2,56")
	conn.addExpectedLine(BoundaryPrefix + "TEST")
	conn.addExpectedLine("foobar")
	conn.addExpectedLine("TEST")

	c := newClient(42, conn)
	message, err := c.NextMessage()

	assert := assert.New(t)

	assert.Nil(message)
	assert.False(err.InvalidCommand())
	assert.True(err.InvalidReceivers())
	assert.False(err.ConnectionError())
}

func TestNextMessageReturnsErrorOnReadError(t *testing.T) {
	conn := new(ClientMockedConnection)

	c := newClient(42, conn)
	message, err := c.NextMessage()

	assert := assert.New(t)

	assert.Nil(message)
	assert.False(err.InvalidCommand())
	assert.False(err.InvalidReceivers())
	assert.True(err.ConnectionError())
}

func TestSendWritesToConnection(t *testing.T) {
	messages := []string{"testMessage1", "test\nMessage2"}

	conn := new(ClientMockedConnection)
	conn.On("Write", messages[0]).Return(nil)
	conn.On("Write", messages[1]).Return(nil)

	c := newClient(42, conn)
	c.Send(&messages[0])
	c.Send(&messages[1])

	time.Sleep(time.Millisecond) // For stupidity points

	conn.AssertNumberOfCalls(t, "Write", 2)
	conn.AssertExpectations(t)
	assert.Len(t, conn.written, 2)
}

/*
 * Mocks
 */

type ClientMockedConnection struct {
	mock.Mock
	lines   []string
	written []string
}

func (c *ClientMockedConnection) Write(message string) error {
	args := c.Called(message)
	c.written = append(c.written, message)
	return args.Error(0)
}

func (c *ClientMockedConnection) Read() (string, error) {
	if len(c.lines) > 0 {
		line := c.lines[0]
		c.lines = c.lines[1:]
		return line, nil
	}

	return "", errors.New("testConnectionReadError")
}

func (c *ClientMockedConnection) Close() {
	c.Called()
}

func (c *ClientMockedConnection) addExpectedLine(line string) {
	c.lines = append(c.lines, line+"\n")
}
