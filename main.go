package main

import (
	"github.com/shibukawa/goqml/goqmllib"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	app          = kingpin.New("goqml", "Go-QML helper")
	buildCommand = app.Command("build", "Build application with icon")
	initCommand  = app.Command("init", "Setup application")
	packCommand  = app.Command("pack", "Create installer")
	setupCommand = app.Command("setup", "Setup go-qml")
)

func main() {
	app.HelpFlag.Short('h')
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case buildCommand.FullCommand():
		goqmllib.Build()
	case packCommand.FullCommand():
		goqmllib.Pack()
	case initCommand.FullCommand():
		goqmllib.Init()
	case setupCommand.FullCommand():
		goqmllib.Setup()
	}
}
