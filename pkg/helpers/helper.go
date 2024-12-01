package helpers

import (
	"os"
	"path"
)

func GetSocketPath() string {
	return path.Join(os.Getenv("XDG_RUNTIME_DIR"), "app/org.keepassxc.KeePassXC/org.keepassxc.KeePassXC.BrowserServer")
}

func GetStoragePath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(userConfigDir, "gokeexc.json"), nil
}
