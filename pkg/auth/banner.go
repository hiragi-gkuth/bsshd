package auth

import (
	"log"
	"time"

	"github.com/hiragi-gkuth/bsshd/pkg/ids"
	"golang.org/x/crypto/ssh"
)

// Banner は，認証前に呼び出される関数
func Banner(conn ssh.ConnMetadata) string {
	log.Println("Banner callback")
	// Bannar は，接続確立，鍵交換後の次に呼ばれるので，idsのKVSを初期化しておく
	ids.KVS[conn.RemoteAddr().String()] = ids.AuthInfo{
		SSHConnMeta: conn,
		AuthAts:     []*time.Time{},
		Success:     false,
		RTT:         nil,
	}
	return "banner"
}
