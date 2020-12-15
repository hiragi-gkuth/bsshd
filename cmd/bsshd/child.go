package main

import (
	"encoding/hex"
	"log"
	"net"

	"github.com/hiragi-gkuth/bsshd/pkg/channel"
	"github.com/hiragi-gkuth/bsshd/pkg/ids"
	"golang.org/x/crypto/ssh"
)

func sshdChild(conn net.Conn, config *ssh.ServerConfig) {
	sessionKey := conn.RemoteAddr().String()
	sshConn, chans, reqs, e := ssh.NewServerConn(conn, config)

	authInfo := ids.KVS[sessionKey]
	authInfo.ShowLogs()
	if e != nil {
		log.Println("establish failed: ", e)
		return
	}

	log.Printf("established connection: %v, %v\n", sshConn.User(), hex.EncodeToString(sshConn.SessionID()))

	go ssh.DiscardRequests(reqs)
	// Service the incoming Channel channel.
	channel.Handler(chans)
}
