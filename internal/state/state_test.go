package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidState(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".state.json")

	data := []byte(`{
		"session_packs": {"sess-1": "glados"},
		"last_played": {"peon": "sounds/Done1.wav"},
		"prompt_times": {"sess-1": [1000, 2000, 3000]},
		"session_starts": {"sess-1": 1000}
	}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if s.SessionPacks["sess-1"] != "glados" {
		t.Errorf("session pack = %q, want %q", s.SessionPacks["sess-1"], "glados")
	}
	if s.LastPlayed["peon"] != "sounds/Done1.wav" {
		t.Errorf("last played = %q, want %q", s.LastPlayed["peon"], "sounds/Done1.wav")
	}
	if len(s.PromptTimes["sess-1"]) != 3 {
		t.Errorf("prompt times count = %d, want 3", len(s.PromptTimes["sess-1"]))
	}
	if s.SessionStarts["sess-1"] != 1000 {
		t.Errorf("session start = %d, want 1000", s.SessionStarts["sess-1"])
	}
}

func TestLoadMissingFile(t *testing.T) {
	s, err := Load("/nonexistent/.state.json")
	if err != nil {
		t.Fatalf("Load() should not error on missing file, got: %v", err)
	}

	if s.SessionPacks == nil {
		t.Error("missing file should return initialized state with non-nil maps")
	}
}

func TestLoadCorruptJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".state.json")

	if err := os.WriteFile(path, []byte("{{invalid"), 0o600); err != nil {
		t.Fatal(err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load() should not error on corrupt file, got: %v", err)
	}

	if s.SessionPacks == nil {
		t.Error("corrupt file should return initialized empty state")
	}
}

func TestLoadEmptyObject(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".state.json")

	if err := os.WriteFile(path, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}

	s, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if s.SessionPacks == nil {
		t.Error("empty object should still have initialized maps")
	}
	if s.LastPlayed == nil {
		t.Error("empty object should still have initialized LastPlayed map")
	}
	if s.PromptTimes == nil {
		t.Error("empty object should still have initialized PromptTimes map")
	}
	if s.LastEvent == nil {
		t.Error("empty object should still have initialized LastEvent map")
	}
	if s.SessionStarts == nil {
		t.Error("empty object should still have initialized SessionStarts map")
	}
}

func TestSaveAndRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".state.json")

	original := New()
	original.SessionPacks["sess-1"] = "glados"
	original.LastPlayed["peon"] = "sounds/Hello1.wav"
	original.PromptTimes["sess-1"] = []int64{100, 200}
	original.SessionStarts["sess-1"] = 100

	if err := original.Save(path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() after Save() error: %v", err)
	}

	if loaded.SessionPacks["sess-1"] != "glados" {
		t.Errorf("round-trip session pack = %q, want %q", loaded.SessionPacks["sess-1"], "glados")
	}
	if loaded.LastPlayed["peon"] != "sounds/Hello1.wav" {
		t.Errorf("round-trip last played = %q, want %q", loaded.LastPlayed["peon"], "sounds/Hello1.wav")
	}
	if len(loaded.PromptTimes["sess-1"]) != 2 {
		t.Errorf("round-trip prompt times count = %d, want 2", len(loaded.PromptTimes["sess-1"]))
	}
}

func TestSaveAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".state.json")

	s := New()
	if err := s.Save(path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

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
