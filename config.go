package cloudinit

import (
	"fmt"

	"melato.org/cloudinit/internal"
)

// Config sections are listed in the order in which they are run.
type Config struct {
	Bootcmd  []Command `yaml:"bootcmd,omitempty"`
	Packages []string  `yaml:"packages,omitempty"`
	Files    []*File   `yaml:"write_files,omitempty"`
	Users    []*User   `yaml:"users,omitempty"`
	Timezone string    `yaml:"timezone,omitempty"`

	// Runcmd is a list of commands to run
	Runcmd []Command `yaml:"runcmd,omitempty"`
}

// Command is string, or []string or []any
// If it is a []string, it is executed using the equivalent of
// execve(3) (with the first arg as the command)
// If it is a []any, it is converted to []string and used as above.
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
	Homedir           string   `yaml:"homedir,omitempty"`
	NoCreateHome      bool     `yaml:"no_create_home,omitempty"`
	Groups            []string `yaml:"groups,omitempty"`
	Gecos             string   `yaml:"gecos,omitempty"`
	SshAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	/* sudo may be true, false, nil, a string, or a []string
	If it is false or nil, it does nothing
	If the directory /etc/sudoers.d/ exists, a file is created there,
	with the line {user} {sudo-string}, for each string value.
	If sudo is true the file content is set to "{user} ALL=(ALL) NOPASSWD:ALL"

	If the directory /etc/doas.d/ exists, a file is created there,
	with the line {sudo-string} {user}, for each string value.
	If sudo is true, the file content is set to "permit nopass {user}"

	doas and sudo configurations are not compatible, so specifying strings instead of true
	makes sense if only one of the above directories exists.
	*/
	Sudo any `yaml:"sudo,omitempty"`
}

// Merge Config b into Config c
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

// Script returns the command script content, if the command is a string.
func Script(command Command) (string, bool) {
	s, isString := command.(string)
	return s, isString
}

// Args returns the command args, if the command is a slice.
func Args(command Command) ([]string, bool) {
	switch list := command.(type) {
	case []string:
		return list, true
	case []any:
		args := make([]string, len(list))
		for i, arg := range list {
			switch v := arg.(type) {
			case string:
				args[i] = v
			default:
				args[i] = fmt.Sprintf("%v", arg)
			}
		}
		return args, true
	default:
		return nil, false
	}
}
