package player

import "testing"

func TestMockPlayerRecordsCalls(t *testing.T) {
	m := &MockPlayer{}

	if err := m.Play("/path/to/sound.wav", 0.5); err != nil {
		t.Fatalf("Play() error: %v", err)
	}

	if len(m.PlayCalls) != 1 {
		t.Fatalf("expected 1 play call, got %d", len(m.PlayCalls))
	}
	if m.PlayCalls[0].File != "/path/to/sound.wav" {
		t.Errorf("file = %q, want %q", m.PlayCalls[0].File, "/path/to/sound.wav")
	}
	if m.PlayCalls[0].Volume != 0.5 {
		t.Errorf("volume = %f, want 0.5", m.PlayCalls[0].Volume)
	}
}

func TestMockPlayerMultipleCalls(t *testing.T) {
	m := &MockPlayer{}

	m.Play("/a.wav", 0.3)
	m.Play("/b.wav", 0.7)

	if len(m.PlayCalls) != 2 {
		t.Fatalf("expected 2 play calls, got %d", len(m.PlayCalls))
	}
}

func TestMockPlayerKillPrevious(t *testing.T) {
	m := &MockPlayer{}

	if err := m.KillPrevious(); err != nil {
		t.Fatalf("KillPrevious() error: %v", err)
	}
	if m.KillCalls != 1 {
		t.Errorf("expected 1 kill call, got %d", m.KillCalls)
	}
}

func TestMockPlayerConfigurableError(t *testing.T) {
	m := &MockPlayer{PlayErr: errMock}

	if err := m.Play("/a.wav", 0.5); err == nil {
		t.Error("Play() should return configured error")
	}
}

func TestPipeWirePlayerCommandArgs(t *testing.T) {
	p := &PipeWirePlayer{}
	args := p.buildArgs("/path/to/sound.wav", 0.3)

	expected := []string{"--media-role=Notification", "--volume=0.3", "/path/to/sound.wav"}
	if len(args) != len(expected) {
		t.Fatalf("args count = %d, want %d", len(args), len(expected))
	}
	for i, want := range expected {
		if args[i] != want {
			t.Errorf("args[%d] = %q, want %q", i, args[i], want)
		}
	}
}
