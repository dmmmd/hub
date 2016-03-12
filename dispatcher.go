package main

import (
	"fmt"
)

type Dispatcher struct {
	name string
	clients map[string]*Client
}

func NewDispather(foo string)	*Dispatcher  {
	return & Dispatcher{name: foo, clients: make(map[string]*Client)}
}

func (d *Dispatcher) Dispatch(message *Message) {
	//d.output <- "Dispatching"

	//d.output <- fmt.Sprintf("message is '%s'", message)
	fmt.Printf("Dispatching message '%s'\n", message.Body())
	for _, id := range message.Receivers() {
		if receiver := d.clients[id]; receiver != nil {
			fmt.Printf("\tto client %s\n", receiver.id)
			notSent++
			receiver.Send(message.Body())
		}
	}
}

func (d *Dispatcher) Subscribe(c *Client) {
	d.clients[c.id] = c
}
