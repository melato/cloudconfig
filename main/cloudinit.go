package main

import (
	_ "embed"
	"fmt"

	"melato.org/cloudinit"
	"melato.org/cloudinit/local"
	"melato.org/cloudinit/ostype"
	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/yaml"
)

//go:embed version
var version string

//go:embed usage.yaml
var usageData []byte

type Run struct {
	ConfigFile string `name:"f" usage:"cloud-config yaml file"`
	OS         string
	os         cloudinit.OSType
}

func (t *Run) Configured() error {
	if t.ConfigFile == "" {
		return fmt.Errorf("missing config file")
	}
	var config *cloudinit.Config
	err := yaml.ReadFile(t.ConfigFile, &config)
	if err != nil {
		return err
	}
	switch t.OS {
	case "":
	case "alpine":
		t.os = &ostype.Alpine{}
	case "debian":
		t.os = &ostype.Debian{}
	default:
		return fmt.Errorf("unrecognized OS.  accepted values are alpine, debian")
	}
	return nil
}

func (t *Run) Run(configFiles ...string) error {
	configurer := cloudinit.NewConfigurer(&local.BaseConfigurer{})
	configurer.OS = t.os
	return configurer.ApplyConfigFiles(configFiles...)
}

func (t *Run) Print(file string) error {
	var config *cloudinit.Config
	err := yaml.ReadFile(file, &config)
	if err != nil {
		return err
	}
	return yaml.Print(config)
}

func Parse(files []string) error {
	for _, file := range files {
		var config *cloudinit.Config
		err := yaml.ReadFile(file, &config)
		if err != nil {
			fmt.Printf("%s ERROR\n", file)
			return err
		}
		fmt.Printf("%s OK\n", file)
	}
	return nil
}

func main() {
	cmd := &command.SimpleCommand{}
	var app Run
	cmd.Command("run").Flags(&app).RunFunc(app.Run)
	cmd.Command("print").Flags(&app).RunFunc(app.Print)
	cmd.Command("parse").RunFunc(Parse)
	cmd.Command("version").RunFunc(func() { fmt.Println(version) })

	usage.Apply(cmd, usageData)
	command.Main(cmd)
}
