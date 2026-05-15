package event

import (
	"github.com/clcollins/peon-ping/internal/state"
)

func CheckSpam(s *state.State, sessionID string, threshold, windowSecs int, clock Clock) bool {
	now := clock.Now().Unix()

	cutoff := now - int64(windowSecs)
	existing := s.PromptTimes[sessionID]

	var recent []int64
	for _, ts := range existing {
		if ts >= cutoff {
			recent = append(recent, ts)
		}
	}

	isSpam := len(recent) >= threshold

	s.PromptTimes[sessionID] = append(recent, now)

	return isSpam
}

func CheckCooldown(s *state.State, sessionID, eventName string, cooldownSecs int, clock Clock) bool {
	now := clock.Now().Unix()

	if events, ok := s.LastEvent[sessionID]; ok {
		if lastTime, ok := events[eventName]; ok {
			if now-lastTime < int64(cooldownSecs) {
				return true
			}
		}
	}

	if s.LastEvent[sessionID] == nil {
		s.LastEvent[sessionID] = make(map[string]int64)
	}
	s.LastEvent[sessionID][eventName] = now

	return false
}
