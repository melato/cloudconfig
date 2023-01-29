package cloudinit

import (
	"melato.org/cloudinit/internal"
)

// Config sections are listed in the order in which they are run.
type Config struct {
	Bootcmd  []Command `yaml:"bootcmd,omitempty"`
	Packages []string  `yaml:"packages,omitempty"`
	Files    []*File   `yaml:"write_files,omitempty"`
	// Runcmd a slice whose elements are string or []string
	Users    []*User   `yaml:"users,omitempty"`
	Timezone string    `yaml:"timezone,omitempty"`
	Runcmd   []Command `yaml:"runcmd,omitempty"`
}

// Command is either a []string or a string
// If it is a []string, it is executed using the equivalent of
// execve(3) (with the first arg as the command)
// If it is a string, it is passed as input to /bin/sh
type Command any

type File struct {
	Path        string `yaml:"path`
	Owner       string `yaml:"owner,omitempty`
	Permissions string `yaml:"permissions,omitempty`
	Content     string `yaml:"content`
}

type User struct {
	Name              string   `yaml:"name"`
	Uid               string   `yaml:"uid,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	Groups            []string `yaml:"groups,omitempty"`
	Gecos             string   `yaml:"gecos,omitempty"`
	SshAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
}

// Merge merges two Configs
// It modifies c so as to include the union of c and b.
// Arrays are appended.  Packages are appended and duplicates are removed.
// If a single value is non-empty in c, it stays as is.  Otherwise, it takes the value from b.
func (c *Config) Merge(b *Config) {
	c.Bootcmd = append(c.Bootcmd, b.Bootcmd...)
	packageSet := make(internal.Set[string])
	packages := make([]string, 0, len(c.Packages)+len(b.Packages))
	for _, packageList := range [][]string{c.Packages, b.Packages} {
		for _, pkg := range packageList {
			if !packageSet.Contains(pkg) {
				packageSet.Put(pkg)
				packages = append(packages, pkg)
			}
		}
	}
	c.Packages = packages
	c.Files = append(c.Files, b.Files...)
	c.Users = append(c.Users, b.Users...)
	if c.Timezone == "" {
		c.Timezone = b.Timezone
	}
	c.Runcmd = append(c.Runcmd, b.Runcmd...)
}
