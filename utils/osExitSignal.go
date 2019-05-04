package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func NewChFromOsExitSignal() <-chan os.Signal {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill, syscall.SIGTERM)
	return s
}

func WaitForOsExitSignal() {
	<-NewChFromOsExitSignal()
	return
}
