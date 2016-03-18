package main

import (
	"fmt"
	"strings"
	"sync"
)

type DispatcherInterface interface {
	Dispatch(message MessageInterface)
	Subscribe(client ClientInterface)
	Unsubscribe(client ClientInterface)
}

type Dispatcher struct {
	clients      map[int64]ClientInterface
	clientsMutex sync.Mutex
}

func newDispather() *Dispatcher {
	return &Dispatcher{clients: make(map[int64]ClientInterface)}
}

func (d *Dispatcher) Dispatch(message MessageInterface) {
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
	body := fmt.Sprintf("[Server] Your ID is %d\n", sender)
	d.sendBody(sender, &body)
}

func (d *Dispatcher) list(sender int64) {
	var clientList []string

	d.lockClients()
	for id, _ := range d.clients {
		if id != sender {
			clientList = append(clientList, fmt.Sprintf("%d", id))
		}
	}
	d.unlockClients()

	body := fmt.Sprintf("[Server] Client IDs are %s\n", strings.Join(clientList, ", "))
	d.sendBody(sender, &body)
}

func (d *Dispatcher) relay(message MessageInterface) {
	for _, id := range message.Receivers() {
		d.sendBody(id, message.Body())
	}
}

func (d *Dispatcher) Subscribe(c ClientInterface) {
	d.lockClients()
	d.clients[c.Id()] = c
	d.unlockClients()
}

func (d *Dispatcher) Unsubscribe(c ClientInterface) {
	d.lockClients()
	delete(d.clients, c.Id())
	d.unlockClients()
}

func (d *Dispatcher) sendBody(receiver int64, body *string) {
	client := d.client(receiver)

	if client != nil {
		client.Send(body)
	}
}

func (d *Dispatcher) client(id int64) ClientInterface {
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
