package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hiragi-gkuth/bsshd/pkg/config"
	"github.com/hiragi-gkuth/bsshd/pkg/ids"
)

const ( // 解析ソフトウェア側のパラメータ固定してる，許して
	entirePeriod = 24 * time.Hour
	divisions    = 24
	subnetMask   = 16
)

func main() {
	// commandline options
	var (
		// debugMode   = flag.Bool("d", false, "デバッグモードを有効にします")
		hostKeyFile  = flag.String("h", "assets/keys/host_ecdsa_key", "ホストキーを指定します")
		port         = flag.Int("p", 22, "sshdが待機するポートを指定します")
		addr         = flag.String("a", "0.0.0.0", "サーバがバインドするアドレスを指定します")
		logServerID  = flag.String("logID", "bsshd", "fluentに知らせるサーバIDを指定します")
		logHost      = flag.String("logHost", "", "fluentのサーバホストを指定します")
		logPort      = flag.Int("logPort", 24224, "fluentのサーバポートを指定します")
		dbHost       = flag.String("dbHost", "0.0.0.0", "DBサーバホストを指定します")
		dbUser       = flag.String("dbUser", "", "DBユーザを指定します")
		dbPass       = flag.String("dbPass", "", "DBパスワードを指定します")
		honeypotMode = flag.Bool("honeypot", false, "ハニーポットサーバとして起動します")
	)
	flag.Parse()
	// consts
	const MaxConn = 64

	// chans
	signalCh := make(chan os.Signal, 1)
	killCh := make(chan os.Signal, 1)
	term := make(chan os.Signal, 1)
	connCh := make(chan net.Conn, MaxConn)

	// listener
	listener := setupServer(*addr, *port)
	defer listener.Close()

	// handlers
	go sigusr1Handler(signalCh)
	go killHandler(killCh)
	go connHandler(listener, connCh)

	// configure bsshd
	bsshdConfig := config.NewServerConfig(*hostKeyFile, *honeypotMode)

	// setup sshd child process manager
	procMgr := NewProcManager(MaxConn, bsshdConfig, *logServerID, *logHost, *logPort)

	// main server accept loop
	for {
		select {
		case conn := <-connCh:
			procMgr.AddConn(conn)

		case <-signalCh:
			log.Print("IDS model refetch signal received. fetching...")
			signal.Stop(signalCh)
			ids.ReconstructIdsModel(*logServerID, *dbHost, *dbUser, *dbPass)
			signal.Notify(signalCh, syscall.SIGUSR1)
			log.Print("done.")

		case <-killCh:
			log.Print("exit")
			procMgr.KillAll()
			return

		case <-term:
			log.Print("term")
		}
	}
}
