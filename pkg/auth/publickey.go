package auth

import "golang.org/x/crypto/ssh"

// PublicKey は，公開鍵認証に関わる機能を提供する
func PublicKey(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	return nil, nil
}
