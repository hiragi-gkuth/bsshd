package config

import (
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

func loadHostKey(filename string) ssh.Signer {
	hostkeyBytes, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Fatal("failed to load host key", e)
	}

	hostkey, e := ssh.ParsePrivateKey(hostkeyBytes)
	if e != nil {
		log.Fatal("failed to parse host key", e)
	}

	return hostkey
}
