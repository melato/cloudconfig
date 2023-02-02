package cloudinit

import (
	"bytes"
	"fmt"
	"path/filepath"
)

const (
	SudoersDir = "/etc/sudoers.d"
	DoasDir    = "/etc/doas.d"
)

func sudoScript(user string, values []string) string {
	if len(values) == 0 {
		values = []string{"ALL=(ALL) NOPASSWD:ALL"}
	}
	var buf bytes.Buffer
	file := filepath.Join(SudoersDir, user)
	fmt.Fprintf(&buf, "if [ -d %s ]; then\ncat <<END > %s\n", SudoersDir, file)
	for _, value := range values {
		fmt.Fprintf(&buf, "%s %s\n", user, value)
	}
	fmt.Fprintf(&buf, "END\n")
	fmt.Fprintf(&buf, "chmod 600 %s\n", file)
	fmt.Fprintf(&buf, "fi\n")
	return buf.String()
}

func doasScript(user string, values []string) string {
	if len(values) == 0 {
		values = []string{"permit nopass"}
	}
	var buf bytes.Buffer
	file := fmt.Sprintf("/etc/doas.d/%s.conf", user)
	fmt.Fprintf(&buf, "if [ -d %s ]; then\ncat <<END > %s\n", DoasDir, file)
	for _, value := range values {
		fmt.Fprintf(&buf, "%s %s\n", value, user)
	}
	fmt.Fprintf(&buf, "END\n")
	fmt.Fprintf(&buf, "chmod 600 %s\n", file)
	fmt.Fprintf(&buf, "fi\n")
	return buf.String()
}
