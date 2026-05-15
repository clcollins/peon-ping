package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Enabled                     bool            `json:"enabled"`
	DefaultPack                 string          `json:"default_pack"`
	Volume                      float64         `json:"volume"`
	DesktopNotifications        bool            `json:"desktop_notifications"`
	Categories                  map[string]bool `json:"categories"`
	AnnoyedThreshold            int             `json:"annoyed_threshold"`
	AnnoyedWindowSeconds        int             `json:"annoyed_window_seconds"`
	SessionStartCooldownSeconds int             `json:"session_start_cooldown_seconds"`
}

func Default() *Config {
	return &Config{
		Enabled:              true,
		DefaultPack:          "peon",
		Volume:               0.5,
		DesktopNotifications: true,
		Categories: map[string]bool{
			"session.start":    true,
			"task.acknowledge": false,
			"task.complete":    true,
			"task.error":       true,
			"input.required":   true,
			"resource.limit":   true,
			"user.spam":        true,
		},
		AnnoyedThreshold:            3,
		AnnoyedWindowSeconds:        10,
		SessionStartCooldownSeconds: 30,
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is from trusted config
	if err != nil {
		return Default(), nil
	}

	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return Default(), nil
	}

	return cfg, nil
}

func (c *Config) IsCategoryEnabled(category string) bool {
	if c.Categories != nil {
		if enabled, exists := c.Categories[category]; exists {
			return enabled
		}
	}

	if category == "task.acknowledge" {
		return false
	}
	return true
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("config: marshal: %w", err)
	}
	data = append(data, '\n')

	tmp := filepath.Join(filepath.Dir(path), ".config.json.tmp")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("config: write tmp: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("config: rename: %w", err)
	}

	return nil
}
