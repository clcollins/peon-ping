package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type State struct {
	SessionPacks  map[string]string           `json:"session_packs,omitempty"`
	LastPlayed    map[string]string           `json:"last_played,omitempty"`
	PromptTimes   map[string][]int64          `json:"prompt_times,omitempty"`
	LastEvent     map[string]map[string]int64 `json:"last_event,omitempty"`
	SessionStarts map[string]int64            `json:"session_starts,omitempty"`
}

func New() *State {
	return &State{
		SessionPacks:  make(map[string]string),
		LastPlayed:    make(map[string]string),
		PromptTimes:   make(map[string][]int64),
		LastEvent:     make(map[string]map[string]int64),
		SessionStarts: make(map[string]int64),
	}
}

func (s *State) init() {
	if s.SessionPacks == nil {
		s.SessionPacks = make(map[string]string)
	}
	if s.LastPlayed == nil {
		s.LastPlayed = make(map[string]string)
	}
	if s.PromptTimes == nil {
		s.PromptTimes = make(map[string][]int64)
	}
	if s.LastEvent == nil {
		s.LastEvent = make(map[string]map[string]int64)
	}
	if s.SessionStarts == nil {
		s.SessionStarts = make(map[string]int64)
	}
}

func Load(path string) (*State, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is from trusted config
	if err != nil {
		return New(), nil
	}

	s := &State{}
	if err := json.Unmarshal(data, s); err != nil {
		return New(), nil
	}

	s.init()
	return s, nil
}

func (s *State) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("state: marshal: %w", err)
	}
	data = append(data, '\n')

	tmp := filepath.Join(filepath.Dir(path), ".state.json.tmp")
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("state: write tmp: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("state: rename: %w", err)
	}

	return nil
}
