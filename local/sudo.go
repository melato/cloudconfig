package local

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

const (
	SudoersDir = "/etc/sudoers.d"
	DoasDir    = "/etc/doas.d"
)

func dirExists(dir string) bool {
	st, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return st.IsDir()
}

func (t *BaseConfigurer) applySudo(user string, values []string) error {
	if len(values) == 0 {
		values = []string{"ALL=(ALL) NOPASSWD:ALL"}
	}
	var buf bytes.Buffer
	for _, value := range values {
		fmt.Fprintf(&buf, "%s %s\n", user, value)
	}
	file := filepath.Join(SudoersDir, user)
	return os.WriteFile(file, buf.Bytes(), os.FileMode(0400))
}

func (t *BaseConfigurer) applyDoas(user string, values []string) error {
	if len(values) == 0 {
		values = []string{"permit nopass"}
	}
	var buf bytes.Buffer
	for _, value := range values {
		fmt.Fprintf(&buf, "%s %s\n", value, user)
	}
	file := filepath.Join(DoasDir, user)
	return os.WriteFile(file, buf.Bytes(), os.FileMode(0400))
}

func (t *BaseConfigurer) ApplySudo(username string, values []string) error {
	hasDoas := dirExists(DoasDir)
	hasSudo := dirExists(SudoersDir)
	if hasDoas && hasSudo && len(values) > 0 {
		return fmt.Errorf("cannot apply the same configuration to both doas and sudo")
	}
	if hasDoas {
		err := t.applyDoas(username, values)
		if err != nil {
			return err
		}
	}
	if hasSudo {
		err := t.applySudo(username, values)
		if err != nil {
			return err
		}
	}
	return nil
}
