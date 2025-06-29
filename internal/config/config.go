package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	DelayMs        int    `json:"delay_ms"`
	LogLevel       string `json:"log_level"`
	MaxLogLines    int    `json:"max_log_lines"`
	WindowWidth    int    `json:"window_width"`
	WindowHeight   int    `json:"window_height"`
	MinimizeToTray bool   `json:"minimize_to_tray"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		DelayMs:        50,
		LogLevel:       "info",
		MaxLogLines:    100,
		WindowWidth:    500,
		WindowHeight:   450,
		MinimizeToTray: true,
	}
}

// ValidateDelay validates the delay value
func (c *Config) ValidateDelay() error {
	if c.DelayMs < 1 || c.DelayMs > 500 {
		return fmt.Errorf("delay must be between 1 and 500 milliseconds")
	}
	return nil
}

// ParseDelay parses a delay string and validates it
func ParseDelay(delayStr string) (int, error) {
	if delayStr == "" {
		return DefaultConfig().DelayMs, nil
	}

	delayMs, err := strconv.Atoi(delayStr)
	if err != nil {
		return 0, fmt.Errorf("invalid delay format: %v", err)
	}

	if delayMs < 1 || delayMs > 500 {
		return 0, fmt.Errorf("delay must be between 1 and 500 milliseconds")
	}

	return delayMs, nil
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %v", err)
	}

	appConfigDir := filepath.Join(configDir, "ClickGuardian")
	if err := os.MkdirAll(appConfigDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create app config directory: %v", err)
	}

	return filepath.Join(appConfigDir, "config.json"), nil
}

// LoadConfig loads configuration from file, or returns default if file doesn't exist
func LoadConfig() *Config {
	configPath, err := GetConfigPath()
	if err != nil {
		return DefaultConfig()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// If file doesn't exist, return default config
		if os.IsNotExist(err) {
			return DefaultConfig()
		}
		return DefaultConfig()
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultConfig()
	}

	// Validate loaded config and use defaults for invalid values
	if config.DelayMs < 5 || config.DelayMs > 500 {
		config.DelayMs = DefaultConfig().DelayMs
	}
	if config.MaxLogLines <= 0 {
		config.MaxLogLines = DefaultConfig().MaxLogLines
	}
	if config.WindowWidth <= 0 {
		config.WindowWidth = DefaultConfig().WindowWidth
	}
	if config.WindowHeight <= 0 {
		config.WindowHeight = DefaultConfig().WindowHeight
	}
	if config.LogLevel == "" {
		config.LogLevel = DefaultConfig().LogLevel
	}

	return &config
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %v", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
