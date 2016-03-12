package main

const MessageTypeRelay string = "relay"
const MessageTypeIdentity string = "identity"
const MessageTypeList string = "list"

type Message struct {
	command   string
	receivers []int64
	body      string
}

func NewRelayMessage(receivers []int64, body string) *Message {
	return &Message{command: MessageTypeRelay, body: body, receivers: receivers}
}

func NewIdentityMessage() *Message {
	return &Message{command: MessageTypeIdentity}
}

func NewListMessage() *Message {
	return &Message{command: MessageTypeList}
}

func (m *Message) Command() string {
	return m.command
}

func (m *Message) Receivers() []int64 {
	return m.receivers
}

func (m *Message) Body() string {
	//return "42"
	return m.body
}
