package services

import (
	"context"
	"fmt"
	"github.com/zlbeidou/ago/app"
	"github.com/zlbeidou/ago/app/example/logger"
	"time"
)

func init() {
	app.RegisterService(&printer{})
}

type printer struct{}

func (p *printer) Init() error {
	return nil
}

func (p *printer) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	fmt.Println("printer working")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stopping printer")
			time.Sleep(time.Second * 2)
			fmt.Println("stopped printer")
			return nil
		case <-ticker.C:
			fmt.Println("printer working")
			logger.LoggerInstance().Log()
		}
	}
}
