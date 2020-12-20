// Package config は，SSHの設定に関する値を管理する
package config

import (
	"github.com/hiragi-gkuth/bsshd/pkg/auth"
	"golang.org/x/crypto/ssh"
)

var common = ssh.Config{}

// NewServerConfig は，bsshdのサーバ設定を行って，返します．
func NewServerConfig(hostKeyPath string) *ssh.ServerConfig {
	config := &ssh.ServerConfig{
		Config:            common,
		ServerVersion:     "SSH-2.0-OpenSSH_8.2p1",
		MaxAuthTries:      6,
		BannerCallback:    auth.Banner,
		PasswordCallback:  auth.Password,
		PublicKeyCallback: auth.PublicKey,
		AuthLogCallback:   auth.Log,
	}

	config.AddHostKey(loadHostKey(hostKeyPath))

	return config
}
