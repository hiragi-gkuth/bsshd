package main

import (
	"log"
	"net"

	"github.com/hiragi-gkuth/bsshd/pkg/config"
	"github.com/hiragi-gkuth/bsshd/pkg/ids"
)

func main() {

	// listen server
	listener, e := net.Listen("tcp4", "0.0.0.0:2222")
	if e != nil {
		log.Fatal("failed to listen for connection: ", e)
	}
	log.Print("bsshd is listening on 0.0.0.0:2222 tcp4")
	defer listener.Close()

	// configure bsshd
	bsshdConfig := config.NewServerConfig("assets/keys/host_ecdsa_key")

	// setup sshd child process manager
	procMgr := NewProcManager(64, bsshdConfig)

	// initial kvs
	ids.InitKVS()

	// main server accept loop
	for {
		// wait for connection
		conn, e := listener.Accept()
		if e != nil {
			log.Fatalf("failed to accept tcp connection %v", e)
			continue
		}
		// accept connection
		log.Printf("accept new connection from %v\n", conn.RemoteAddr().String())
		procMgr.AddConn(conn, bsshdConfig)
	}
}
