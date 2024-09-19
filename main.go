package main

import (
	"os"

	"github.com/supermicah/go-framework-admin/cmd"
	"github.com/urfave/cli/v2"
)

// Usage: go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "v1.0.0"

// @title goframeworkadmin
// @version v1.0.0
// @description An API service based on golang.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schemes http https
// @basePath /
func main() {
	app := cli.NewApp()
	app.Name = "goframeworkadmin"
	app.Version = VERSION
	app.Usage = "An API service based on golang."
	app.Commands = []*cli.Command{
		cmd.StartCmd(),
		cmd.StopCmd(),
		cmd.VersionCmd(VERSION),
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
