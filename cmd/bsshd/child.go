package main

import (
	"encoding/hex"
	"log"
	"net"
	"time"

	"github.com/hiragi-gkuth/bsshd/pkg/channel"
	"github.com/hiragi-gkuth/bsshd/pkg/ids"
	"golang.org/x/crypto/ssh"
)

func sshdChild(conn net.Conn, config *ssh.ServerConfig, logger ids.BitrisAuthLogger) {
	authSession := ids.GetAuthSession()
	key := conn.RemoteAddr().String()
	// RTT計測のため，SSHコネクション確立前に，時間を保存しておく
	authInfo := ids.NewAuthInfo()
	authInfo.BeforeEstablishAt = time.Now()
	authSession.Set(key, authInfo)
	// SSHセッションの確立を試みる
	sshConn, chans, reqs, e := ssh.NewServerConn(conn, config)
	authInfo, _ = authSession.Get(key)
	authInfo.ShowLogs()
	logger.Send(authInfo)
	if e != nil { // 失敗したら終了
		log.Println("establish failed: ", e)
		return
	}

	log.Printf("established connection: %v, %v\n", sshConn.User(), hex.EncodeToString(sshConn.SessionID()))

	go ssh.DiscardRequests(reqs)
	// Service the incoming Channel channel.
	channel.Handler(chans)
}
