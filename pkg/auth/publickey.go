package auth

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)
// PublickeyAuthenticationError は，公開鍵認証が失敗したときに投げられるエラー
type PublickeyAuthenticationError struct {
	user     string
	message  string
}


// PublicKey は，公開鍵認証に関わる機能を提供する
func PublicKey(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	return nil, PublickeyAuthenticationError{}
}

func (rcv PublickeyAuthenticationError) Error() string {
	return fmt.Sprintf("Failed login for %v using %v", rcv.user, rcv.message)
}
