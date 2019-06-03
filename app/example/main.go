package main

import (
	"github.com/zlbeidou/ago/app"
	_ "github.com/zlbeidou/ago/app/example/services"
)

func main() {
	app.Init()
	app.Run()
}
