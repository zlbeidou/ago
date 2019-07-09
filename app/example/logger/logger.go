package logger

import (
	"fmt"
	"github.com/zlbeidou/ago/app"
	"sync"
	"time"
)

type logger struct{}

var (
	loggerInstance     *logger
	onceLoggerInstance sync.Once
)

func (f *logger) Log() {
	// log to file buffer
}

// LoggerInstance get instance of logger
func LoggerInstance() *logger {
	onceLoggerInstance.Do(func() {
		loggerInstance = &logger{}

		// rear do two thing:
		// 1. flush
		// 2. stop goroutine
		app.RearStarted()
		go func() {
			defer app.RearStopped()

			ticker := time.NewTicker(time.Minute)
			for {
				select {
				case <-app.Done():
					time.Sleep(time.Second * 2)
					fmt.Println("log flush")
					// flush log file buffer
					return
				case <-ticker.C:
					// flush log file buffer
				}
			}
		}()
	})

	return loggerInstance
}
