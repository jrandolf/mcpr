package clients

import (
	"os"
	"path/filepath"
)

// Path functions as variables for testing
var (
	getCursorConfigPath = getCursorConfigPathImpl
	getCursorLocalPath  = getCursorLocalPathImpl
)

func init() {
	RegisterClient(&Client{
		Name:          "cursor",
		DisplayName:   "Cursor",
		GlobalPath:    func() (string, error) { return getCursorConfigPath() },
		LocalPath:     func() (string, error) { return getCursorLocalPath() },
		SupportsLocal: true,
		SyncFunc:      syncToMCPConfig,
	})
}

func getCursorConfigPathImpl() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cursor", "mcp.json"), nil
}

func getCursorLocalPathImpl() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".cursor", "mcp.json"), nil
}
