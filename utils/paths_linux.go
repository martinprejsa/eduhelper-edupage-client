package utils

import (
	"os"
	"path/filepath"
)

func GetRootDir() string {
	path, _ := os.UserHomeDir()
	return filepath.Join(path, ".eduhelper")
}
