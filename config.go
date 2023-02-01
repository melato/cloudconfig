package cloudinit

// Comment is the cloud-config comment that indicates that a .yaml file is a cloud-config file.
const Comment = "#cloud-config"

// Config sections are listed in the order in which they are run.
type Config struct {
	Bootcmd        []Command `yaml:"bootcmd,omitempty"`
	PackageUpdate  bool      `yaml:"package_update",omitempty`
	PackageUpgrade bool      `yaml:"package_upgrade",omitempty`
	Packages       []string  `yaml:"packages,omitempty"`
	Files          []*File   `yaml:"write_files,omitempty"`
	Users          []*User   `yaml:"users,omitempty"`
	Timezone       string    `yaml:"timezone,omitempty"`

	// Runcmd is a list of commands to run
	Runcmd []Command `yaml:"runcmd,omitempty"`
}

// Command is string, or []string or []any
// If it is a []string, it is executed using the equivalent of
// execve(3) (with the first arg as the command)
// If it is a []any, it is converted to []string and used as above.
// If it is a string, it is passed as input to /bin/sh
type Command any

type File struct {
	Path        string `yaml:"path`
	Owner       string `yaml:"owner,omitempty`
	Permissions string `yaml:"permissions,omitempty`
	Content     string `yaml:"content`
}

type User struct {
	Name              string   `yaml:"name"`
	Uid               string   `yaml:"uid,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	Homedir           string   `yaml:"homedir,omitempty"`
	NoCreateHome      bool     `yaml:"no_create_home,omitempty"`
	PrimaryGroup      string   `yaml:"primary_group,omitempty"`
	Groups            string   `yaml:"groups,omitempty"`
	Gecos             string   `yaml:"gecos,omitempty"`
	SshAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	/* sudo may be true, false, nil, a string, or a []string
	If it is false or nil, it does nothing
	If the directory /etc/sudoers.d/ exists, a file is created there,
	with the line {user} {sudo-string}, for each string value.
	If sudo is true the file content is set to "{user} ALL=(ALL) NOPASSWD:ALL"

	If the directory /etc/doas.d/ exists, a file is created there,
	with the line {sudo-string} {user}, for each string value.
	If sudo is true, the file content is set to "permit nopass {user}"

	doas and sudo configurations are not compatible, so specifying strings instead of true
	makes sense if only one of the above directories exists.
	*/
	Sudo any `yaml:"sudo,omitempty"`
}
