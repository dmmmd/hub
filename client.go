package main

import (
	"fmt"
	"strings"
)

type Client struct {
	id string
	//output chan string
	dispatcher *Dispatcher
	outbox chan string
}

func NewClient(id string, d *Dispatcher) *Client {
	c := & Client{id: id, dispatcher: d, outbox: make(chan string, 255)}
	return c
}

func (c *Client) Say(message *Message) {
	//c.output <- fmt.Sprintf("Client %d says '%s' to %d", c.id, message, id)
	receivers := strings.Join(message.Receivers(), ", ")
	fmt.Printf("Client %s says '%s' to %s\n", c.id, message.Body(), receivers)
	c.dispatcher.Dispatch(message)
}

func (c *Client) Send(message string) {
	c.outbox <- message
	go func() {
		message := <- c.outbox
		fmt.Printf("Client %s receiving %s\n", c.id, message)
		notSent--
	}()
	//c.output <- fmt.Sprintf("Client %d receiving %s", c.id, message)
}