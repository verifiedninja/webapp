package filesystem

import (
	"os"
)

func FolderExists(f string) bool {
	return FileExists(f)
}

func FileExists(f string) bool {
	if _, err := os.Stat(f); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
