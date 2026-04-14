package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"transok/backend/consts"
)

/* Get environment variables */
func GetEnv() string {
	return consts.APP_INFO["env"]
}

func GetBasePath() string {
	env := GetEnv()
	fmt.Println("Environment:", env)
	if env == "" {
		env = "dev" // Default to development environment
	}

	if env == "dev" {
		return "data"
	}
	// Get the base storage path suitable for the current operating system
	var basePath string
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		basePath = filepath.Join(appData, "transok")
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		basePath = filepath.Join(homeDir, "Library", "Application Support", "transok")
	default: // Linux and other Unix-like systems
		basePath = "/var/lib/transok"
		// If not root, use user directory
		if os.Getuid() != 0 {
			homeDir, _ := os.UserHomeDir()
			basePath = filepath.Join(homeDir, ".transok")
		}
	}

	return basePath

}
