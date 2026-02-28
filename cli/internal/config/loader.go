package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigDir  = ".wordmon"
	DefaultConfigFile = "wordmon.yaml"
)

// Manager handles loading and saving configuration
type Manager struct {
	configPath string
}

// NewManager creates a new config manager with the default config path
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, DefaultConfigDir, DefaultConfigFile)
	return &Manager{configPath: configPath}, nil
}

// NewManagerWithPath creates a new config manager with a custom config path
func NewManagerWithPath(configPath string) *Manager {
	return &Manager{configPath: configPath}
}

// GetConfigPath returns the path to the config file
func (m *Manager) GetConfigPath() string {
	return m.configPath
}

// GetConfigDir returns the directory containing the config file
func (m *Manager) GetConfigDir() string {
	return filepath.Dir(m.configPath)
}

// ConfigExists checks if the config file exists
func (m *Manager) ConfigExists() bool {
	_, err := os.Stat(m.configPath)
	return err == nil
}

// Load reads and parses the configuration file
func (m *Manager) Load() (*Config, error) {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save writes the configuration to disk using atomic writes
func (m *Manager) Save(config *Config) error {
	// Ensure config directory exists
	configDir := m.GetConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to temporary file
	tmpPath := m.configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write temp config file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, m.configPath); err != nil {
		os.Remove(tmpPath) // Cleanup on failure
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

// Initialize creates a new config file with default values
func (m *Manager) Initialize() error {
	if m.ConfigExists() {
		return fmt.Errorf("config file already exists at %s", m.configPath)
	}

	config := DefaultConfig()
	return m.Save(config)
}
