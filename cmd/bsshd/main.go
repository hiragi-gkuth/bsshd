package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/hiragi-gkuth/bsshd/pkg/config"
	"github.com/hiragi-gkuth/bsshd/pkg/ids"
)

func main() {
	// commandline options
	var (
		// debugMode   = flag.Bool("d", false, "デバッグモードを有効にします")
		hostKeyFile = flag.String("h", "assets/keys/host_ecdsa_key", "ホストキーを指定します")
		port        = flag.Int("p", 22, "sshdが待機するポートを指定します")
		bindAddr    = flag.String("a", "0.0.0.0", "サーバがバインドするアドレスを指定します")
		logServerID = flag.String("li", "bsshd", "fluentに知らせるサーバIDを指定します")
		logHost     = flag.String("lh", "", "fluentのサーバホストを指定します")
		logPort     = flag.Int("lp", 24224, "fluentのサーバポートを指定します")
		// ipsMode     = flag.Bool("ips", false, "IPS(侵入検知システム)を有効化します")
		// idsMode     = flag.Bool("ids", false, "IDS(侵入防止システム)を有効化します")
	)

	flag.Parse()

	// listen server
	listener, e := net.Listen("tcp", fmt.Sprintf("%s:%d", *bindAddr, *port))
	if e != nil {
		log.Fatal("failed to listen for connection: ", e)
	}
	log.Printf("bsshd is listening on %v", *port)
	defer listener.Close()

	// configure bsshd
	bsshdConfig := config.NewServerConfig(*hostKeyFile)

	// initialize bitris logger
	var logger ids.BitrisAuthLogger = nil
	if *logHost != "" {
		logger = ids.NewBitrisAuthLogger(*logServerID, *logHost, *logPort)
	} else {
		log.Print("Launch bsshd without logging")
	}

	// setup sshd child process manager
	procMgr := NewProcManager(64, bsshdConfig, logger)

	// main server accept loop
	for {
		// wait for connection
		conn, e := listener.Accept()
		if e != nil {
			log.Fatalf("failed to accept tcp connection %v", e)
			continue
		}
		// accept connection
		log.Printf("accept new connection from %v\n", conn.RemoteAddr().String())
		procMgr.AddConn(conn)
	}
}
