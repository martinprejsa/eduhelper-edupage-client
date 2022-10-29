package utils

import "path/filepath"

func GetCredentialsFilePath() string {
	return filepath.Join(GetRootDir(), "credentials")
}
