package cloudinit

import (
	"fmt"
)

func FirstLineIs(data []byte, line string) bool {
	n := len(line)
	if len(data) < n {
		return false
	}
	if len(data) > n {
		c := rune(data[n])
		if (c != '\r') && (c != '\n') {
			return false
		}
	}
	for i := 0; i < n; i++ {
		if data[i] != line[i] {
			return false
		}
	}
	return true
}

// HasComment returns true if the first line in the provided data is Comment.
func HasComment(data []byte) bool {
	return FirstLineIs(data, Comment)
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

func toStrings(a any) ([]string, error) {
	switch v := a.(type) {
	case bool:
		if v {
			return nil, nil
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
