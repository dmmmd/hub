package main

import (
	"fmt"
)

type Dispatcher struct {
	name string
	//output chan string
	nextId int
	clients [255]*Client
}

func NewDispather(foo string)	*Dispatcher  {
	d := & Dispatcher{name: foo}
	//d.output = output
	//go d.dispatch()
	return d
}

func (d *Dispatcher) Dispatch(message string) {
	//d.output <- "Dispatching"

	//d.output <- fmt.Sprintf("message is '%s'", message)
	fmt.Printf("Dispatching message '%s'\n", message)
	for i, receiver := range d.clients {
		if (receiver == nil) {
			continue
		}

		fmt.Printf("\tto client %d\n", i)
		notSent++
		go receiver.Send(message)
	}
}

func (d *Dispatcher) Subscribe(c *Client) {
	d.nextId++
	id := d.nextId
	d.clients[id] = c
}
