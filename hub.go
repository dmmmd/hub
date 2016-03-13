package main

import (
	"fmt"
	"net"
	"os"
)

const serverPort = 8888

func main() {
	listener := createListener()
	defer listener.Close()

	registry := new(Registry)
	dispatcher := NewDispather()

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Can't accept connection: %s", err.Error())
			os.Exit(1)
		}

		go serveConnection(registry, dispatcher, connection)
	}
}

func serveConnection(registry *Registry, dispatcher *Dispatcher, connection *net.TCPConn) {
	client := NewClient(registry.NextId(), connection)
	dispatcher.Subscribe(client)

	for {
		message := client.NextMessage()
		dispatcher.Dispatch(message)
	}

	connection.Close()
}

func createListener() *net.TCPListener {
	ip := net.IPv4(127, 0, 0, 1)
	addr := &net.TCPAddr{Port: serverPort, IP: ip}
	listener, err := net.ListenTCP("tcp", addr)

	if err != nil {
		fmt.Printf("Can't start listening: %s", err.Error())
		os.Exit(1)
	}

	return listener
}
