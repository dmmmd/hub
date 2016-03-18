package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDispatcherImplementsDispatcherInterface(t *testing.T) {
	var _ DispatcherInterface = (*Dispatcher)(nil)
}

func TestDispatchRespondsWithIdentity(t *testing.T) {
	clientId := int64(42)
	client := newMockedClient(clientId)

	d := newDispather()
	d.Subscribe(client)
	d.Dispatch(newIdentityMessage(clientId))

	assert := assert.New(t)
	assert.Len(client.received, 1)
	assertSentMessageContainsClient(assert, client.received[0], clientId)
}

func TestDispatchNotSendingIdentityToDisconnectedClient(t *testing.T) {
	clientId := int64(42)
	client := newMockedClient(clientId)

	d := newDispather()
	d.Subscribe(client)
	d.Unsubscribe(client)

	d.Dispatch(newIdentityMessage(clientId))

	assert.Len(t, client.received, 0)
}

func TestDispatchRespondsWithList(t *testing.T) {
	clientId1 := int64(42)
	clientId2 := int64(100500)
	clientId3 := int64(9001)

	client1 := newMockedClient(clientId1)
	client2 := newMockedClient(clientId2)
	client3 := newMockedClient(clientId3)

	assert := assert.New(t)

	d := newDispather()
	d.Subscribe(client1)
	d.Dispatch(newListMessage(clientId1))

	var m string

	assert.Len(client1.received, 1)
	assert.Len(client2.received, 0) // Nothing sent there
	assert.Len(client3.received, 0)
	m = client1.received[0]
	assertSentMessageNotContainsClient(assert, m, clientId1)
	assertSentMessageNotContainsClient(assert, m, clientId2)
	assertSentMessageNotContainsClient(assert, m, clientId3)

	// Now lets connect more clients

	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Dispatch(newListMessage(clientId1))

	assert.Len(client1.received, 2)
	assert.Len(client2.received, 0) // Still nothing sent there
	assert.Len(client3.received, 0)
	m = client1.received[1]
	assertSentMessageNotContainsClient(assert, m, clientId1)
	assertSentMessageContainsClient(assert, m, clientId2)
	assertSentMessageContainsClient(assert, m, clientId3)

	d.Dispatch(newListMessage(clientId2))

	assert.Len(client1.received, 2)
	assert.Len(client2.received, 1)
	assert.Len(client3.received, 0)
	m = client2.received[0]
	assertSentMessageContainsClient(assert, m, clientId1)
	assertSentMessageNotContainsClient(assert, m, clientId2)
	assertSentMessageContainsClient(assert, m, clientId3)
}

func TestDispatchRespondsWithShorterListAfterUnsubscribing(t *testing.T) {
	clientId1 := int64(42)
	clientId2 := int64(100500)
	clientId3 := int64(9001)

	client1 := newMockedClient(clientId1)
	client2 := newMockedClient(clientId2)
	client3 := newMockedClient(clientId3)

	assert := assert.New(t)

	d := newDispather()
	d.Subscribe(client1)
	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Dispatch(newListMessage(clientId1))

	var m string

	assert.Len(client1.received, 1)
	m = client1.received[0]
	assertSentMessageContainsClient(assert, m, clientId2)
	assertSentMessageContainsClient(assert, m, clientId3)

	// Now lets disconnect client2

	d.Unsubscribe(client2)

	d.Dispatch(newListMessage(clientId1))

	assert.Len(client1.received, 2)
	m = client1.received[1]
	assertSentMessageNotContainsClient(assert, m, clientId2) // Must go away
	assertSentMessageContainsClient(assert, m, clientId3)
}

func TestDispatchSendsRelayMessageToReceivers(t *testing.T) {
	body1 := "test message\nbody 1"
	body2 := "test message\nbody 2"
	body3 := "test message\nbody 3"

	clientId1 := int64(42)
	clientId2 := int64(100500)
	clientId3 := int64(9001)

	client1 := newMockedClient(clientId1)
	client2 := newMockedClient(clientId2)
	client3 := newMockedClient(clientId3)

	assert := assert.New(t)

	d := newDispather()
	d.Subscribe(client1)
	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Dispatch(newRelayMessage(clientId1, []int64{clientId2, clientId3, clientId1}, &body1))
	d.Dispatch(newRelayMessage(clientId1, []int64{clientId2}, &body2))
	d.Dispatch(newRelayMessage(clientId3, []int64{clientId2, clientId1}, &body3))

	assert.Len(client1.received, 2)
	assert.Equal(body1, client1.received[0])
	assert.Equal(body3, client1.received[1])

	assert.Len(client2.received, 3)
	assert.Equal(body1, client2.received[0])
	assert.Equal(body2, client2.received[1])
	assert.Equal(body3, client2.received[2])

	assert.Len(client3.received, 1)
	assert.Equal(body1, client3.received[0])
}

func TestDispatchIgnoresNonExistingReceivers(t *testing.T) {
	body1 := "test message\nbody 1"

	clientId1 := int64(42)
	clientId2 := int64(100500)
	clientId3 := int64(9001)

	client1 := newMockedClient(clientId1)
	client2 := newMockedClient(clientId2)
	client3 := newMockedClient(clientId3)

	assert := assert.New(t)

	d := newDispather()
	d.Subscribe(client1)
	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Dispatch(newRelayMessage(clientId1, []int64{clientId2, 987654321, clientId3}, &body1))

	assert.Len(client1.received, 0)

	assert.Len(client2.received, 1)
	assert.Equal(body1, client2.received[0])

	assert.Len(client3.received, 1)
	assert.Equal(body1, client3.received[0])
}

func TestDispatchNotSendingToUnsubscribedClients(t *testing.T) {
	body1 := "test message\nbody 1"

	clientId1 := int64(42)
	clientId2 := int64(100500)
	clientId3 := int64(9001)

	client1 := newMockedClient(clientId1)
	client2 := newMockedClient(clientId2)
	client3 := newMockedClient(clientId3)

	assert := assert.New(t)

	d := newDispather()
	d.Subscribe(client1)
	d.Subscribe(client2)
	d.Subscribe(client3)

	d.Unsubscribe(client2)

	d.Dispatch(newRelayMessage(clientId1, []int64{clientId2, clientId3}, &body1))

	assert.Len(client1.received, 0)
	assert.Len(client2.received, 0)

	assert.Len(client3.received, 1)
	assert.Equal(body1, client3.received[0])
}

/*
 * Helpers
 */

func assertSentMessageContainsClient(assert *assert.Assertions, sent string, clientId int64) {
	assert.Contains(sent, fmt.Sprintf("%d", clientId))
}

func assertSentMessageNotContainsClient(assert *assert.Assertions, sent string, clientId int64) {
	assert.NotContains(sent, fmt.Sprintf("%d", clientId))
}

/*
 * Mocks
 */

type MockedClient struct {
	id int64
	mock.Mock
	received []string
}

func newMockedClient(id int64) *MockedClient {
	return &MockedClient{id: id}
}

func (c *MockedClient) Id() int64 {
	return c.id
}

func (c *MockedClient) Send(message *string) {
	c.received = append(c.received, *message)
}

func (c *MockedClient) NextMessage() (MessageInterface, *ClientError) {
	return nil, nil // Irrelevant for this test
}
