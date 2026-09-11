package utils

import (
	"os"
	"path"
)

func GetSocketPath() string {
	if runtimeDir := os.Getenv("XDG_RUNTIME_DIR"); runtimeDir != "" {
		return path.Join(runtimeDir, "org.keepassxc.KeePassXC.BrowserServer")
	}
	return "/run/user/1000/org.keepassxc.KeePassXC.BrowserServer"
}

func GetConfigPath() string {
	return path.Join(os.Getenv("HOME"), ".config", "gopassxc.json")
}
