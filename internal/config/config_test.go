package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if !cfg.Enabled {
		t.Error("default config should be enabled")
	}
	if cfg.DefaultPack != "peon" {
		t.Errorf("default pack = %q, want %q", cfg.DefaultPack, "peon")
	}
	if cfg.Volume != 0.5 {
		t.Errorf("default volume = %f, want 0.5", cfg.Volume)
	}
	if !cfg.DesktopNotifications {
		t.Error("default desktop_notifications should be true")
	}
	if cfg.AnnoyedThreshold != 3 {
		t.Errorf("default annoyed_threshold = %d, want 3", cfg.AnnoyedThreshold)
	}
	if cfg.AnnoyedWindowSeconds != 10 {
		t.Errorf("default annoyed_window_seconds = %d, want 10", cfg.AnnoyedWindowSeconds)
	}
	if cfg.SessionStartCooldownSeconds != 30 {
		t.Errorf("default session_start_cooldown_seconds = %d, want 30", cfg.SessionStartCooldownSeconds)
	}
}

func TestDefaultCategories(t *testing.T) {
	cfg := Default()

	expectations := map[string]bool{
		"session.start":   true,
		"task.acknowledge": false,
		"task.complete":    true,
		"task.error":       true,
		"input.required":   true,
		"resource.limit":   true,
		"user.spam":        true,
	}

	for cat, want := range expectations {
		got, exists := cfg.Categories[cat]
		if !exists {
			t.Errorf("category %q missing from defaults", cat)
			continue
		}
		if got != want {
			t.Errorf("category %q = %v, want %v", cat, got, want)
		}
	}
}

func TestLoadValidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	data := []byte(`{
		"enabled": false,
		"default_pack": "glados",
		"volume": 0.8,
		"desktop_notifications": false,
		"categories": {"task.complete": false},
		"annoyed_threshold": 5,
		"annoyed_window_seconds": 20,
		"session_start_cooldown_seconds": 60
	}`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Enabled {
		t.Error("loaded config should be disabled")
	}
	if cfg.DefaultPack != "glados" {
		t.Errorf("loaded pack = %q, want %q", cfg.DefaultPack, "glados")
	}
	if cfg.Volume != 0.8 {
		t.Errorf("loaded volume = %f, want 0.8", cfg.Volume)
	}
	if cfg.DesktopNotifications {
		t.Error("loaded desktop_notifications should be false")
	}
	if cfg.AnnoyedThreshold != 5 {
		t.Errorf("loaded annoyed_threshold = %d, want 5", cfg.AnnoyedThreshold)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.json")
	if err != nil {
		t.Fatalf("Load() should not error on missing file, got: %v", err)
	}

	if cfg.DefaultPack != "peon" {
		t.Errorf("missing file should return defaults, got pack = %q", cfg.DefaultPack)
	}
	if !cfg.Enabled {
		t.Error("missing file should return defaults (enabled=true)")
	}
}

func TestLoadCorruptJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	if err := os.WriteFile(path, []byte("{not valid json"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() should not error on corrupt file, got: %v", err)
	}

	if cfg.DefaultPack != "peon" {
		t.Errorf("corrupt file should return defaults, got pack = %q", cfg.DefaultPack)
	}
}

func TestIsCategoryEnabled(t *testing.T) {
	tests := []struct {
		name     string
		category string
		cats     map[string]bool
		want     bool
	}{
		{
			name:     "explicitly enabled category",
			category: "task.complete",
			cats:     map[string]bool{"task.complete": true},
			want:     true,
		},
		{
			name:     "explicitly disabled category",
			category: "task.complete",
			cats:     map[string]bool{"task.complete": false},
			want:     false,
		},
		{
			name:     "task.acknowledge defaults off when missing",
			category: "task.acknowledge",
			cats:     map[string]bool{},
			want:     false,
		},
		{
			name:     "other categories default on when missing",
			category: "task.complete",
			cats:     map[string]bool{},
			want:     true,
		},
		{
			name:     "unknown category defaults on",
			category: "something.new",
			cats:     map[string]bool{},
			want:     true,
		},
		{
			name:     "nil categories map defaults on",
			category: "task.complete",
			cats:     nil,
			want:     true,
		},
		{
			name:     "nil categories task.acknowledge defaults off",
			category: "task.acknowledge",
			cats:     nil,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Categories: tt.cats}
			got := cfg.IsCategoryEnabled(tt.category)
			if got != tt.want {
				t.Errorf("IsCategoryEnabled(%q) = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}

func TestSaveAndRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	original := Default()
	original.DefaultPack = "sc_kerrigan"
	original.Volume = 0.3

	if err := original.Save(path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() after Save() error: %v", err)
	}

	if loaded.DefaultPack != "sc_kerrigan" {
		t.Errorf("round-trip pack = %q, want %q", loaded.DefaultPack, "sc_kerrigan")
	}
	if loaded.Volume != 0.3 {
		t.Errorf("round-trip volume = %f, want 0.3", loaded.Volume)
	}
	if !loaded.Enabled {
		t.Error("round-trip should preserve enabled=true")
	}
}

func TestSaveAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	cfg := Default()
	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Verify no temp files left behind
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		t.Errorf("expected 1 file after save, got %d: %v", len(entries), names)
	}
}
