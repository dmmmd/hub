package main

import (
	"fmt"
)

func main() {
	output := make(chan string, 255)

	d := NewDispather(output, "testDispatcher")

	c1 := NewClient(1, output, d)
	c2 := NewClient(2, output, d)
	c3 := NewClient(3, output, d)

	c1.Say(42, "first message")

	d.Subscribe(c1)
	d.Subscribe(c2)
	d.Subscribe(c3)

	//fmt.Printf("c1 id is %d\n", c1.id)
	//fmt.Printf("c2 id is %d\n", c2.id)
	//fmt.Printf("c3 id is %d\n", c3.id)

	fmt.Print("Talking\n")
	c1.Say(42, "second message")

	//for message := range output {
	//	fmt.Printf("%s\n", message)
	//}
}
