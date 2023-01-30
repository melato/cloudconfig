package local

import (
	"melato.org/cloudinit"
)

type OS interface {
	// InstallPackageCommand returns a command that installs a package.
	// The command should be a single string that is passed as input to sh.
	InstallPackageCommand(pkg string) string

	// AddUserCommand returns a command that creates the user.
	// The command should not configure sudo or ssh keys
	// The command is executed with the equivalent of execve(3)
	AddUserCommand(u *cloudinit.User) []string

	// NeedUserPasswords returns true if a new user should be
	// assigned a password.
	// In debian a locked user can login with ssh,
	// so the user does not need a password.
	// In alpine, a locked user cannot login,
	// so we need to keep the user unlocked and assign them an impossible
	// password in order to enable passwordless ssh login
	NeedUserPasswords() bool
}
