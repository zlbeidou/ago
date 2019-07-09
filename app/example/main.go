package main

import (
	"github.com/zlbeidou/ago/app"
	"github.com/zlbeidou/ago/app/example/services"
)

func main() {
	app.RegisterService(&services.Printer{})
	app.RegisterRear(services.Cleaner)
	app.Init()
	app.Run()
}
