package event

import (
	"testing"
	"time"

	"github.com/clcollins/peon-ping/internal/state"
)

func TestCheckSpam(t *testing.T) {
	baseTime := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		times     []int64
		threshold int
		window    int
		want      bool
	}{
		{
			name:      "first prompt is not spam",
			times:     nil,
			threshold: 3,
			window:    10,
			want:      false,
		},
		{
			name:      "under threshold is not spam",
			times:     []int64{baseTime.Unix() - 2, baseTime.Unix() - 1},
			threshold: 3,
			window:    10,
			want:      false,
		},
		{
			name: "at threshold within window IS spam",
			times: []int64{
				baseTime.Unix() - 3,
				baseTime.Unix() - 2,
				baseTime.Unix() - 1,
			},
			threshold: 3,
			window:    10,
			want:      true,
		},
		{
			name: "over threshold but outside window is not spam",
			times: []int64{
				baseTime.Unix() - 30,
				baseTime.Unix() - 20,
				baseTime.Unix() - 15,
			},
			threshold: 3,
			window:    10,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state.New()
			s.PromptTimes["sess-1"] = tt.times

			clock := FixedClock{T: baseTime}
			got := CheckSpam(s, "sess-1", tt.threshold, tt.window, clock)
			if got != tt.want {
				t.Errorf("CheckSpam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckSpamUpdatesState(t *testing.T) {
	baseTime := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)
	s := state.New()
	clock := FixedClock{T: baseTime}

	CheckSpam(s, "sess-1", 3, 10, clock)

	times := s.PromptTimes["sess-1"]
	if len(times) != 1 {
		t.Fatalf("expected 1 prompt time after CheckSpam, got %d", len(times))
	}
	if times[0] != baseTime.Unix() {
		t.Errorf("prompt time = %d, want %d", times[0], baseTime.Unix())
	}
}

func TestCheckCooldown(t *testing.T) {
	baseTime := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		lastTime int64
		cooldown int
		want     bool
	}{
		{
			name:     "first event is allowed",
			lastTime: 0,
			cooldown: 30,
			want:     false,
		},
		{
			name:     "within cooldown is suppressed",
			lastTime: baseTime.Unix() - 10,
			cooldown: 30,
			want:     true,
		},
		{
			name:     "after cooldown is allowed",
			lastTime: baseTime.Unix() - 60,
			cooldown: 30,
			want:     false,
		},
		{
			name:     "exactly at cooldown boundary is allowed",
			lastTime: baseTime.Unix() - 30,
			cooldown: 30,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state.New()
			if tt.lastTime > 0 {
				s.LastEvent["sess-1"] = map[string]int64{
					"SessionStart": tt.lastTime,
				}
			}

			clock := FixedClock{T: baseTime}
			got := CheckCooldown(s, "sess-1", "SessionStart", tt.cooldown, clock)
			if got != tt.want {
				t.Errorf("CheckCooldown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckCooldownUpdatesState(t *testing.T) {
	baseTime := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)
	s := state.New()
	clock := FixedClock{T: baseTime}

	CheckCooldown(s, "sess-1", "SessionStart", 30, clock)

	ts, ok := s.LastEvent["sess-1"]["SessionStart"]
	if !ok {
		t.Fatal("expected SessionStart timestamp in state after CheckCooldown")
	}
	if ts != baseTime.Unix() {
		t.Errorf("event timestamp = %d, want %d", ts, baseTime.Unix())
	}
}
