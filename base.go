package cloudconfig

import (
	"io/fs"
)

type BaseConfigurer interface {
	// RunScript runs sh with the given input
	RunScript(input string) error

	// RunCommand runs program args[0], with args args
	RunCommand(args ...string) error

	// WriteFile writes a file, like os.WriteFile.  It should not try to create any directories.
	WriteFile(path string, data []byte, perm fs.FileMode) error

	// AppendFile appends to a file.  It should not try to create any directories.
	AppendFile(path string, data []byte, perm fs.FileMode) error

	// FileExists checks if the file exists.\
	// It should return an error only if an unexpected error occurs,
	// not if the file simply does not exist.
	FileExists(path string) (bool, error)
}

// Optional interface to return a user's home directory.
// It is used only if File.HomeDir is not specified
// Configurer.UserHomeDir() provides a default implementation.
type BaseUserHomeDir interface {
	// UserHomeDir returns a user's home directory.
	UserHomeDir(username string) (string, error)
}

// Optional interface to configure sudo privileges to a user
// If not implemented, Configurer.ApplySudo() is used.
type BaseApplySudo interface {
	// ApplySudo configures sudo privileges for a user
	// values is the list of sudo privileges
	// If values is empty, all sudo privileges should be applied.
	// There is no mechanism to specify no privileges
	ApplySudo(username string, values []string) error
}
