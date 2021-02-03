package ids

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
)

// BitrisAuthLog は，認証情報を格納する型
type BitrisAuthLog struct {
	f        *fluent.Fluent
	ServerID string
}

// BitrisAuthLogger は，認証情報を管理するインタフェース
type BitrisAuthLogger interface {
	Send(ai AuthInfo)
}

// NewBitrisAuthLogger は，Bitrisのロガーを設定して返す
func NewBitrisAuthLogger(serverID, host string, port int) BitrisAuthLogger {
	config := fluent.Config{
		FluentHost:    host,
		FluentPort:    port,
		FluentNetwork: "tcp",
	}
	f, e := fluent.New(config)
	if e != nil {
		log.Panicln("failed to init fluentd logger:", e)
		return nil
	}

	return &BitrisAuthLog{
		f:        f,
		ServerID: serverID,
	}
}

// Send は，認証情報をログに変換してDBサーバに送信する
func (bal BitrisAuthLog) Send(ai AuthInfo) {
	conn := ai.SSHConnMeta
	var ( // Fluentdに送信するための，認証セッションごとに共通な変数を用意
		sessionID   = hex.EncodeToString(conn.SessionID())
		clientVer   = string(conn.ClientVersion())
		ip, port, _ = net.SplitHostPort(conn.RemoteAddr().String())
		rtt         = ai.InitialTime().Seconds()
	)

	for i := 0; i < ai.AttemptCount; i++ {
		var ( // user, password は，任意の文字列が含まれ，フォーマットが壊される可能性があるため，hexに変換
			userHex     = hex.EncodeToString([]byte(conn.User()))
			passwordHex = hex.EncodeToString([]byte(ai.Passwords[i]))
			result      = ai.Results[i]
			authtime    = ai.AuthTimes()[i].Seconds()
			authAt      = ai.AuthAts[i+1] // 最初のやつは"none"メソッドのタイミングで発生するものなので除外
		)

		// send data
		message := map[string]string{
			"server_id": bal.ServerID,
			"sessionid": sessionID,
			"clientver": clientVer,
			"result":    result,
			"user":      userHex,
			"password":  passwordHex,
			"ip":        ip,
			"port":      port,
			"authtime":  fmt.Sprint(authtime),
			"rtt":       fmt.Sprint(rtt),
			"unixtime":  fmt.Sprint(authAt.Unix()),
			"usec":      fmt.Sprint(authAt.Nanosecond() / 1000),
		}
		hostname, _ := os.Hostname()
		tag := hostname + ".auth.info"

		bal.f.Post(tag, message)
	}
}
