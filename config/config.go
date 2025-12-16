package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = "mcpr.json"

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Name    string            `json:"name"`
	Type    string            `json:"type"` // "stdio" or "http"
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// SyncedClient represents a client that has been synced
type SyncedClient struct {
	Name    string   `json:"name"`              // Client name (e.g., "claude-desktop")
	Local   bool     `json:"local"`             // Whether synced to local config
	Servers []string `json:"servers,omitempty"` // Specific servers synced (empty = all)
}

// Config holds all configured MCP servers
type Config struct {
	Servers       []MCPServer    `json:"servers"`
	SyncedClients []SyncedClient `json:"synced_clients,omitempty"`
	path          string         // path where config was loaded from or will be saved to
}

// findConfigInParents searches for config file in current and parent directories
func findConfigInParents() (string, bool) {
	dir, err := os.Getwd()
	if err != nil {
		return "", false
	}

	for {
		configPath := filepath.Join(dir, configFileName)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", false
}

// getGlobalConfigPath returns the global config path at ~/.config/mcpr/config.json
func getGlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "mcpr", "config.json"), nil
}

// GetConfigPath returns the path to the mcpr config file
// It searches in the following order:
// 1. Current directory and parent directories for mcpr.json
// 2. ~/.config/mcpr/config.json
func GetConfigPath() (string, error) {
	// First check parent directories
	if path, found := findConfigInParents(); found {
		return path, nil
	}

	// Fall back to global config
	return getGlobalConfigPath()
}

// GetWriteConfigPath returns the path where new config should be written
// Prefers local directory if mcpr.json exists, otherwise uses global config
func GetWriteConfigPath(preferLocal bool) (string, error) {
	if preferLocal {
		// Check if local config exists
		if path, found := findConfigInParents(); found {
			return path, nil
		}
		// Create in current directory
		return configFileName, nil
	}
	return getGlobalConfigPath()
}

// Load reads the config from disk
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// Return empty config, will be saved to global path
		globalPath, _ := getGlobalConfigPath()
		return &Config{Servers: []MCPServer{}, path: globalPath}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	cfg.path = path

	return &cfg, nil
}

// LoadFromPath reads the config from a specific path
func LoadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{Servers: []MCPServer{}, path: path}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	cfg.path = path

	return &cfg, nil
}

// Path returns the path where this config was loaded from or will be saved to
func (c *Config) Path() string {
	return c.path
}

// SetPath sets the path where this config will be saved
func (c *Config) SetPath(path string) {
	c.path = path
}

// Save writes the config to disk
func (c *Config) Save() error {
	if c.path == "" {
		path, err := getGlobalConfigPath()
		if err != nil {
			return err
		}
		c.path = path
	}

	// Ensure directory exists
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// AddServer adds a new MCP server to the config
func (c *Config) AddServer(server MCPServer) error {
	for _, s := range c.Servers {
		if s.Name == server.Name {
			return fmt.Errorf("server %q already exists", server.Name)
		}
	}
	c.Servers = append(c.Servers, server)
	return nil
}

// RemoveServer removes an MCP server from the config by name
func (c *Config) RemoveServer(name string) error {
	for i, s := range c.Servers {
		if s.Name == name {
			c.Servers = append(c.Servers[:i], c.Servers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("server %q not found", name)
}

// GetServer retrieves a server by name
func (c *Config) GetServer(name string) (*MCPServer, error) {
	for _, s := range c.Servers {
		if s.Name == name {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("server %q not found", name)
}

// ListServers returns all configured servers
func (c *Config) ListServers() []MCPServer {
	return c.Servers
}

// AddSyncedClient adds or updates a synced client record
func (c *Config) AddSyncedClient(clientName string, local bool, servers []string) {
	// Check if client already exists and update it
	for i, sc := range c.SyncedClients {
		if sc.Name == clientName && sc.Local == local {
			c.SyncedClients[i].Servers = servers
			return
		}
	}
	// Add new synced client
	c.SyncedClients = append(c.SyncedClients, SyncedClient{
		Name:    clientName,
		Local:   local,
		Servers: servers,
	})
}

// RemoveSyncedClient removes a synced client record
func (c *Config) RemoveSyncedClient(clientName string, local bool) {
	for i, sc := range c.SyncedClients {
		if sc.Name == clientName && sc.Local == local {
			c.SyncedClients = append(c.SyncedClients[:i], c.SyncedClients[i+1:]...)
			return
		}
	}
}

// GetSyncedClients returns all synced client records
func (c *Config) GetSyncedClients() []SyncedClient {
	return c.SyncedClients
}

// GetSyncedClient returns a specific synced client by name and local flag
func (c *Config) GetSyncedClient(clientName string, local bool) *SyncedClient {
	for _, sc := range c.SyncedClients {
		if sc.Name == clientName && sc.Local == local {
			return &sc
		}
	}
	return nil
}
