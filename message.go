package main

const MessageTypeRelay string = "relay"
const MessageTypeIdentity string = "identity"
const MessageTypeList string = "list"

type MessageInterface interface {
	Command() string
	Sender() uint64
	Receivers() []uint64
	Body() *string
}

type Message struct {
	command   string
	sender    uint64
	receivers []uint64
	body      *string
}

func newRelayMessage(sender uint64, receivers []uint64, body *string) *Message {
	return &Message{command: MessageTypeRelay, body: body, sender: sender, receivers: receivers}
}

func newIdentityMessage(sender uint64) *Message {
	return &Message{command: MessageTypeIdentity, sender: sender}
}

func newListMessage(sender uint64) *Message {
	return &Message{command: MessageTypeList, sender: sender}
}

func (m *Message) Command() string {
	return m.command
}

func (m *Message) Sender() uint64 {
	return m.sender
}

func (m *Message) Receivers() []uint64 {
	return m.receivers
}

func (m *Message) Body() *string {
	return m.body
}
