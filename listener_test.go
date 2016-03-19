package main

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestListsAndIdentities(t *testing.T) {
	port := 1234

	l := NewListener(port)
	go l.Listen()

	input1, output1 := getHelperClient(port)
	input2, output2 := getHelperClient(port)
	input3, output3 := getHelperClient(port)

	input1 <- MessageTypeIdentity

	input2 <- MessageTypeList
	input2 <- MessageTypeList

	input3 <- MessageTypeList
	input3 <- MessageTypeIdentity
	input3 <- MessageTypeIdentity

	assert := assert.New(t)
	finished := make(chan int, 3) // 3 Clients

	go func() {
		r := <-output1
		assert.Contains(r, "Your ID is 1")
		finished <- 1
	}()

	go func() {
		r := <-output2
		assert.Contains(r, "Client IDs are")
		assert.Contains(r, "1")
		assert.NotContains(r, "2")
		assert.Contains(r, "3")

		r = <-output2
		assert.Contains(r, "Client IDs are")
		assert.Contains(r, "1")
		assert.NotContains(r, "2")
		assert.Contains(r, "3")

		finished <- 2
	}()

	go func() {
		r := <-output3
		assert.Contains(r, "Client IDs are")
		assert.Contains(r, "1")
		assert.Contains(r, "2")
		assert.NotContains(r, "3")

		r = <-output3
		assert.Contains(r, "Your ID is 3")

		r = <-output3
		assert.Contains(r, "Your ID is 3")

		finished <- 3
	}()

	<-finished
	<-finished
	<-finished
}

func TestRelays(t *testing.T) {
	port := 1235

	l := NewListener(port)
	go l.Listen()

	input1, output1 := getHelperClient(port)
	input2, output2 := getHelperClient(port)
	input3, output3 := getHelperClient(port)

	body1 := "test message 1"
	body2 := "test END message 2\numad?"
	body3 := "test message 3\n\nEND"

	message1 := fmt.Sprintf("%s\n3\n%s", MessageTypeRelay, body1)
	message2 := fmt.Sprintf("%s\n2,1\n%s", MessageTypeRelay, body2)
	message3 := fmt.Sprintf("%s\ninvalid receivers\nignored body", MessageTypeRelay)
	message4 := fmt.Sprintf("invalid type\n1\nignored body")
	message5 := fmt.Sprintf("%s\n1,100500\n%s", MessageTypeRelay, body3)

	input1 <- message1
	input3 <- message2
	input1 <- message3
	input2 <- message4
	input3 <- message5

	assert := assert.New(t)
	finished := make(chan int, 3) // 3 Clients

	go func() {
		r := <-output1
		//fmt.Printf("RESULT: c1: '%s' -----c1-----\n", r)
		assert.Contains(r, body2)

		r = <-output1
		//fmt.Printf("RESULT: c1: '%s' -----c1-----\n", r)
		assert.Contains(r, body3)

		finished <- 1
	}()

	go func() {
		r := <-output2
		//fmt.Printf("RESULT: c2: '%s' -----c2-----\n", r)
		assert.Contains(r, body2)

		finished <- 2
	}()

	go func() {
		r := <-output3
		//fmt.Printf("RESULT: c3: '%s' -----c3-----\n", r)
		assert.Contains(r, body1)

		finished <- 3
	}()

	<-finished
	<-finished
	<-finished
}

/*
 * Helpers
 */

func getHelperClient(port int) (chan string, chan string) {
	c, _ := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	input := make(chan string, 255)
	output := make(chan string, 255)

	bytes := make([]byte, 100500)
	var len int
	reader := bufio.NewReader(c)

	go func() {
		for {
			// Sending
			m := <-input
			fmt.Fprint(c, m+"\n")

			// Reading response
			//fmt.Printf("About to read\n")
			len, _ = reader.Read(bytes)
			output <- string(bytes[:len])
		}
	}()

	return input, output
}
