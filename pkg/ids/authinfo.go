// Package ids は，認証試行を遮断する仕組みを提供する
package ids

import (
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// AuthInfo は，認証情報を含む
type AuthInfo struct {
	SSHConnMeta       ssh.ConnMetadata
	AttemptCount      int
	BeforeEstablishAt time.Time
	AfterEstablishAt  time.Time
	AuthAts           []time.Time
	Passwords         []string
	Results           []string
}

// NewAuthInfo は，新しい認証情報を返す
func NewAuthInfo() AuthInfo {
	return AuthInfo{
		SSHConnMeta:  nil,
		AttemptCount: 0,
		AuthAts:      []time.Time{},
		Passwords:    []string{},
		Results:      []string{},
	}
}

// AuthTimes は，認証試行の結果から，認証時間のリストを返す．
func (ai AuthInfo) AuthTimes() []time.Duration {
	authTimes := make([]time.Duration, 0)
	for i := 0; i < len(ai.AuthAts)-1; i++ {
		authTime := ai.AuthAts[i+1].Sub(ai.AuthAts[i])
		authTimes = append(authTimes, authTime)
	}
	return authTimes
}

// InitialTime は，TCPコネクション確立後，認証開始前までの時間を計算して返す
func (ai AuthInfo) InitialTime() time.Duration {
	d := ai.AfterEstablishAt.Sub(ai.BeforeEstablishAt)

	d, _ = time.ParseDuration(fmt.Sprintf("%vns", d.Nanoseconds()/5))
	// SSHの接続確立開始から，Bannarまでだいたい5往復ぐらいする．
	// (計測開始)
	// 1. Notify ServerVersion -> Notify ClientVersion
	// 2. -> Server: KexInit    -> Client: KexInit
	// 3. -> Client: DH KexInit -> Server: DH KexInit
	// 4. -> Server: NewKeys    -> Client: NewKeys (ここまでパケットキャプチャで確認)
	// 5. -> Client: `ssh-userauth` Request -> Server: `SSH2_MSG_SERVICE_ACCEPT` Acception
	// 6. -> Server: `input_userauth_banner` (このタイミングで計測終了)

	// 3. のタイミングで，クライアントが2連続でパケットを投げる．これはTCPのACKを待ってるはず．
	// TCPなので，全ての通信ごとにACKが投げられるし，それを待機するかはWindowSizeによる．
	// 2. のタイミングでは，クライアント側実装によって，あるいはMTUによって，パケットが断片化する可能性が高い事を留意
	// TCPにおけるサーバ側からのRTT計測は困難であることを留意．
	return d
}

// ShowLogs は，AuthInfoをもとにログを出力する
func (ai AuthInfo) ShowLogs() {
	fmt.Printf(
		`ID: %v
Ver: %v
IP: %v
User: %v
Pass: %v
AuthAts: %v
Results: %v
AuthTimes: %v
InitialTime: %s
`,
		hex.EncodeToString(ai.SSHConnMeta.SessionID()),
		string(ai.SSHConnMeta.ClientVersion()),
		ai.SSHConnMeta.RemoteAddr().String(),
		ai.SSHConnMeta.User(),
		ai.Passwords,
		ai.AuthAts[len(ai.AuthAts)-1],
		ai.Results,
		ai.AuthTimes(),
		ai.InitialTime(),
	)
}
