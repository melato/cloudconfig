package os

import (
	"melato.org/cloudinit"
)

type Debian struct {
}

// NeedPasswords returns false for Debian, because we can disable
// the account (therefore disabling password login), and still  use it
// for passwordless ssh login.
func (t *Debian) NeedUserPasswords() bool { return false }

func (t *Debian) InstallPackageCommand(pkg string) string {
	return "DEBIAN_FRONTEND=noninteractive apt-get -y install " + pkg
}

func (t *Debian) AddUserCommand(u *cloudinit.User) []string {
	args := []string{"adduser", u.Name, "--disabled-password", "--gecos", u.Gecos}
	if u.Uid != "" {
		args = append(args, "--uid", u.Uid)
	}
	if u.Shell != "" {
		args = append(args, "--shell", u.Shell)
	}
	if u.NoCreateHome {
		args = append(args, "--no-create-home")
	} else if u.Homedir != "" {
		args = append(args, "--home", u.Homedir)
	}
	if u.PrimaryGroup != "" {
		args = append(args, "--group", u.PrimaryGroup)
	}
	return args
}

func (t *Debian) SetTimezoneCommand(timezone string) []string {
	return []string{"timedatectl", "set-timezone", timezone}
}
