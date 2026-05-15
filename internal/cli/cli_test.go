package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToggleCreatesPaused(t *testing.T) {
	dir := t.TempDir()

	msg, err := Toggle(dir)
	if err != nil {
		t.Fatalf("Toggle() error: %v", err)
	}

	pausedPath := filepath.Join(dir, ".paused")
	if _, err := os.Stat(pausedPath); err != nil {
		t.Error(".paused file should exist after toggle on")
	}
	if !strings.Contains(msg, "paused") {
		t.Errorf("message should mention paused, got %q", msg)
	}
}

func TestToggleRemovesPaused(t *testing.T) {
	dir := t.TempDir()
	pausedPath := filepath.Join(dir, ".paused")
	os.WriteFile(pausedPath, []byte(""), 0644)

	msg, err := Toggle(dir)
	if err != nil {
		t.Fatalf("Toggle() error: %v", err)
	}

	if _, err := os.Stat(pausedPath); !os.IsNotExist(err) {
		t.Error(".paused file should be removed after toggle off")
	}
	if !strings.Contains(msg, "resumed") {
		t.Errorf("message should mention resumed, got %q", msg)
	}
}

func TestStatusShowsInfo(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "config.json")
	os.WriteFile(cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`), 0644)

	packDir := filepath.Join(dir, "packs", "peon")
	os.MkdirAll(packDir, 0755)
	os.WriteFile(filepath.Join(packDir, "manifest.json"), []byte("{}"), 0644)

	output := Status(cfgPath, dir)

	checks := []string{"enabled", "peon", "0.5"}
	for _, want := range checks {
		if !strings.Contains(strings.ToLower(output), want) {
			t.Errorf("Status() output should contain %q, got:\n%s", want, output)
		}
	}
}

func TestStatusShowsPaused(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "config.json")
	os.WriteFile(cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`), 0644)
	os.WriteFile(filepath.Join(dir, ".paused"), []byte(""), 0644)

	output := Status(cfgPath, dir)
	if !strings.Contains(strings.ToLower(output), "paused") {
		t.Errorf("Status() should show paused state, got:\n%s", output)
	}
}

func TestUseValidPack(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "config.json")
	os.WriteFile(cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`), 0644)

	packsDir := filepath.Join(dir, "packs")
	for _, name := range []string{"peon", "glados"} {
		packDir := filepath.Join(packsDir, name)
		os.MkdirAll(packDir, 0755)
		os.WriteFile(filepath.Join(packDir, "manifest.json"), []byte("{}"), 0644)
	}

	err := Use("glados", cfgPath, packsDir)
	if err != nil {
		t.Fatalf("Use() error: %v", err)
	}
}

func TestUseInvalidPack(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "config.json")
	os.WriteFile(cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`), 0644)

	packsDir := filepath.Join(dir, "packs")
	os.MkdirAll(filepath.Join(packsDir, "peon"), 0755)
	os.WriteFile(filepath.Join(packsDir, "peon", "manifest.json"), []byte("{}"), 0644)

	err := Use("nonexistent", cfgPath, packsDir)
	if err == nil {
		t.Error("Use() should error on nonexistent pack")
	}
}

func TestListPacks(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"peon", "glados", "nopack"} {
		packDir := filepath.Join(dir, name)
		os.MkdirAll(packDir, 0755)
		if name != "nopack" {
			os.WriteFile(filepath.Join(packDir, "manifest.json"), []byte("{}"), 0644)
		}
	}

	packs, err := List(dir)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(packs) != 2 {
		t.Errorf("List() returned %d packs, want 2 (should skip nopack)", len(packs))
	}
}

func TestVolumeValid(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	os.WriteFile(cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`), 0644)

	err := Volume(0.8, cfgPath)
	if err != nil {
		t.Fatalf("Volume() error: %v", err)
	}
}

func TestVolumeOutOfRange(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	os.WriteFile(cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`), 0644)

	if err := Volume(1.5, cfgPath); err == nil {
		t.Error("Volume(1.5) should error")
	}
	if err := Volume(-0.1, cfgPath); err == nil {
		t.Error("Volume(-0.1) should error")
	}
}

func TestHelp(t *testing.T) {
	output := Help()
	if output == "" {
		t.Error("Help() should return non-empty string")
	}
	if !strings.Contains(output, "toggle") {
		t.Error("Help() should mention toggle command")
	}
}
