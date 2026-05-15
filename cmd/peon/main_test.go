package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/clcollins/peon-ping/internal/event"
	"github.com/clcollins/peon-ping/internal/notifier"
	"github.com/clcollins/peon-ping/internal/player"
)

func writeTestFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	writeTestFile(t, filepath.Join(dir, "config.json"), []byte(`{
		"enabled": true,
		"default_pack": "peon",
		"volume": 0.5,
		"desktop_notifications": true,
		"categories": {
			"session.start": true,
			"task.acknowledge": false,
			"task.complete": true,
			"task.error": true,
			"input.required": true,
			"resource.limit": true,
			"user.spam": true
		},
		"annoyed_threshold": 3,
		"annoyed_window_seconds": 10,
		"session_start_cooldown_seconds": 30
	}`))

	writeTestFile(t, filepath.Join(dir, ".state.json"), []byte("{}"))

	packDir := filepath.Join(dir, "packs", "peon", "sounds")
	if err := os.MkdirAll(packDir, 0o750); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(dir, "packs", "peon", "manifest.json"), []byte(`{
		"name": "peon",
		"categories": {
			"session.start": {"sounds": [{"file": "sounds/Hello1.wav"}]},
			"task.complete": {"sounds": [{"file": "sounds/Done1.wav"}]},
			"task.error": {"sounds": [{"file": "sounds/Error1.wav"}]},
			"input.required": {"sounds": [{"file": "sounds/Input1.wav"}]},
			"resource.limit": {"sounds": [{"file": "sounds/Limit1.wav"}]},
			"user.spam": {"sounds": [{"file": "sounds/Spam1.wav"}]}
		}
	}`))

	for _, f := range []string{"Hello1.wav", "Done1.wav", "Error1.wav", "Input1.wav", "Limit1.wav", "Spam1.wav"} {
		writeTestFile(t, filepath.Join(packDir, f), []byte("fake wav"))
	}

	return dir
}

func TestProcessHookSessionStart(t *testing.T) {
	dir := setupTestDir(t)
	mp := &player.MockPlayer{}
	mn := &notifier.MockNotifier{}
	clock := event.FixedClock{T: fixedTime()}

	ev := &event.HookEvent{
		EventName: "SessionStart",
		SessionID: "test-session",
	}

	err := processHook(ev, dir, mp, mn, &clock)
	if err != nil {
		t.Fatalf("processHook() error: %v", err)
	}

	if len(mp.PlayCalls) != 1 {
		t.Fatalf("expected 1 play call, got %d", len(mp.PlayCalls))
	}
	if !strings.Contains(mp.PlayCalls[0].File, "Hello1.wav") {
		t.Errorf("expected session.start sound, got %q", mp.PlayCalls[0].File)
	}
}

func TestProcessHookStop(t *testing.T) {
	dir := setupTestDir(t)
	mp := &player.MockPlayer{}
	mn := &notifier.MockNotifier{}
	clock := event.FixedClock{T: fixedTime()}

	ev := &event.HookEvent{
		EventName: "Stop",
		SessionID: "test-session",
	}

	err := processHook(ev, dir, mp, mn, &clock)
	if err != nil {
		t.Fatalf("processHook() error: %v", err)
	}

	if len(mp.PlayCalls) != 1 {
		t.Fatalf("expected 1 play call, got %d", len(mp.PlayCalls))
	}
	if !strings.Contains(mp.PlayCalls[0].File, "Done1.wav") {
		t.Errorf("expected task.complete sound, got %q", mp.PlayCalls[0].File)
	}
}

func TestProcessHookPermissionRequest(t *testing.T) {
	dir := setupTestDir(t)
	mp := &player.MockPlayer{}
	mn := &notifier.MockNotifier{}
	clock := event.FixedClock{T: fixedTime()}

	ev := &event.HookEvent{
		EventName: "PermissionRequest",
		SessionID: "test-session",
	}

	err := processHook(ev, dir, mp, mn, &clock)
	if err != nil {
		t.Fatalf("processHook() error: %v", err)
	}

	if len(mp.PlayCalls) != 1 {
		t.Fatalf("expected 1 play call, got %d", len(mp.PlayCalls))
	}
	if len(mn.SendCalls) != 1 {
		t.Fatalf("expected 1 notification call, got %d", len(mn.SendCalls))
	}
}

func TestProcessHookPaused(t *testing.T) {
	dir := setupTestDir(t)
	writeTestFile(t, filepath.Join(dir, ".paused"), []byte(""))

	mp := &player.MockPlayer{}
	mn := &notifier.MockNotifier{}
	clock := event.FixedClock{T: fixedTime()}

	ev := &event.HookEvent{
		EventName: "Stop",
		SessionID: "test-session",
	}

	err := processHook(ev, dir, mp, mn, &clock)
	if err != nil {
		t.Fatalf("processHook() error: %v", err)
	}

	if len(mp.PlayCalls) != 0 {
		t.Error("should not play when paused")
	}
}

func TestProcessHookDisabledCategory(t *testing.T) {
	dir := setupTestDir(t)
	mp := &player.MockPlayer{}
	mn := &notifier.MockNotifier{}
	clock := event.FixedClock{T: fixedTime()}

	ev := &event.HookEvent{
		EventName: "UserPromptSubmit",
		SessionID: "test-session",
	}

	err := processHook(ev, dir, mp, mn, &clock)
	if err != nil {
		t.Fatalf("processHook() error: %v", err)
	}

	if len(mp.PlayCalls) != 0 {
		t.Error("should not play when category (task.acknowledge) is disabled")
	}
}

func TestProcessHookCompactSuppressed(t *testing.T) {
	dir := setupTestDir(t)
	mp := &player.MockPlayer{}
	mn := &notifier.MockNotifier{}
	clock := event.FixedClock{T: fixedTime()}

	ev := &event.HookEvent{
		EventName: "SessionStart",
		SessionID: "test-session",
		Source:    "compact",
	}

	err := processHook(ev, dir, mp, mn, &clock)
	if err != nil {
		t.Fatalf("processHook() error: %v", err)
	}

	if len(mp.PlayCalls) != 0 {
		t.Error("should not play for compact SessionStart")
	}
}

func TestProcessHookUnknownEvent(t *testing.T) {
	dir := setupTestDir(t)
	mp := &player.MockPlayer{}
	mn := &notifier.MockNotifier{}
	clock := event.FixedClock{T: fixedTime()}

	ev := &event.HookEvent{
		EventName: "SomethingUnknown",
		SessionID: "test-session",
	}

	err := processHook(ev, dir, mp, mn, &clock)
	if err != nil {
		t.Fatalf("processHook() error: %v", err)
	}

	if len(mp.PlayCalls) != 0 {
		t.Error("should not play for unknown event")
	}
}

func TestProcessHookUpdatesState(t *testing.T) {
	dir := setupTestDir(t)
	mp := &player.MockPlayer{}
	mn := &notifier.MockNotifier{}
	clock := event.FixedClock{T: fixedTime()}

	ev := &event.HookEvent{
		EventName: "SessionStart",
		SessionID: "test-session",
	}

	if err := processHook(ev, dir, mp, mn, &clock); err != nil {
		t.Fatalf("processHook() error: %v", err)
	}

	statePath := filepath.Join(dir, ".state.json")
	data, err := os.ReadFile(statePath) //nolint:gosec // test file path
	if err != nil {
		t.Fatalf("failed to read state: %v", err)
	}
	if string(data) == "{}" {
		t.Error("state should be updated after processing")
	}
}
