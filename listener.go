package main

import (
	"fmt"
	"net"
	"os"
)

type Listener struct {
	port   int
	server *Server
	stop   bool
}

func NewListener(port int) *Listener {
	server := newServer(newClientFactory(newIdSequence()), newDispather())
	return &Listener{port: port, server: server, stop: false}
}

func (l *Listener) Listen() {
	listener := l.createListener()
	defer listener.Close()

	for !l.stop {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Can't accept connection: %s", err.Error())
		} else {
			conn := newConnection(connection)
			go l.server.Serve(conn)
		}
	}
}

func (l *Listener) createListener() *net.TCPListener {
	ip := net.IPv4(127, 0, 0, 1)
	addr := &net.TCPAddr{Port: l.port, IP: ip}
	listener, err := net.ListenTCP("tcp", addr)

	if err != nil {
		fmt.Printf("Can't start listening: %s", err.Error())
		os.Exit(1)
	}

	return listener
}

func (l *Listener) Stop() {
	l.stop = true
}
