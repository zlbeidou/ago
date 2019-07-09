package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

// Service is basic service, only has Init func
type Service interface {
	// Init is used to initialize this service, this func will be called by app.Init()
	Init() error
}

// BackgroundService "inheriting" from Service, has Run func
type BackgroundService interface {
	Service

	// Run should keep running until get ctx.Done
	Run(ctx context.Context) error
}

var (
	services []Service

	rears                []func()
	rearCtx, stopAllRear = context.WithCancel(context.Background())
	rearWg               sync.WaitGroup
)

// RegisterPprofPort register pprof port, port should be greater than 0
func RegisterPprofPort(port int) {
	if port > 0 {
		go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}
}

// RegisterService register service instance, will cache service in list
func RegisterService(instance Service) {
	services = append(services, instance)
}

// RegisterRear register rear function, rear function will be called when program exiting
func RegisterRear(rearFunc func()) {
	rears = append(rears, rearFunc)
}

// Done detect if app is done
func Done() <-chan struct{} {
	return rearCtx.Done()
}

// RearStarted tell app a rear function is started
func RearStarted() {
	rearWg.Add(1)
}

// RearStopped tell app a rear function is stopped
func RearStopped() {
	rearWg.Done()
}

// Init all services by order
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

// Run all backend service at the same time.
// Run will catch exit signal, if get signal, try to done all service via ctx,
// when service exit, execute all registered rear function, exit when all rear function done.
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
		go func(rear func()) {
			rear()
			wg.Done()
		}(rear)
	}
	wg.Wait()
	rearWg.Wait()
}
