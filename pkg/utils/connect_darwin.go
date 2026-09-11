package utils

import (
	"os"
	"path"
)

func GetSocketPath() string {
	return path.Join(os.Getenv("TMPDIR"), "org.keepassxc.KeePassXC.BrowserServer")
}

func GetConfigPath() string {
	return path.Join(os.Getenv("HOME"), ".config", "gopassxc.json")
}
