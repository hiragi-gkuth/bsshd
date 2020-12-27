package auth

import (
	"log"
	"time"

	"github.com/hiragi-gkuth/bsshd/pkg/ids"
	"golang.org/x/crypto/ssh"
)

// Log は，認証試行が行われたときに呼ばれる関数
func Log(conn ssh.ConnMetadata, method string, e error) {
	log.Printf("Login attempt by %v using %s\n", conn.User(), method)

	key := conn.RemoteAddr().String()
	switch method {
	case "none": // initial access
		if authInfo, ok := ids.AuthSession.Get(key); !ok {
			log.Fatal("ids session kvs must not be nil.")
		} else {
			authInfo.AuthAts = append(authInfo.AuthAts, time.Now())
			authInfo.SSHConnMeta = conn
			ids.AuthSession.Set(key, authInfo)
		}
	case "password": // when password
		authInfo, _ := ids.AuthSession.Get(key)
		authInfo.AuthAts = append(authInfo.AuthAts, time.Now())
		authInfo.AttemptCount++
		authInfo.SSHConnMeta = conn
		ids.AuthSession.Set(key, authInfo)
	}

	if e != nil {
		log.Print("  failed for reason: ", e)
	}
}
