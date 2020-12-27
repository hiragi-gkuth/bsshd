package auth

import (
	"time"

	"github.com/hiragi-gkuth/bsshd/pkg/ids"
	"golang.org/x/crypto/ssh"
)

// Banner は，認証前に呼び出される関数
func Banner(conn ssh.ConnMetadata) string {
	key := conn.RemoteAddr().String()
	// Bannar は，接続確立，鍵交換後の次に呼ばれるので，idsのKVSを初期化しておく
	authInfo, _ := ids.AuthSession.Get(key)
	authInfo.AfterEstablishAt = time.Now()
	authInfo.SSHConnMeta = conn
	ids.AuthSession.Set(key, authInfo)
	return ""
}
