package utils

import "os"

// FileExists returns if file at 'path' exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// DirExists returns if a directory at 'path' exists
func DirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
