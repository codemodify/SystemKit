package Helpers

import (
	"os"

)

// FileOrFolderExists -
func FileOrFolderExists(name string) bool {
	_, error := os.Stat(name)
	return !os.IsNotExist(error)
}

