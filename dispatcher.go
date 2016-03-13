package main

import (
	"fmt"
)

type Dispatcher struct {
	clients map[int64]*Client
}

func NewDispather() *Dispatcher {
	return &Dispatcher{clients: make(map[int64]*Client)}
}

func (d *Dispatcher) Dispatch(message *Message) {
	switch message.Command() {
	case MessageTypeRelay:
		d.relay(message)
	case MessageTypeIdentity:
		d.identify(message.Sender())
	case MessageTypeList:
		d.list(message.Sender())
	}

}

func (d *Dispatcher) identify(sender int64) {
	d.sendBody(sender, fmt.Sprintf("Your ID is %d", sender))
}

func (d *Dispatcher) list(sender int64) {
	var clientList string
	for id, _ := range d.clients {
		if id != sender {
			clientList += fmt.Sprintf("%d, ", id)
		}
	}

	d.sendBody(sender, fmt.Sprintf("Client IDs are %s", clientList))
}

func (d *Dispatcher) relay(message *Message) {
	for _, id := range message.Receivers() {
		d.sendBody(id, message.Body())
	}
}

func (d *Dispatcher) Subscribe(c *Client) {
	d.clients[c.id] = c
}

func (d *Dispatcher) sendBody(receiver int64, body string) {
	if client := d.clients[receiver]; client != nil {
		client.Send(body + "\n")
	}

}
