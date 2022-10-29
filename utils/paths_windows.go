package utils

func GetRootDir() string {
	return filepath.Join(os.Getenv("APPDATA"), ".eduhelper")
}
