package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/clcollins/peon-ping/internal/config"
)

func Toggle(peonDir string) (string, error) {
	pausedPath := filepath.Join(peonDir, ".paused")

	if _, err := os.Stat(pausedPath); err == nil {
		if err := os.Remove(pausedPath); err != nil {
			return "", fmt.Errorf("cli: remove paused: %w", err)
		}
		return "Sounds resumed", nil
	}

	if err := os.WriteFile(pausedPath, []byte(""), 0644); err != nil {
		return "", fmt.Errorf("cli: create paused: %w", err)
	}
	return "Sounds paused", nil
}

func Status(cfgPath string, peonDir string) string {
	cfg, _ := config.Load(cfgPath)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Enabled: %v\n", cfg.Enabled))

	pausedPath := filepath.Join(peonDir, ".paused")
	if _, err := os.Stat(pausedPath); err == nil {
		b.WriteString("Status: PAUSED\n")
	} else {
		b.WriteString("Status: active\n")
	}

	b.WriteString(fmt.Sprintf("Pack: %s\n", cfg.DefaultPack))
	b.WriteString(fmt.Sprintf("Volume: %.1f\n", cfg.Volume))

	packsDir := filepath.Join(peonDir, "packs")
	packs, _ := List(packsDir)
	b.WriteString(fmt.Sprintf("Installed packs: %s\n", strings.Join(packs, ", ")))

	b.WriteString("Categories:\n")
	cats := make([]string, 0, len(cfg.Categories))
	for cat := range cfg.Categories {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	for _, cat := range cats {
		status := "on"
		if !cfg.Categories[cat] {
			status = "off"
		}
		b.WriteString(fmt.Sprintf("  %s: %s\n", cat, status))
	}

	return b.String()
}

func Use(packName string, cfgPath string, packsDir string) error {
	manifest := filepath.Join(packsDir, packName, "manifest.json")
	if _, err := os.Stat(manifest); err != nil {
		packs, _ := List(packsDir)
		return fmt.Errorf("cli: pack %q not found (available: %s)", packName, strings.Join(packs, ", "))
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("cli: load config: %w", err)
	}

	cfg.DefaultPack = packName
	if err := cfg.Save(cfgPath); err != nil {
		return fmt.Errorf("cli: save config: %w", err)
	}

	return nil
}

func List(packsDir string) ([]string, error) {
	entries, err := os.ReadDir(packsDir)
	if err != nil {
		return nil, nil
	}

	var packs []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		manifest := filepath.Join(packsDir, e.Name(), "manifest.json")
		if _, err := os.Stat(manifest); err == nil {
			packs = append(packs, e.Name())
		}
	}

	sort.Strings(packs)
	return packs, nil
}

func Volume(vol float64, cfgPath string) error {
	if vol < 0.0 || vol > 1.0 {
		return fmt.Errorf("cli: volume must be between 0.0 and 1.0, got %.1f", vol)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("cli: load config: %w", err)
	}

	cfg.Volume = vol
	if err := cfg.Save(cfgPath); err != nil {
		return fmt.Errorf("cli: save config: %w", err)
	}

	return nil
}

func Help() string {
	return `peon-ping: Sound notifications for Claude Code

Usage:
  peon toggle          Toggle sounds on/off
  peon status          Show current configuration
  peon use <pack>      Switch sound pack
  peon list            List available packs
  peon volume <0-1>    Set volume (0.0 to 1.0)
  peon help            Show this help

When invoked with no arguments, reads Claude Code hook JSON from stdin.`
}
