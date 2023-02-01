package cloudinit

import (
	"bufio"
	"bytes"
	"fmt"

	"melato.org/cloudinit/internal"
)

// FirstLine returns the first line in the data
func FirstLine(data []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

// HasComment returns true if the first line in the provided data is Comment.
func HasComment(data []byte) bool {
	return FirstLine(data) == Comment
}

// Script returns the command script content, if the command is a string.
func CommandScript(command Command) (string, bool) {
	s, isString := command.(string)
	return s, isString
}

// Args returns the command args, if the command is a slice.
func CommandArgs(command Command) ([]string, bool) {
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
