package os

import (
	"os"
	"path/filepath"
)

// Readonly variables for executable path and directory
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

func GetExeDir() string {
	return execDir
}
