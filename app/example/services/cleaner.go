package services

import (
	"fmt"
	"github.com/zlbeidou/ago/app"
	"time"
)

func init() {
	app.RegisterRear(func() {
		time.Sleep(time.Second * 2)
		fmt.Println("cleaner")
	})
}
