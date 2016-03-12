package main

import (
	"fmt"
)

var notSent = 0

func main() {
	//output := make(chan string, 255)

	d := NewDispather("testDispatcher")

	c1 := NewClient(1, d)
	c2 := NewClient(2, d)
	c3 := NewClient(3, d)

	c1.Say(42, "first message")

	d.Subscribe(c1)
	d.Subscribe(c2)
	d.Subscribe(c3)

	//fmt.Printf("c1 id is %d\n", c1.id)
	//fmt.Printf("c2 id is %d\n", c2.id)
	//fmt.Printf("c3 id is %d\n", c3.id)

	fmt.Print("Talking\n")
	c1.Say(42, "second message")
	c2.Say(42, "third message")

	//fmt.Printf("Due for sending: %d\n", notSent)

	//for {
	//	for message := range output {
	//		fmt.Printf("[OUTPUT] %s\n", message)
	//	}
	//
	//}

	for notSent > 0 {
	}

	//fmt.Printf("Due for sending: %d\n", notSent)
}
