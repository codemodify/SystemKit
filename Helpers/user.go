package Helpers

import "os/user"

// IsRoot -
func IsRoot() bool {
	u, err := user.Current()

	if err != nil {
		return false
	}

	// On unix systems, root user either has the UID 0,
	// the GID 0 or both.
	return u.Uid == "0" || u.Gid == "0"
}

// HomeDir -
func HomeDir(returnIfError string) string {
	u, err := user.Current()
	if err != nil {
		return returnIfError
	}

	return u.HomeDir
}

// UserName -
func UserName(returnIfError string) string {
	u, err := user.Current()
	if err != nil {
		return returnIfError
	}

	return u.Username
}
