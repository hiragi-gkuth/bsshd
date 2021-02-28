package main

import (
	"fmt"
	"log"
	"net"
)

func setupServer(addr string, port int) net.Listener {
	listener, e := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if e != nil {
		log.Fatal("failed to listen for connection: ", e)
	}
	log.Printf("bsshd is listening on %v", port)
	return listener
}

func connHandler(listener net.Listener, connCh chan net.Conn) {
	for { // wait for connection
		conn, e := listener.Accept()
		if e != nil {
			log.Fatalf("failed to accept tcp connection %v", e)
			continue
		}
		// accept connection
		log.Printf("accept new connection from %v\n", conn.RemoteAddr().String())

		// send chan
		connCh <- conn
	}
}
