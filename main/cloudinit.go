package main

import (
	_ "embed"
	"fmt"

	"melato.org/cloudinit"
	"melato.org/cloudinit/local"
	"melato.org/cloudinit/local/os"
	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/yaml"
)

//go:embed version
var version string

//go:embed usage.yaml
var usageData []byte

type App struct {
	ConfigFile string `name:"f" usage:"cloud-config yaml file"`
	OS         string
	config     *cloudinit.Config
	os         local.OS
}

func (t *App) Configured() error {
	if t.ConfigFile == "" {
		return fmt.Errorf("missing config file")
	}
	var config *cloudinit.Config
	err := yaml.ReadFile(t.ConfigFile, &config)
	if err != nil {
		return err
	}
	t.config = config
	switch t.OS {
	case "":
	case "alpine":
		t.os = &os.Alpine{}
	case "debian":
		t.os = &os.Debian{}
	default:
		return fmt.Errorf("unrecognized OS.  accepted values are alpine, debian")
	}
	return nil
}

func (t *App) Run() error {
	runner := &local.Runner{OS: t.os}
	return runner.Run(t.config)
}

func main() {
	cmd := &command.SimpleCommand{}
	var app App
	cmd.Command("run").Flags(&app).RunFunc(app.Run)
	cmd.Command("version").RunFunc(func() { fmt.Println(version) })

	usage.Apply(cmd, usageData)
	command.Main(cmd)
}
