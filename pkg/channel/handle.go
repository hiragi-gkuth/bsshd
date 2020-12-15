package channel

import (
	"log"

	"golang.org/x/crypto/ssh"
)

// Handler は，SSHのchannelを制御する
func Handler(chans <-chan ssh.NewChannel) {
	for newChan := range chans {
		switch t := newChan.ChannelType(); t {
		case "session":
			c, _, e := newChan.Accept()
			if e != nil {
				log.Fatalf("Could not accept channel (%s)", e)
				return
			}
			handleSession(c)
		default:
			newChan.Reject(ssh.UnknownChannelType, "Unknown channel type")
		}
	}
}
