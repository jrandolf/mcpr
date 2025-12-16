package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Path functions as variables for testing
var (
	getWindsurfConfigPath = getWindsurfConfigPathImpl
	getWindsurfLocalPath  = getWindsurfLocalPathImpl
)

func init() {
	RegisterClient(&Client{
		Name:          "windsurf",
		DisplayName:   "Windsurf",
		GlobalPath:    func() (string, error) { return getWindsurfConfigPath() },
		LocalPath:     func() (string, error) { return getWindsurfLocalPath() },
		SupportsLocal: true,
		SyncFunc:      syncToMCPConfig,
	})
}

func getWindsurfConfigPathImpl() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Windsurf", "User", "globalStorage", "windsurf.mcp", "mcp.json"), nil
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Windsurf", "User", "globalStorage", "windsurf.mcp", "mcp.json"), nil
	case "linux":
		return filepath.Join(home, ".config", "Windsurf", "User", "globalStorage", "windsurf.mcp", "mcp.json"), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func getWindsurfLocalPathImpl() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".windsurf", "mcp.json"), nil
}
