package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Path functions as variables for testing
var (
	getClineConfigPath = getClineConfigPathImpl
)

func init() {
	RegisterClient(&Client{
		Name:          "cline",
		DisplayName:   "Cline",
		GlobalPath:    func() (string, error) { return getClineConfigPath() },
		LocalPath:     nil,
		SupportsLocal: false,
		SyncFunc:      syncToMCPConfig,
	})
}

func getClineConfigPathImpl() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json"), nil
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json"), nil
	case "linux":
		return filepath.Join(home, ".config", "Code", "User", "globalStorage", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json"), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
