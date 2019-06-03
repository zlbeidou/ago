package app

import (
	"os"
	"syscall"
)

func reloadSig() os.Signal {
	return syscall.SIGTERM
}
