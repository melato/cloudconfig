package cloudinit

type CloudConfig struct {
	Timezone string   `yaml:"timezone,omitempty"`
	Packages []string `yaml:"packages,omitempty"`
	Files    []*File  `yaml:"write-files,omitempty"`
	Runcmd   any      `yaml:"runcmd,omitempty"`
	Users    []*User  `yaml:"users,omitempty"`
}

type User struct {
	Name              string   `yaml:"name"`
	Uid               string   `yaml:"uid,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	Groups            []string `yaml:"groups,omitempty"`
	SshAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
}

type File struct {
	Path        string `yaml:"path`
	Owner       string `yaml:"owner,omitempty`
	Permissions string `yaml:"permissions,omitempty`
	Content     string `yaml:"content`
}
