package main

import (
	"fmt"
	"sync"
)

type Dispatcher struct {
	clients      map[int64]*Client
	clientsMutex sync.Mutex
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

	d.lockClients()
	for id, _ := range d.clients {
		if id != sender {
			clientList += fmt.Sprintf("%d, ", id)
		}
	}
	d.unlockClients()

	d.sendBody(sender, fmt.Sprintf("Client IDs are %s", clientList))
}

func (d *Dispatcher) relay(message *Message) {
	for _, id := range message.Receivers() {
		d.sendBody(id, message.Body())
	}
}

func (d *Dispatcher) Subscribe(c *Client) {
	d.lockClients()
	d.clients[c.id] = c
	d.unlockClients()
}

func (d *Dispatcher) sendBody(receiver int64, body string) {
	client := d.client(receiver)

	if client != nil {
		client.Send(body + "\n")
	}
}

func (d *Dispatcher) client(id int64) *Client {
	d.lockClients()
	client := d.clients[id]
	d.unlockClients()
	return client
}

func (d *Dispatcher) lockClients() {
	d.clientsMutex.Lock()
}

func (d *Dispatcher) unlockClients() {
	d.clientsMutex.Unlock()
}
