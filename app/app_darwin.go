package app

import (
	"os"
	"syscall"
)

func reloadSig() os.Signal {
	return syscall.SIGUSR2
}
