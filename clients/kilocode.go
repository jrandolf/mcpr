package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Path functions as variables for testing
var (
	getKiloCodeConfigPath = getKiloCodeConfigPathImpl
	getKiloCodeLocalPath  = getKiloCodeLocalPathImpl
)

func init() {
	RegisterClient(&Client{
		Name:          "kilo-code",
		DisplayName:   "Kilo Code",
		GlobalPath:    func() (string, error) { return getKiloCodeConfigPath() },
		LocalPath:     func() (string, error) { return getKiloCodeLocalPath() },
		SupportsLocal: true,
		SyncFunc:      syncToMCPConfig,
	})
}

func getKiloCodeConfigPathImpl() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "globalStorage", "kilocode.kilo-code", "settings", "mcp_settings.json"), nil
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "Code", "User", "globalStorage", "kilocode.kilo-code", "settings", "mcp_settings.json"), nil
	case "linux":
		return filepath.Join(home, ".config", "Code", "User", "globalStorage", "kilocode.kilo-code", "settings", "mcp_settings.json"), nil
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func getKiloCodeLocalPathImpl() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".kilocode", "mcp.json"), nil
}
