package local

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"melato.org/cloudinit"
)

type Runner struct {
	OS OS
}

func (t *Runner) runCommandExec(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("command has 0 args")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (t *Runner) runCommandSh(script string) error {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = strings.NewReader(script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (t *Runner) runCommand(command cloudinit.Command) error {
	switch c := command.(type) {
	case string:
		return t.runCommandSh(c)
	case []string:
		return t.runCommandExec(c)
	default:
		return fmt.Errorf("invalid command type: %T", command)
	}
}

func (t *Runner) RunCommands(commands []cloudinit.Command) error {
	for _, command := range commands {
		err := t.runCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}

// may return user.UnknownUserError, user.UnknownGroupError, or other error.
func (t *Runner) findUidGid(owner string) (uid int, gid int, err error) {
	parts := strings.Split(owner, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid owner (user:group): %s", owner)
	}
	var userid = parts[0]
	var groupid = parts[1]

	parse := func(userid, groupid string) (uid int, gid int, err error) {
		// if user, group are numeric, don't look them up
		uid, err = strconv.Atoi(userid)
		if err == nil {
			gid, err = strconv.Atoi(groupid)
		}
		return uid, gid, err
	}

	uid, gid, err = parse(userid, groupid)
	if err != nil {
		return uid, gid, err
	}

	// lookup uid, gid
	var u *user.User
	var g *user.Group
	u, err = user.Lookup(userid)
	if err == nil {
		g, err = user.LookupGroup(groupid)
	}
	if err != nil {
		return 0, 0, err
	}
	return parse(u.Uid, g.Gid)
}

func (t *Runner) WriteFile(file *cloudinit.File) error {
	dir := filepath.Dir(file.Path)
	if dir == "." {
		return fmt.Errorf("relative file path: %s", file.Path)
	}
	err := os.MkdirAll(dir, os.FileMode(0775))
	if err != nil {
		return err
	}

	path := filepath.Clean(file.Path)
	var mode uint64
	if file.Permissions == "" {
		mode = 0o664
	} else {
		mode, err = strconv.ParseUint(file.Permissions, 8, 32)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path, []byte(file.Content), os.FileMode(mode))
	if err != nil {
		return err
	}

	if file.Owner != "" {
		uid, gid, err := t.findUidGid(file.Owner)
		if err == nil {
			err = os.Chown(path, uid, gid)
		} else {
			err = exec.Command("chown", file.Owner, path).Run()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Runner) InstallPackages(packages []string) error {
	if len(packages) == 0 {
		return nil
	}
	if t.OS == nil {
		return fmt.Errorf("cannot install packages.  Missing OS")
	}
	commands := make([]cloudinit.Command, 0, len(packages))
	for _, pkg := range packages {
		commands = append(commands, t.OS.InstallPackageCommand(pkg))
	}
	return t.RunCommands(commands)
}

func (t *Runner) AddUsers(users []*cloudinit.User) error {
	if len(users) == 0 {
		return nil
	}
	if t.OS == nil {
		return fmt.Errorf("cannot create users.  Missing OS")
	}
	commands := make([]cloudinit.Command, 0, len(users))
	for _, user := range users {
		commands = append(commands, t.OS.AddUserCommand(user))
	}
	return t.RunCommands(commands)
}

func (t *Runner) Run(c *cloudinit.Config) error {
	err := t.RunCommands(c.Bootcmd)
	if err != nil {
		return err
	}
	err = t.InstallPackages(c.Packages)
	if err != nil {
		return err
	}
	for _, file := range c.Files {
		err = t.WriteFile(file)
		if err != nil {
			return err
		}
	}
	err = t.AddUsers(c.Users)
	if err != nil {
		return err
	}
	err = t.RunCommands(c.Runcmd)
	if err != nil {
		return err
	}
	return nil
}
