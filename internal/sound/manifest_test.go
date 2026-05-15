package sound

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadManifestValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.json")

	data := []byte(`{
		"name": "peon",
		"display_name": "Orc Peon",
		"categories": {
			"session.start": {
				"sounds": [
					{"file": "sounds/Hello1.wav", "label": "Ready to work"},
					{"file": "sounds/Hello2.wav", "label": "Hmmm?"}
				]
			},
			"task.complete": {
				"sounds": [
					{"file": "sounds/Done1.wav", "label": "Work complete"}
				]
			}
		}
	}`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	m, err := LoadManifest(path)
	if err != nil {
		t.Fatalf("LoadManifest() error: %v", err)
	}

	if m.Name != "peon" {
		t.Errorf("name = %q, want %q", m.Name, "peon")
	}
	if len(m.Categories) != 2 {
		t.Errorf("categories count = %d, want 2", len(m.Categories))
	}
	if len(m.Categories["session.start"].Sounds) != 2 {
		t.Errorf("session.start sounds count = %d, want 2", len(m.Categories["session.start"].Sounds))
	}
}

func TestLoadManifestMissing(t *testing.T) {
	_, err := LoadManifest("/nonexistent/manifest.json")
	if err == nil {
		t.Error("LoadManifest() should error on missing file")
	}
}

func TestLoadManifestCorrupt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "manifest.json")

	if err := os.WriteFile(path, []byte("{bad json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadManifest(path)
	if err == nil {
		t.Error("LoadManifest() should error on corrupt JSON")
	}
}
