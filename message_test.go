package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageImplementsMessageInterface(t *testing.T) {
	var _ MessageInterface = (*Message)(nil)
}

func TestRelayCommandReturnsRelay(t *testing.T) {
	m := newRelayMessage(9001, []int64{100500, 42}, "irrelevant")

	assert.Equal(t, MessageTypeRelay, m.Command())
}

func TestRelaySenderReturnsSender(t *testing.T) {
	sender := int64(42)
	m := newRelayMessage(sender, []int64{100500, 42}, "irrelevant")

	assert.Equal(t, sender, m.Sender())
}

func TestRelayBodyReturnsBody(t *testing.T) {
	body := "testBody\nline2 foobar"
	m := newRelayMessage(9001, []int64{100500, 42}, body)

	assert.Equal(t, body, m.Body())
}

func TestRelayReceiversReturnsReceivers(t *testing.T) {
	receivers := []int64{100500, 42}
	m := newRelayMessage(9001, receivers, "irrelevant")
	r := m.Receivers()

	assert.Len(t, r, 2)
	assert.EqualValues(t, 100500, r[0])
	assert.EqualValues(t, 42, r[1])
}

func TestIdentityCommandReturnsIdentity(t *testing.T) {
	m := newIdentityMessage(42)

	assert.Equal(t, MessageTypeIdentity, m.Command())
}

func TestIdentitySenderReturnsSender(t *testing.T) {
	sender := int64(42)
	m := newIdentityMessage(sender)

	assert.Equal(t, sender, m.Sender())
}

func TestIdentityBodyIsEmpty(t *testing.T) {
	m := newIdentityMessage(42)

	assert.Empty(t, m.Body())
}

func TestIdentityReceiversAreEmpty(t *testing.T) {
	m := newIdentityMessage(42)
	r := m.Receivers()

	assert.Empty(t, r)
}

func TestListCommandReturnsList(t *testing.T) {
	m := newListMessage(42)

	assert.Equal(t, MessageTypeList, m.Command())
}

func TestListSenderReturnsSender(t *testing.T) {
	sender := int64(42)
	m := newListMessage(sender)

	assert.Equal(t, sender, m.Sender())
}

func TestListBodyIsEmpty(t *testing.T) {
	m := newListMessage(42)

	assert.Empty(t, m.Body())
}

func TestListReceiversAreEmpty(t *testing.T) {
	m := newListMessage(42)
	r := m.Receivers()

	assert.Empty(t, r)
}
