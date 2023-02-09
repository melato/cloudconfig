package main

import (
	_ "embed"
	"fmt"
	"os"

	"melato.org/cloudconfig"
	"melato.org/cloudconfig/local"
	"melato.org/cloudconfig/ostype"
	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/yaml"
)

//go:embed version
var version string

//go:embed usage.yaml
var usageData []byte

type Run struct {
	OS string
	os cloudconfig.OSType
}

func (t *Run) Configured() error {
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

func (t *Run) Apply(configFiles ...string) error {
	configurer := cloudconfig.NewConfigurer(&local.BaseConfigurer{})
	configurer.OS = t.os
	configurer.Log = os.Stdout
	if len(configFiles) == 1 && configFiles[0] == "-" {
		return configurer.ApplyStdin()
	} else {
		return configurer.ApplyConfigFiles(configFiles...)
	}
}

func (t *Run) Print(file string) error {
	var config *cloudconfig.Config
	err := yaml.ReadFile(file, &config)
	if err != nil {
		return err
	}
	return yaml.Print(config)
}

func Parse(files []string) error {
	for _, file := range files {
		var config *cloudconfig.Config
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
	cmd.Command("apply").Flags(&app).RunFunc(app.Apply)
	cmd.Command("print").Flags(&app).RunFunc(app.Print)
	cmd.Command("parse").RunFunc(Parse)
	cmd.Command("version").RunFunc(func() { fmt.Println(version) })

	usage.Apply(cmd, usageData)
	command.Main(cmd)
}
