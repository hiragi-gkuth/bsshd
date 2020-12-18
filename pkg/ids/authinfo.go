// Package ids は，認証試行を遮断する仕組みを提供する
package ids

import (
	"encoding/hex"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

// AuthInfo は，認証情報を含む
type AuthInfo struct {
	SSHConnMeta  ssh.ConnMetadata
	AttemptCount int
	AuthAts      []*time.Time
	Passwords    []string
	Results      []string
	RTT          *time.Duration
}

// AuthTimes は，認証試行の結果から，認証時間のリストを返す．
func (ai AuthInfo) AuthTimes() []*time.Duration {
	authTimes := make([]*time.Duration, 0)
	for i := 0; i < len(ai.AuthAts)-1; i++ {
		authTime := ai.AuthAts[i+1].Sub(*ai.AuthAts[i])
		authTimes = append(authTimes, &authTime)
	}
	return authTimes
}

// ShowLogs は，AuthInfoをもとにログを出力する
func (ai AuthInfo) ShowLogs() {
	log.Printf(
		`SessionID: %v
		ClientVer: %v
		ServerVer: %v
		IP: %v
		User: %v
		Passwords: %v
		AuthAts: %v
		Results: %v
		AuthTimes: %v`,
		hex.EncodeToString(ai.SSHConnMeta.SessionID()),
		string(ai.SSHConnMeta.ClientVersion()),
		string(ai.SSHConnMeta.ServerVersion()),
		ai.SSHConnMeta.RemoteAddr().String(),
		ai.SSHConnMeta.User(),
		ai.Passwords,
		ai.AuthAts[len(ai.AuthAts)-1],
		ai.Results,
		ai.AuthTimes(),
	)
}
