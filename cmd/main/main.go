package main

import (
	"github.com/kgugunava/effective_mobile_golang/internal/app"
)


func main() {
	app := app.NewApp()
	app.Router.Run(app.Cfg.ServerAddress)
}