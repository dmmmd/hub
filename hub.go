package main

var notSent = 0

func main() {
	//output := make(chan string, 255)

	d := NewDispather("testDispatcher")

	c1 := NewClient("One", d)
	c2 := NewClient("Two", d)
	c3 := NewClient("Three", d)

	message1 := NewMessage([]string{"Two", "Three"}, "first message 1")
	message2 := NewMessage([]string{"One", "Two", "Three"}, "second message 2")
	message3 := NewMessage([]string{"One", "Three"}, "third message 3")
	message4 := NewMessage([]string{"One", "Two"}, "fourth message 4")

	c1.Say(message1)

	d.Subscribe(c1)
	d.Subscribe(c2)
	d.Subscribe(c3)

	c1.Say(message2)
	c2.Say(message3)
	c3.Say(message4)

	//fmt.Printf("Due for sending: %d\n", notSent)

	for notSent > 0 {
	}

	//fmt.Printf("Due for sending: %d\n", notSent)
}
