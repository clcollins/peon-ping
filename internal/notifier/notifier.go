package notifier

import (
	"errors"
	"fmt"
	"os/exec"
)

type Notifier interface {
	Send(title, message, urgency string) error
}

type DesktopNotifier struct{}

func (n *DesktopNotifier) buildArgs(title, message, urgency string) []string {
	return []string{
		fmt.Sprintf("--urgency=%s", urgency),
		title,
		message,
	}
}

func (n *DesktopNotifier) Send(title, message, urgency string) error {
	args := n.buildArgs(title, message, urgency)
	cmd := exec.Command("notify-send", args...) //nolint:gosec // args are constructed internally
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("notifier: notify-send: %w", err)
	}
	return nil
}

var errMock = errors.New("mock error")

type SendCall struct {
	Title   string
	Message string
	Urgency string
}

type MockNotifier struct {
	SendCalls []SendCall
	SendErr   error
}

func (m *MockNotifier) Send(title, message, urgency string) error {
	m.SendCalls = append(m.SendCalls, SendCall{Title: title, Message: message, Urgency: urgency})
	return m.SendErr
}
