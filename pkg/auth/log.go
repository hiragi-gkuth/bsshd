package auth

import (
	"log"
	"time"

	"github.com/hiragi-gkuth/bsshd/pkg/ids"
	"golang.org/x/crypto/ssh"
)

// Log は，認証試行が行われたときに呼ばれる関数
func Log(conn ssh.ConnMetadata, method string, e error) {
	log.Println("Log callback")
	log.Printf("Login attempt by %v using %s\n", conn.User(), method)

	key := conn.RemoteAddr().String()
	switch method {
	case "none": // initial access
		if authInfo, ok := ids.KVS[key]; !ok {
			log.Fatal("ids session kvs must not be nil.")
		} else {
			now := time.Now()
			authInfo.AuthAts = append(authInfo.AuthAts, &now)
			authInfo.SSHConnMeta = conn
			ids.KVS[key] = authInfo
		}
	case "password": // when password
		now := time.Now()
		authInfo := ids.KVS[key]
		authInfo.AuthAts = append(authInfo.AuthAts, &now)
		authInfo.SSHConnMeta = conn
		ids.KVS[key] = authInfo
	}

	if e != nil {
		log.Println("  failed for reason: ", e)
	}
}
