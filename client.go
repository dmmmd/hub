package main

import "fmt"

type Client struct {
	id int
	//output chan string
	dispatcher *Dispatcher
}

func NewClient(id int, d *Dispatcher) *Client {
	c := & Client{id: id, dispatcher: d}
	return c
}

func (c *Client) Say(id int, message string) {
	//c.output <- fmt.Sprintf("Client %d says '%s' to %d", c.id, message, id)
	fmt.Printf("Client %d says '%s' to %d\n", c.id, message, id)
	c.dispatcher.Dispatch(message)
}

func (c *Client) Send(message string) {
	//c.output <- fmt.Sprintf("Client %d receiving %s", c.id, message)
	fmt.Printf("Client %d receiving %s\n", c.id, message)
	notSent--
}