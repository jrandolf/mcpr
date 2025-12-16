package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Path functions as variables for testing
var (
	getZencoderConfigPath = getZencoderConfigPathImpl
)

func init() {
	RegisterClient(&Client{
		Name:          "zencoder",
		DisplayName:   "ZenCoder",
		GlobalPath:    func() (string, error) { return getZencoderConfigPath() },
		LocalPath:     nil,
		SupportsLocal: false,
		SyncFunc:      syncToMCPConfig,
	})
}

func getZencoderConfigPathImpl() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "globalStorage", "zencoderAI.zencoder", "mcp_settings.json"), nil
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Code", "User", "globalStorage", "zencoderAI.zencoder", "mcp_settings.json"), nil
	case "linux":
		return filepath.Join(home, ".config", "Code", "User", "globalStorage", "zencoderAI.zencoder", "mcp_settings.json"), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
