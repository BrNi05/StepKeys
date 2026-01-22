package os

import (
	"os"
	"path/filepath"
)

var (
	execPath string
	execDir  string
)

func init() {
	// The path to the executable (can be a symlink)
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Resolve symlinks
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		panic(err)
	}

	execDir = filepath.Dir(exePath)
	execPath = exePath
}

// Returns the execDir
func GetExecDir() string {
	return execDir
}

// Returns the execPath
func GetExecPath() string {
	return execPath
}
