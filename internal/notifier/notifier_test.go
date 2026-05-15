package notifier

import "testing"

func TestMockNotifierRecordsCalls(t *testing.T) {
	m := &MockNotifier{}

	if err := m.Send("Title", "Message", "normal"); err != nil {
		t.Fatalf("Send() error: %v", err)
	}

	if len(m.SendCalls) != 1 {
		t.Fatalf("expected 1 send call, got %d", len(m.SendCalls))
	}
	if m.SendCalls[0].Title != "Title" {
		t.Errorf("title = %q, want %q", m.SendCalls[0].Title, "Title")
	}
	if m.SendCalls[0].Message != "Message" {
		t.Errorf("message = %q, want %q", m.SendCalls[0].Message, "Message")
	}
	if m.SendCalls[0].Urgency != "normal" {
		t.Errorf("urgency = %q, want %q", m.SendCalls[0].Urgency, "normal")
	}
}

func TestMockNotifierMultipleCalls(t *testing.T) {
	m := &MockNotifier{}

	_ = m.Send("A", "B", "low")
	_ = m.Send("C", "D", "critical")

	if len(m.SendCalls) != 2 {
		t.Fatalf("expected 2 send calls, got %d", len(m.SendCalls))
	}
}

func TestMockNotifierConfigurableError(t *testing.T) {
	m := &MockNotifier{SendErr: errMock}

	if err := m.Send("A", "B", "normal"); err == nil {
		t.Error("Send() should return configured error")
	}
}

func TestDesktopNotifierCommandArgs(t *testing.T) {
	n := &DesktopNotifier{}
	args := n.buildArgs("Peon-Ping", "Task complete", "normal")

	expected := []string{"--urgency=normal", "Peon-Ping", "Task complete"}
	if len(args) != len(expected) {
		t.Fatalf("args count = %d, want %d", len(args), len(expected))
	}
	for i, want := range expected {
		if args[i] != want {
			t.Errorf("args[%d] = %q, want %q", i, args[i], want)
		}
	}
}
