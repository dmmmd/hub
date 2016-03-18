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
	acceptConnections(registry, dispatcher, listener)
}

func acceptConnections(registry *Registry, dispatcher *Dispatcher, listener *net.TCPListener) {
	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Can't accept connection: %s", err.Error())
		} else {
			conn := newConnection(connection)
			go serveConnection(registry, dispatcher, conn)
		}
	}
}

func serveConnection(registry *Registry, dispatcher *Dispatcher, connection ConnectionInterface) {
	client := newClient(registry.NextId(), connection)
	dispatcher.Subscribe(client)
	defer connection.Close()

	for {
		message, err := client.NextMessage()

		switch {
		case err == nil:
			dispatcher.Dispatch(message)
		case err.ConnectionError():
			dispatcher.Unsubscribe(client)
			// case err.InvalidMessage(): // Just continue
		}
	}
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
