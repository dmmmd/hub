package main

import (
	"fmt"
)

type Client struct {
	id int64
	//output chan string
	dispatcher *Dispatcher
	outbox     chan string
}

func NewClient(id int64, d *Dispatcher) *Client {
	c := &Client{id: id, dispatcher: d, outbox: make(chan string, 255)}
	return c
}

func (c *Client) Say(message *Message) {
	//c.output <- fmt.Sprintf("Client %d says '%s' to %d", c.id, message, id)
	//receivers := strings.Join(message.Receivers(), ", ")
	receivers := message.Receivers()
	fmt.Printf("Client %d says '%s' to %d receivers\n", c.id, message.Body(), len(receivers))
	c.dispatcher.Dispatch(message)
}

func (c *Client) Send(message string) {
	c.outbox <- message
	go func() {
		message := <-c.outbox
		fmt.Printf("Client %d receiving %s\n", c.id, message)
		notSent--
	}()
	//c.output <- fmt.Sprintf("Client %d receiving %s", c.id, message)
}
