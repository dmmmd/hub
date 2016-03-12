package main

type Message struct {
	receivers []string
	body string
}

func NewMessage(receivers []string, body string) *Message {
	return &Message{body: body, receivers: receivers}
}

func (m *Message) Receivers() []string {
	return m.receivers
}

func (m *Message) Body() string {
	return m.body
}
