package main

import (
	"flag"
)

func main() {
	port := flag.Int("port", 8888, "Port to listen")
	flag.Parse()

	l := NewListener(*port)
	l.Listen()
}
