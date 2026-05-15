package sound

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/clcollins/peon-ping/internal/config"
	"github.com/clcollins/peon-ping/internal/state"
)

func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("sound: read manifest: %w", err)
	}

	m := &Manifest{}
	if err := json.Unmarshal(data, m); err != nil {
		return nil, fmt.Errorf("sound: parse manifest: %w", err)
	}

	return m, nil
}

func PickSound(m *Manifest, category string, lastPlayed string, rng *rand.Rand) (string, error) {
	cat, ok := m.Categories[category]
	if !ok {
		return "", fmt.Errorf("sound: category %q not found in manifest", category)
	}

	if len(cat.Sounds) == 0 {
		return "", fmt.Errorf("sound: category %q has no sounds", category)
	}

	if len(cat.Sounds) == 1 {
		return cat.Sounds[0].File, nil
	}

	var candidates []ManifestSound
	for _, s := range cat.Sounds {
		if s.File != lastPlayed {
			candidates = append(candidates, s)
		}
	}

	if len(candidates) == 0 {
		candidates = cat.Sounds
	}

	choice := candidates[rng.Intn(len(candidates))]
	return choice.File, nil
}

func ResolvePack(cfg *config.Config, s *state.State, sessionID string, packsDir string) string {
	if override, ok := s.SessionPacks[sessionID]; ok {
		manifest := filepath.Join(packsDir, override, "manifest.json")
		if _, err := os.Stat(manifest); err == nil {
			return override
		}
	}

	return cfg.DefaultPack
}
