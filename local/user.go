package local

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type Sudo struct {
}

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

func toStrings(a any, trueValue string) ([]string, error) {
	switch v := a.(type) {
	case nil:
		return nil, nil
	case bool:
		if v {
			return []string{trueValue}, nil
		} else {
			return nil, nil
		}
	case string:
		return []string{v}, nil
	case []string:
		return v, nil
	case []any:
		list := make([]string, len(v))
		for i, arg := range v {
			s, isString := arg.(string)
			if !isString {
				break
			}
			list[i] = s
		}
		return list, nil
	}
	return nil, fmt.Errorf("cannot convert to string list: %v", a)
}

func (t *Sudo) IsEnabled() bool {
	return dirExists(SudoersDir)
}

func (t *Sudo) Configure(user string, sudo any) error {
	values, err := toStrings(sudo, "ALL=(ALL) NOPASSWD:ALL")
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	var buf bytes.Buffer
	for _, value := range values {
		fmt.Fprintf(&buf, "%s %s\n", user, value)
	}
	file := filepath.Join(SudoersDir, user)
	return os.WriteFile(file, buf.Bytes(), os.FileMode(0400))
}

type Doas struct {
}

func (t *Doas) IsEnabled() bool {
	return dirExists(DoasDir)
}

func (t *Doas) Configure(user string, sudo any) error {
	values, err := toStrings(sudo, "permit nopass")
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	var buf bytes.Buffer
	for _, value := range values {
		fmt.Fprintf(&buf, "%s %s\n", value, user)
	}
	file := filepath.Join(DoasDir, user)
	return os.WriteFile(file, buf.Bytes(), os.FileMode(0400))
}
