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

	sequence := newIdSequence()
	clientFactory := newClientFactory(sequence)
	dispatcher := newDispather()
	server := newServer(clientFactory, dispatcher)

	acceptConnections(server, listener)
}

func acceptConnections(server *Server, listener *net.TCPListener) {
	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Can't accept connection: %s", err.Error())
		} else {
			conn := newConnection(connection)
			go server.Serve(conn)
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
