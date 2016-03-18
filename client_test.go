package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
	"errors"
)

func TestClientImplementsClientInterface(t *testing.T) {
	var _ ClientInterface = (*Client)(nil)
}

func TestIdReturnsId(t *testing.T) {
	id := int64(42)
	conn := new(MockedConnection)

	c := newClient(id, conn)
	assert.Equal(t, id, c.Id())
}

func TestNextMessageReturnsIdentityMessage(t *testing.T) {
	conn := new(MockedConnection)

	conn.AddExpectedLine(MessageTypeIdentity)
	conn.AddExpectedLine("")
	conn.AddExpectedLine("")

	c := newClient(42, conn)
	message, _ := c.NextMessage()

	assert.Equal(t, MessageTypeIdentity, message.Command())
}

func TestNextMessageReturnsListMessage(t *testing.T) {
	conn := new(MockedConnection)

	conn.AddExpectedLine(MessageTypeList)
	conn.AddExpectedLine("")
	conn.AddExpectedLine("")

	c := newClient(42, conn)
	message, _ := c.NextMessage()

	assert.Equal(t, MessageTypeList, message.Command())
}

func TestNextMessageReturnsRelayMessage(t *testing.T) {
	conn := new(MockedConnection)

	conn.AddExpectedLine(MessageTypeRelay)
	conn.AddExpectedLine("100500,42,56")
	conn.AddExpectedLine("foobar 1")
	conn.AddExpectedLine("foobar 2")
	conn.AddExpectedLine("")
	conn.AddExpectedLine("foobar 3")
	conn.AddExpectedLine("")
	conn.AddExpectedLine("")

	c := newClient(42, conn)
	message, _ := c.NextMessage()

	assert := assert.New(t)

	assert.Equal(MessageTypeRelay, message.Command())
	assert.Equal("foobar 1\nfoobar 2\n\nfoobar 3", message.Body())

	receivers := message.Receivers()
	assert.Len(receivers, 3)
	assert.Contains(receivers, int64(42))
	assert.Contains(receivers, int64(56))
	assert.Contains(receivers, int64(100500))
}

func TestNextMessageReturnsErrorOnInvalidCommand(t *testing.T) {
	conn := new(MockedConnection)

	conn.AddExpectedLine("testInvalidCommand")
	conn.AddExpectedLine("100500,42,56")
	conn.AddExpectedLine("foobar")
	conn.AddExpectedLine("")
	conn.AddExpectedLine("")

	c := newClient(42, conn)
	message, err := c.NextMessage()

	assert := assert.New(t)

	assert.Nil(message)
	assert.True(err.InvalidCommand())
	assert.False(err.InvalidReceivers())
	assert.False(err.ConnectionError())
}

func TestNextMessageReturnsErrorOnInvalidReceivers(t *testing.T) {
	conn := new(MockedConnection)

	conn.AddExpectedLine(MessageTypeRelay)
	conn.AddExpectedLine("100500,4foo2,56")
	conn.AddExpectedLine("foobar")
	conn.AddExpectedLine("")
	conn.AddExpectedLine("")

	c := newClient(42, conn)
	message, err := c.NextMessage()

	assert := assert.New(t)

	assert.Nil(message)
	assert.False(err.InvalidCommand())
	assert.True(err.InvalidReceivers())
	assert.False(err.ConnectionError())
}

func TestNextMessageReturnsErrorOnReadError(t *testing.T) {
	conn := new(MockedConnection)

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

	conn := new(MockedConnection)
	conn.On("Write", messages[0]).Return(nil)
	conn.On("Write", messages[1]).Return(nil)

	c := newClient(42, conn)
	c.Send(messages[0])
	c.Send(messages[1])

	time.Sleep(time.Millisecond) // For stupidity points

	conn.AssertNumberOfCalls(t, "Write", 2)
	conn.AssertExpectations(t)
	assert.Len(t, conn.written, 2)
}

/*
 * Mocks
 */

type MockedConnection struct {
	mock.Mock
	lines []string
	written []string
}

func (c *MockedConnection) Write(message string) error {
	args := c.Called(message)
	c.written = append(c.written, message)
	return args.Error(0)
}

func (c *MockedConnection) Read() (string, error) {
	if len(c.lines) > 0 {
		line := c.lines[0]
		c.lines = c.lines[1:]
		return line, nil
	}

	return "", errors.New("testConnectionReadError")
}

func (c *MockedConnection) Close() {
	c.Called()
}

func (c *MockedConnection) AddExpectedLine(line string) {
	c.lines = append(c.lines, line + "\n")
}