package cli

import (
	_ "embed"
	"fmt"
	"os"

	"melato.org/cloudconfig"
	"melato.org/cloudconfig/local"
	"melato.org/cloudconfig/ostype"
)

type App struct {
	OS string
	os cloudconfig.OSType
}

func (t *App) Configured() error {
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

func (t *App) Apply(configFiles ...string) error {
	base := &local.BaseConfigurer{}
	base.SetLogWriter(os.Stdout)
	configurer := cloudconfig.NewConfigurer(base)
	configurer.OS = t.os
	configurer.Log = os.Stdout
	if len(configFiles) == 1 && configFiles[0] == "-" {
		return configurer.ApplyStdin()
	} else {
		return configurer.ApplyConfigFiles(configFiles...)
	}
}

func Print(file string) error {
	config, err := cloudconfig.ReadFile(file)
	if err != nil {
		return err
	}
	data, err := cloudconfig.Marshal(config)
	if err != nil {
		return err
	}
	os.Stdout.Write(data)
	return nil
}

func Parse(files []string) error {
	for _, file := range files {
		_, err := cloudconfig.ReadFile(file)
		if err != nil {
			fmt.Printf("%s ERROR\n", file)
			return err
		}
		fmt.Printf("%s OK\n", file)
	}
	return nil
}

func Packages(files []string) error {
	for _, file := range files {
		c, err := cloudconfig.ReadFile(file)
		if err != nil {
			return fmt.Errorf("%s %e\n", file, err)
		}
		for _, p := range c.Packages {
			fmt.Printf("%s\n", p)
		}
	}
	return nil
}
