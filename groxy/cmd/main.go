package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// listener
	var listenAddress string
	listenAddress = "localhost:4242"
	log.Printf("Server starting up\n")
	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server up and listening on %s", listenAddress)

	// accept loop that looks for and accepts all connections. NO BLOCKING CODE INSIDE OF ACCEPT LOOP!
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)

		}
		go proxy(conn)

	}
}


