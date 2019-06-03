package app

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

type Service interface {
	Init() error
}

type BackgroundService interface {
	Service
	Run(ctx context.Context) error
}

var (
	services []Service

	rears                []func()
	rearCtx, stopAllRear = context.WithCancel(context.Background())
	rearWg               sync.WaitGroup
)

func RegisterService(instance Service) {
	services = append(services, instance)
}

func RegisterRear(rearFunc func()) {
	rears = append(rears, rearFunc)
}

func Done() <-chan struct{} {
	return rearCtx.Done()
}

func RearStarted() {
	rearWg.Add(1)
}

func RearStopped() {
	rearWg.Done()
}

func Init() {
	for _, service := range services {
		err := service.Init()
		if err != nil {
			log.Println("[CRIC] service Init fail:", err)
			os.Exit(1)
		}
	}
}

func reload() {
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Start()
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() { // catch exit signal
		s := make(chan os.Signal, 1)
		signal.Notify(s)

		for {
			sig := <-s
			if sig == reloadSig() { // need to reload
				signal.Stop(s)
				reload()
				cancel()
			} else {
				switch sig {
				case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
					signal.Stop(s)
					cancel()
				}
			}
		}
	}()

	var wg sync.WaitGroup
	for _, service := range services {
		if backgroundService, ok := service.(BackgroundService); ok {
			wg.Add(1)
			go func() {
				err := backgroundService.Run(ctx)
				if err != nil {
					log.Println("[ERRO] backgroundService Run err:", err)
				}

				wg.Done()
			}()
		}
	}
	wg.Wait()

	stopAllRear()
	for _, rear := range rears {
		wg.Add(1)
		go func() {
			rear()
			wg.Done()
		}()
	}
	wg.Wait()
	rearWg.Wait()
}
