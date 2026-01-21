package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	DefaultConfigDir  = ".config/debtq"
	DefaultConfigFile = "config.json"
)

// Config holds application configuration
type Config struct {
	ObsidianVaultPath string `json:"obsidian_vault_path"`
	DataFile          string `json:"data_file"`
	Currency          string `json:"currency"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		ObsidianVaultPath: filepath.Join(homeDir, "Documents", "obsidian-notes", "debtq"),
		DataFile:          filepath.Join(homeDir, DefaultConfigDir, "data.json"),
		Currency:          "INR",
	}
}

// GetConfigPath returns the config file path
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, DefaultConfigDir, DefaultConfigFile), nil
}

// Load loads configuration from file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		cfg := DefaultConfig()
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// EnsureObsidianDir ensures the obsidian directory exists
func (c *Config) EnsureObsidianDir() error {
	return os.MkdirAll(c.ObsidianVaultPath, 0755)
}
