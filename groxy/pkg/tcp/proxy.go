package tcp

import (
	"io"
	"net"
	"time"

	log "github.com/richpeaua/gomplementations/groxy/pkg/log"

)

// Proxy forwards TCP requests to a destination address
type Proxy struct {
	destination *net.TCPAddr
	terminationDelay time.Duration

}

type NewProxy(address string, terminationDelay time.Duration) (*Proxy, error) {
	dstAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	return &Proxy{ destination: dstAddress, terminationDelay: terminationDelay }, nil
}

func TCPConn(conn net.Conn) {
	

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