package main

import (
	"log"
	"net"

	"github.com/hiragi-gkuth/bsshd/pkg/ids"
	"golang.org/x/crypto/ssh"
)

// ProcMgr は，bsshdの子プロセスを管理する
type ProcMgr struct {
	MaxConn    int
	SshdConfig *ssh.ServerConfig
	Logger     ids.BitrisAuthLogger
}

// NewProcManager は，新しいプロセスマネージャを返す
func NewProcManager(maxConn int, sshdConfig *ssh.ServerConfig, serverID string, logHost string, logPort int) *ProcMgr {
	// initialize bitris logger
	var logger ids.BitrisAuthLogger = nil
	if logHost != "" {
		logger = ids.NewBitrisAuthLogger(serverID, logHost, logPort)
	} else {
		log.Print("Launch bsshd without logging")
	}
	return &ProcMgr{
		MaxConn:    maxConn,
		SshdConfig: sshdConfig,
		Logger:     logger,
	}
}

// AddConn は，新しいssh接続を追加する
func (pm ProcMgr) AddConn(conn net.Conn) {
	go sshdChild(conn, pm.SshdConfig, pm.Logger)
}

// KillAll は，今持ってる認証プロセスを強制的に全て殺す
func (pm ProcMgr) KillAll() {
}

// Manage は，それぞれのsshdプロセスとの通信を行う
func (pm ProcMgr) Manage() {
}
