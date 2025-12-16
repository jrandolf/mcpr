package clients

import (
	"os"
	"path/filepath"
)

// Path functions as variables for testing
var (
	getGeminiConfigPath = getGeminiConfigPathImpl
	getGeminiLocalPath  = getGeminiLocalPathImpl
)

func init() {
	RegisterClient(&Client{
		Name:          "gemini",
		DisplayName:   "Gemini CLI",
		GlobalPath:    func() (string, error) { return getGeminiConfigPath() },
		LocalPath:     func() (string, error) { return getGeminiLocalPath() },
		SupportsLocal: true,
		SyncFunc:      syncToSettingsWithMcpServers,
	})
}

func getGeminiConfigPathImpl() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gemini", "settings.json"), nil
}

func getGeminiLocalPathImpl() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".gemini", "settings.json"), nil
}
