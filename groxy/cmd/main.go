package main

import (
	"bytes"
	"encoding/binary"
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

func proxy(conn net.Conn) {
	defer conn.Close()

	upstream, err := net.Dial("tcp", "google.com:443")
	if err != nil {
		log.Print(err)
		return
	}
	defer upstream.Close()

	// io and net packages use splice system call to copy data directly between two file descriptors
	go io.Copy(upstream, conn) // should track go routines, it's fine here because of the defer.
	io.Copy(conn, upstream)

}

func copyToStderr(conn net.Conn) {
	defer conn.Close()
	for {
		var buff [128]byte
		conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // timeout to kill idle connections to release kernel resources i.e. file desciptors
		n, err := conn.Read(buff[:])
		if err != nil {
			log.Print(err)
			return
		}
		os.Stderr.Write(buff[:n])
	}
}

