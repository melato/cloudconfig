package main

import (
	_ "embed"
	"fmt"

	"melato.org/cloudconfig/cli"
	"melato.org/command"
	"melato.org/command/usage"
)

//go:embed version
var version string

//go:embed usage.yaml
var usageData []byte

func main() {
	cmd := &command.SimpleCommand{}
	var app cli.App
	cmd.Command("apply").Flags(&app).RunFunc(app.Apply)
	cmd.Command("print").Flags(&app).RunFunc(app.Print)
	cmd.Command("parse").RunFunc(cli.Parse)
	cmd.Command("version").RunFunc(func() { fmt.Println(version) })

	usage.Apply(cmd, usageData)
	command.Main(cmd)
}
