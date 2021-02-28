package main

import (
	"os"
	"os/signal"
	"syscall"
)

func sigusr1Handler(sigs chan os.Signal) {
	signal.Notify(sigs, syscall.SIGUSR1)
}

func killHandler(sigs chan os.Signal) {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
}
