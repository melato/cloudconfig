package local

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type BaseConfigurer struct {
	Log io.Writer
}

func (t *BaseConfigurer) SetLogWriter(w io.Writer) {
	t.Log = w
}

func (t *BaseConfigurer) RunCommand(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("command has 0 args")
	}
	fmt.Printf("log\n: %x\n", t.Log)
	cmd := exec.Command(args[0], args[1:]...)
	fmt.Printf("%s\n", cmd.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = t.Log
	cmd.Stderr = t.Log
	return cmd.Run()
}

func (t *BaseConfigurer) RunScript(script string) error {
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = strings.NewReader(script)
	cmd.Stdout = t.Log
	cmd.Stderr = t.Log
	return cmd.Run()
}

// may return user.UnknownUserError, user.UnknownGroupError, or other error.
func (t *BaseConfigurer) findUidGid(owner string) (uid int, gid int, err error) {
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

func (t *BaseConfigurer) ensureDirExists(path string) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		err := os.MkdirAll(dir, fs.FileMode(0775))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *BaseConfigurer) WriteFile(path string, data []byte, perm fs.FileMode) error {
	err := t.ensureDirExists(path)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

func (t *BaseConfigurer) AppendFile(path string, data []byte, perm fs.FileMode) error {
	err := t.ensureDirExists(path)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return f.Close()
}

func (t *BaseConfigurer) FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// UserHomeDir default implementation
// returns /home/{username} or /root for the root user
func (t *BaseConfigurer) UserHomeDir(username string) (string, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}
