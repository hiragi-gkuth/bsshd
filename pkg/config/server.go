// Package config は，SSHの設定に関する値を管理する
package config

import (
	"github.com/hiragi-gkuth/bsshd/pkg/auth"
	"golang.org/x/crypto/ssh"
)

var common = ssh.Config{}

// NewServerConfig は，bsshdのサーバ設定を行って，返します．
func NewServerConfig(hostKeyPath string, honeypotMode bool) *ssh.ServerConfig {
	var passwordAuthFunc func(ssh.ConnMetadata, []byte) (*ssh.Permissions, error)

	if honeypotMode {
		passwordAuthFunc = auth.PasswordHoneyPot
	} else {
		passwordAuthFunc = auth.Password
	}

	config := &ssh.ServerConfig{
		Config:            common,
		ServerVersion:     "SSH-2.0-OpenSSH_8.2p1",
		MaxAuthTries:      6,
		BannerCallback:    auth.Banner,
		PasswordCallback:  passwordAuthFunc,
		PublicKeyCallback: auth.PublicKey,
		AuthLogCallback:   auth.Log,
	}

	config.AddHostKey(loadHostKey(hostKeyPath))

	return config
}
