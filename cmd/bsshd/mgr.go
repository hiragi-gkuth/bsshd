package main

import (
	"net"

	"golang.org/x/crypto/ssh"
)

// ProcMgr は，bsshdの子プロセスを管理する
type ProcMgr struct {
	MaxConn    int
	SshdConfig *ssh.ServerConfig
}

// NewProcManager は，新しいプロセスマネージャを返す
func NewProcManager(maxConn int, sshdConfig *ssh.ServerConfig) *ProcMgr {
	return &ProcMgr{
		MaxConn:    maxConn,
		SshdConfig: sshdConfig,
	}
}

// AddConn は，新しいssh接続を追加する
func (pm ProcMgr) AddConn(conn net.Conn, config *ssh.ServerConfig) {
	go sshdChild(conn, config)
}

// KillAll は，今持ってる認証プロセスを強制的に全て殺す
func (pm ProcMgr) KillAll() {
}

// Manage は，それぞれのsshdプロセスとの通信を行う
func (pm ProcMgr) Manage() {
}
