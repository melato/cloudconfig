package cloudconfig

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Trace bool

type Configurer struct {
	Base        BaseConfigurer
	OS          OSType
	Log         io.Writer
	createdDirs map[string]struct{}
}

// NewConfigurer creates a Configurer
func NewConfigurer(base BaseConfigurer) *Configurer {
	t := &Configurer{Base: base}
	t.createdDirs = make(map[string]struct{})
	return t
}

func (t *Configurer) logf(format string, args ...any) {
	if t.Log != nil {
		fmt.Fprintf(t.Log, format, args...)
	}
}

func (t *Configurer) ensureDirExists(dir string) error {
	if dir == "/" || dir == "." {
		return nil
	}
	_, exists := t.createdDirs[dir]
	if exists {
		return nil
	}
	err := t.Base.RunCommand("mkdir", "-p", dir)
	if err != nil {
		return err
	}
	for d := dir; !(d == "." || d == "/"); d = filepath.Dir(d) {
		t.createdDirs[d] = struct{}{}
	}
	return nil
}

func (t *Configurer) WriteFile(f *File) error {
	var perm fs.FileMode
	if f.Permissions != "" {
		mode, err := strconv.ParseInt(f.Permissions, 8, 32)
		if err != nil {
			return err
		}
		perm = fs.FileMode(mode)
	} else {
		// does cloud-init specify default permissions?
		perm = fs.FileMode(0644)
	}
	t.logf("write file: %s\n", f.Path)
	dir := filepath.Dir(f.Path)
	err := t.ensureDirExists(dir)
	if err != nil {
		return err
	}
	if f.Append {
		err = t.Base.AppendFile(f.Path, []byte(f.Content), perm)
	} else {
		err = t.Base.WriteFile(f.Path, []byte(f.Content), perm)
	}
	if err != nil {
		return err
	}
	if f.Owner != "" {
		err := t.Base.RunCommand("chown", f.Owner, f.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Configurer) WriteFiles(files []*File, defered bool) error {
	for _, f := range files {
		if f.Defer == defered {
			err := t.WriteFile(f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Configurer) Apply(config *Config) error {
	if t.Base == nil {
		return fmt.Errorf("missing base configurer")
	}
	var err error
	err = t.WriteFiles(config.Files, false)
	if err != nil {
		return err
	}
	err = t.InstallPackages(config.Packages)
	if err != nil {
		return err
	}
	err = t.AddUsers(config.Users)
	if err != nil {
		return err
	}
	if config.Timezone != "" {
		if t.OS == nil {
			return requireOSError("cannot set timezone")
		}
		command := t.OS.SetTimezoneCommand(config.Timezone)
		err := t.RunCommands(Commands{command})
		if err != nil {
			return err
		}
	}
	err = t.WriteFiles(config.Files, true)
	if err != nil {
		return err
	}
	err = t.RunCommands(config.Runcmd)
	if err != nil {
		return err
	}
	return nil
}

// ApplyConfigFiles reads cloud-config files and applies them.
// It reads them all first before applying any of them.
func (t *Configurer) ApplyConfigFiles(files ...string) error {
	configs := make([]*Config, len(files))
	for i, file := range files {
		config, err := ReadFile(file)
		if err != nil {
			return err
		}
		configs[i] = config
	}
	for i, config := range configs {
		err := t.Apply(config)
		if err != nil {
			return fmt.Errorf("%s: %w", files[i], err)
		}
	}
	return nil
}

// Apply reads a cloud-config file from stdin and applies it
func (t *Configurer) ApplyStdin() error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, os.Stdin)
	if err != nil {
		return fmt.Errorf("stdin: %w", err)
	}
	data := buf.Bytes()
	config, err := Unmarshal(data)
	if err != nil {
		return err
	}
	return t.Apply(config)
}

func (t *Configurer) RunCommands(commands Commands) error {
	for _, command := range commands {
		err := t.runCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Configurer) runCommand(command any) error {
	script, isScript := CommandScript(command)
	if isScript {
		t.logf("script << ---\n")
		t.logf("%s\n---\n", script)
		return t.Base.RunScript(script)
	}
	args, isArgs := CommandArgs(command)
	if isArgs {
		if len(args) == 0 {
			return fmt.Errorf("empty command")
		}
		if t.Log != nil {
			t.logf("%s\n", strings.Join(args, " "))
		}
		return t.Base.RunCommand(args...)
	}
	return fmt.Errorf("invalid command type: %T", command)
}

func (t *Configurer) InstallPackages(packages []string) error {
	if len(packages) == 0 {
		return nil
	}
	if t.OS == nil {
		return requireOSError("cannot install packages")
	}
	commands := make(Commands, 0, len(packages))
	for _, pkg := range packages {
		commands = append(commands, t.OS.InstallPackageCommand(pkg))
	}
	return t.RunCommands(commands)
}

func requireOSError(msg string) error {
	return fmt.Errorf("%s.  Missing OS", msg)
}

// ApplySudo default implementation
// It supports sudo and doas.
// It runs scripts that create files /etc/sudoers.d/{username} or /etc/doas.d/{username},
// if these directories exist
// This method is used only if BaseConfigurer does not implement ApplySudo
func (t *Configurer) ApplySudo(username string, values []string) error {
	commands := make(Commands, 0, 2)
	commands = append(commands, sudoScript(username, values))
	commands = append(commands, doasScript(username, values))
	return t.RunCommands(commands)
}

func (t *Configurer) chpasswdScript(pass string, users []string) string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "chpasswd -e << END\n")
	for _, user := range users {
		fmt.Fprintf(&buf, "%s:%s\n", user, pass)
	}
	fmt.Fprintf(&buf, "END\n")
	return buf.String()
}

func (t *Configurer) AddUsers(users []*User) error {
	if len(users) == 0 {
		return nil
	}
	if t.OS == nil {
		return requireOSError("cannot create users")
	}
	commands := make(Commands, 0, len(users))
	for _, u := range users {
		commands = append(commands, t.OS.AddUserCommand(u))
		groups := strings.Split(u.Groups, ",")
		for _, group := range groups {
			group = strings.TrimSpace(group)
			if group != "" {
				commands = append(commands, []string{"adduser", u.Name, group})
			}
		}
	}
	err := t.RunCommands(commands)
	if err != nil {
		return err
	}
	if t.OS.NeedUserPasswords() {
		names := make([]string, len(users))
		for i, u := range users {
			names[i] = u.Name
		}

		script := t.chpasswdScript("*", names)
		err := t.Base.RunScript(script)
		if err != nil {
			return err
		}
	}

	for _, u := range users {
		if u.Sudo != nil && u.Sudo != false {
			values, err := toStrings(u.Sudo)
			if err != nil {
				return fmt.Errorf("invalid sudo value for user %s: %w", u.Name, err)
			}
			applySudo, ok := t.Base.(BaseApplySudo)
			if !ok {
				applySudo = t
			}
			err = applySudo.ApplySudo(u.Name, values)
			if err != nil {
				return err
			}
		}
	}
	for _, u := range users {
		err := t.SetAuthorizedKeys(u)
		if err != nil {
			return err
		}
	}
	return nil
}

// UserHomeDir default implementation
// returns /home/{username} or /root for the root user
func (t *Configurer) UserHomeDir(username string) (string, error) {
	if username == "root" {
		return "/root", nil
	} else {
		return filepath.Join("/home", username), nil
	}
}

func (t *Configurer) SetAuthorizedKeys(u *User) error {
	if len(u.SshAuthorizedKeys) == 0 {
		return nil
	}
	var err error
	homeDir := u.Homedir
	if homeDir == "" {
		impl, ok := t.Base.(BaseUserHomeDir)
		if !ok {
			impl = t
		}
		homeDir, err = impl.UserHomeDir(u.Name)
		if err != nil {
			return err
		}
	}
	dir := filepath.Join(homeDir, ".ssh")
	file := filepath.Join(dir, "authorized_keys")
	exists, err := t.Base.FileExists(file)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	var buf bytes.Buffer
	for _, key := range u.SshAuthorizedKeys {
		fmt.Fprintf(&buf, "%s\n", key)
	}
	err = t.Base.WriteFile(file, buf.Bytes(), os.FileMode(0600))
	if err != nil {
		return err
	}
	owner := u.Name + ":" + u.Name
	for _, command := range [][]string{
		[]string{"chmod", "0755", dir},
		[]string{"chown", owner, dir},
		[]string{"chown", owner, file},
	} {
		if t.Log != nil {
			t.logf("%s\n", strings.Join(command, " "))
		}
		err := t.Base.RunCommand(command...)
		if err != nil {
			return err
		}
	}
	return nil
}
