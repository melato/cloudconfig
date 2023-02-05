package ostype

import (
	"melato.org/cloudconfig"
)

type Alpine struct {
}

// NeedPasswords returns true for alpine.
// In order to use passwordless authentication with alpinelinux,
// the account must be enabled (unlike other distributions, which allow disabled
// accounts to be used with ssh, but not login with password).
// Since we keep the account enabled, we require it to have a password.
func (t *Alpine) NeedUserPasswords() bool { return true }

func (t *Alpine) InstallPackageCommand(pkg string) string {
	return "apk add " + pkg
}

func (t *Alpine) AddUserCommand(u *cloudconfig.User) []string {
	args := []string{"adduser", "-g", u.Gecos, "-D"}
	if u.Uid != "" {
		args = append(args, "-u", u.Uid)
	}
	if u.Shell != "" {
		args = append(args, "-s", u.Shell)
	}
	if u.NoCreateHome {
		args = append(args, "-H")
	} else if u.Homedir != "" {
		args = append(args, "-h", u.Homedir)
	}
	if u.PrimaryGroup != "" {
		args = append(args, "-G", u.PrimaryGroup)
	}
	args = append(args, u.Name)
	return args
}

func (t *Alpine) SetTimezoneCommand(timezone string) []string {
	return []string{"setup-timezone", "-z", timezone}
}
