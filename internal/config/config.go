package config

import (
	"fmt"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	DelayMs      int    `json:"delay_ms"`
	LogLevel     string `json:"log_level"`
	MaxLogLines  int    `json:"max_log_lines"`
	WindowWidth  int    `json:"window_width"`
	WindowHeight int    `json:"window_height"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		DelayMs:      50,
		LogLevel:     "info",
		MaxLogLines:  100,
		WindowWidth:  500,
		WindowHeight: 450,
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
