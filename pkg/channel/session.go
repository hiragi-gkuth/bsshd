// Package channel は，SSHのchannelを扱う
package channel

import (
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
)

// HandleSession は，SSHのSessionを扱う
func handleSession(sshChannel ssh.Channel) {
	zsh := exec.Command("zsh")

	log.Print("shell executed")

	close := func() {
		sshChannel.Close()
		_, err := zsh.Process.Wait()
		if err != nil {
			log.Printf("Failed to exit bash (%s)", err)
		}
		log.Printf("Session closed")
	}

	p, e := pty.Start(zsh)
	if e != nil {
		log.Fatalln("Could not start pty", e)
		close()
		return
	}
	var once sync.Once
	go func() {
		io.Copy(sshChannel, p)
		once.Do(close)
	}()
	go func() {
		io.Copy(p, sshChannel)
		once.Do(close)
	}()
}
