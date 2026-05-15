package sound

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/clcollins/peon-ping/internal/config"
	"github.com/clcollins/peon-ping/internal/state"
)

func testManifest() *Manifest {
	return &Manifest{
		Name: "peon",
		Categories: map[string]ManifestCategory{
			"session.start": {
				Sounds: []ManifestSound{
					{File: "sounds/Hello1.wav", Label: "Ready to work"},
					{File: "sounds/Hello2.wav", Label: "Hmmm?"},
					{File: "sounds/Hello3.wav", Label: "What you want?"},
				},
			},
			"task.complete": {
				Sounds: []ManifestSound{
					{File: "sounds/Done1.wav", Label: "Work complete"},
				},
			},
		},
	}
}

func TestPickSoundFromCategory(t *testing.T) {
	m := testManifest()
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // deterministic seed for test

	file, err := PickSound(m, "session.start", "", rng)
	if err != nil {
		t.Fatalf("PickSound() error: %v", err)
	}
	if file == "" {
		t.Error("PickSound() returned empty file")
	}
}

func TestPickSoundAvoidsLastPlayed(t *testing.T) {
	m := testManifest()

	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		rng := rand.New(rand.NewSource(int64(i))) //nolint:gosec // deterministic seed for test
		file, err := PickSound(m, "session.start", "sounds/Hello1.wav", rng)
		if err != nil {
			t.Fatalf("PickSound() error: %v", err)
		}
		seen[file] = true
	}

	if seen["sounds/Hello1.wav"] {
		t.Error("PickSound() should avoid last-played file when alternatives exist")
	}
}

func TestPickSoundSingleReturnsOnly(t *testing.T) {
	m := testManifest()
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // deterministic seed for test

	file, err := PickSound(m, "task.complete", "sounds/Done1.wav", rng)
	if err != nil {
		t.Fatalf("PickSound() error: %v", err)
	}
	if file != "sounds/Done1.wav" {
		t.Errorf("single sound should return %q, got %q", "sounds/Done1.wav", file)
	}
}

func TestPickSoundMissingCategory(t *testing.T) {
	m := testManifest()
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // deterministic seed for test

	_, err := PickSound(m, "nonexistent.category", "", rng)
	if err == nil {
		t.Error("PickSound() should error on missing category")
	}
}

func TestPickSoundEmptyCategory(t *testing.T) {
	m := &Manifest{
		Name: "empty",
		Categories: map[string]ManifestCategory{
			"session.start": {Sounds: []ManifestSound{}},
		},
	}
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // deterministic seed for test

	_, err := PickSound(m, "session.start", "", rng)
	if err == nil {
		t.Error("PickSound() should error on empty category")
	}
}

func setupPackDir(t *testing.T, dir, name string) {
	t.Helper()
	packDir := filepath.Join(dir, name)
	if err := os.MkdirAll(packDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packDir, "manifest.json"), []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestResolvePackDefault(t *testing.T) {
	dir := t.TempDir()
	setupPackDir(t, dir, "peon")

	cfg := config.Default()
	s := state.New()

	pack := ResolvePack(cfg, s, "sess-1", dir)
	if pack != "peon" {
		t.Errorf("ResolvePack() = %q, want %q", pack, "peon")
	}
}

func TestResolvePackSessionOverride(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"peon", "glados"} {
		setupPackDir(t, dir, name)
	}

	cfg := config.Default()
	s := state.New()
	s.SessionPacks["sess-1"] = "glados"

	pack := ResolvePack(cfg, s, "sess-1", dir)
	if pack != "glados" {
		t.Errorf("ResolvePack() = %q, want %q", pack, "glados")
	}
}

func TestResolvePackSessionOverrideMissing(t *testing.T) {
	dir := t.TempDir()
	setupPackDir(t, dir, "peon")

	cfg := config.Default()
	s := state.New()
	s.SessionPacks["sess-1"] = "nonexistent"

	pack := ResolvePack(cfg, s, "sess-1", dir)
	if pack != "peon" {
		t.Errorf("ResolvePack() should fall back to default, got %q", pack)
	}
}
