package main

import (
	"fmt"
	"net"
	"os"
	"flag"
)

func main() {
	port := flag.Int("port", 8888, "Port to listen")
	flag.Parse()

	NewHub(*port)
}

type Hub struct {
	port int
	server *Server
}

func NewHub(port int) *Hub {
	server := newServer(newClientFactory(newIdSequence()), newDispather())
	hub := &Hub{port: port, server: server}
	hub.acceptConnections()

	return hub
}

func (h *Hub) acceptConnections() {
	listener := h.createListener()
	defer listener.Close()

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Can't accept connection: %s", err.Error())
		} else {
			conn := newConnection(connection)
			go h.server.Serve(conn)
		}
	}
}

func (h *Hub) createListener() *net.TCPListener {
	ip := net.IPv4(127, 0, 0, 1)
	addr := &net.TCPAddr{Port: h.port, IP: ip}
	listener, err := net.ListenTCP("tcp", addr)

	if err != nil {
		fmt.Printf("Can't start listening: %s", err.Error())
		os.Exit(1)
	}

	return listener
}
