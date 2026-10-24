package cloudconfig

import (
	"io"
	"io/fs"
	"strings"
	"testing"
)

// recordBase records the scripts and commands that are run.
type recordBase struct {
	calls []string
}

func (t *recordBase) SetLogWriter(io.Writer) {}
func (t *recordBase) RunScript(input string) error {
	t.calls = append(t.calls, input)
	return nil
}
func (t *recordBase) RunCommand(args ...string) error {
	t.calls = append(t.calls, strings.Join(args, " "))
	return nil
}
func (t *recordBase) WriteFile(path string, data []byte, perm fs.FileMode) error  { return nil }
func (t *recordBase) AppendFile(path string, data []byte, perm fs.FileMode) error { return nil }
func (t *recordBase) FileExists(path string) (bool, error)                        { return false, nil }

// testOS returns commands that are easy to recognize.
type testOS struct {
	noUpdate bool
}

func (t *testOS) NeedUserPasswords() bool { return false }
func (t *testOS) UpdatePackagesCommand() string {
	if t.noUpdate {
		return ""
	}
	return "update"
}
func (t *testOS) UpgradePackagesCommand() string          { return "upgrade" }
func (t *testOS) InstallPackageCommand(pkg string) string { return "install " + pkg }
func (t *testOS) AddUserCommand(u *User) []string         { return nil }
func (t *testOS) SetTimezoneCommand(tz string) []string   { return nil }

func applyYaml(t *testing.T, os OSType, yaml string) ([]string, error) {
	t.Helper()
	config, err := Unmarshal([]byte(yaml))
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	base := &recordBase{}
	c := NewConfigurer(base)
	c.OS = os
	err = c.Apply(config)
	return base.calls, err
}

func TestPackageUpdateUpgrade(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected string
	}{
		{"none", "packages:\n- a\n", "install a"},
		{"false", "package_update: false\npackage_upgrade: false\npackages:\n- a\n", "install a"},
		{"update", "package_update: true\npackages:\n- a\n", "update,install a"},
		{"upgrade", "package_upgrade: true\npackages:\n- a\n", "upgrade,install a"},
		{"both", "package_update: true\npackage_upgrade: true\npackages:\n- a\n- b\nruncmd:\n- echo\n",
			"update,upgrade,install a,install b,echo"},
		{"no packages", "package_update: true\npackage_upgrade: true\n", "update,upgrade"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls, err := applyYaml(t, &testOS{}, "#cloud-config\n"+test.yaml)
			if err != nil {
				t.Fatal(err)
			}
			if got := strings.Join(calls, ","); got != test.expected {
				t.Fatalf("got %q, expected %q", got, test.expected)
			}
		})
	}
}

func TestPackageUpdateErrors(t *testing.T) {
	_, err := applyYaml(t, nil, "#cloud-config\npackage_update: true\n")
	if err == nil {
		t.Fatalf("expected error without OS")
	}
	_, err = applyYaml(t, &testOS{noUpdate: true}, "#cloud-config\npackage_update: true\n")
	if err == nil {
		t.Fatalf("expected error when update is not supported")
	}
}

func TestFileTags(t *testing.T) {
	config, err := Unmarshal([]byte("#cloud-config\nwrite_files:\n- path: /a\n  owner: root\n  permissions: '0600'\n  content: x\n  append: true\n  defer: true\n"))
	if err != nil {
		t.Fatal(err)
	}
	f := config.Files[0]
	if f.Path != "/a" || f.Owner != "root" || f.Permissions != "0600" || f.Content != "x" || !f.Append || !f.Defer {
		t.Fatalf("unexpected file: %+v", *f)
	}
}
