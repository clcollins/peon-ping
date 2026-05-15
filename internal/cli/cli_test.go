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

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func mkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o750); err != nil {
		t.Fatal(err)
	}
}

func TestToggleRemovesPaused(t *testing.T) {
	dir := t.TempDir()
	pausedPath := filepath.Join(dir, ".paused")
	writeFile(t, pausedPath, []byte(""))

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
	writeFile(t, cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`))

	packDir := filepath.Join(dir, "packs", "peon")
	mkdirAll(t, packDir)
	writeFile(t, filepath.Join(packDir, "manifest.json"), []byte("{}"))

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
	writeFile(t, cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`))
	writeFile(t, filepath.Join(dir, ".paused"), []byte(""))

	output := Status(cfgPath, dir)
	if !strings.Contains(strings.ToLower(output), "paused") {
		t.Errorf("Status() should show paused state, got:\n%s", output)
	}
}

func TestUseValidPack(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "config.json")
	writeFile(t, cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`))

	packsDir := filepath.Join(dir, "packs")
	for _, name := range []string{"peon", "glados"} {
		packDir := filepath.Join(packsDir, name)
		mkdirAll(t, packDir)
		writeFile(t, filepath.Join(packDir, "manifest.json"), []byte("{}"))
	}

	err := Use("glados", cfgPath, packsDir)
	if err != nil {
		t.Fatalf("Use() error: %v", err)
	}
}

func TestUseInvalidPack(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "config.json")
	writeFile(t, cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`))

	packsDir := filepath.Join(dir, "packs")
	mkdirAll(t, filepath.Join(packsDir, "peon"))
	writeFile(t, filepath.Join(packsDir, "peon", "manifest.json"), []byte("{}"))

	err := Use("nonexistent", cfgPath, packsDir)
	if err == nil {
		t.Error("Use() should error on nonexistent pack")
	}
}

func TestListPacks(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"peon", "glados", "nopack"} {
		packDir := filepath.Join(dir, name)
		mkdirAll(t, packDir)
		if name != "nopack" {
			writeFile(t, filepath.Join(packDir, "manifest.json"), []byte("{}"))
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
	writeFile(t, cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`))

	err := Volume(0.8, cfgPath)
	if err != nil {
		t.Fatalf("Volume() error: %v", err)
	}
}

func TestVolumeOutOfRange(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	writeFile(t, cfgPath, []byte(`{"enabled":true,"default_pack":"peon","volume":0.5}`))

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
