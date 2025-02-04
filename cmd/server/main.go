package main

import "github.com/teamcubation/go-items-challenge/cmd/server/modules"

func main() {
	app := modules.NewApp()
	// app.Run() starts the application and blocks until it stops.
	app.Run()
}
