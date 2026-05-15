package player

import (
	"errors"
	"fmt"
	"os/exec"
)

type Player interface {
	Play(soundFile string, volume float64) error
	KillPrevious() error
}

type PipeWirePlayer struct {
	cmd *exec.Cmd
}

func (p *PipeWirePlayer) buildArgs(soundFile string, volume float64) []string {
	return []string{
		"--media-role=Notification",
		fmt.Sprintf("--volume=%.1f", volume),
		soundFile,
	}
}

func (p *PipeWirePlayer) Play(soundFile string, volume float64) error {
	_ = p.KillPrevious()

	args := p.buildArgs(soundFile, volume)
	p.cmd = exec.Command("pw-play", args...) //nolint:gosec // args are constructed internally
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("player: pw-play: %w", err)
	}

	go func() { _ = p.cmd.Wait() }()
	return nil
}

func (p *PipeWirePlayer) KillPrevious() error {
	if p.cmd != nil && p.cmd.Process != nil {
		_ = p.cmd.Process.Kill()
		p.cmd = nil
	}
	return nil
}

var errMock = errors.New("mock error")

type PlayCall struct {
	File   string
	Volume float64
}

type MockPlayer struct {
	PlayCalls []PlayCall
	KillCalls int
	PlayErr   error
}

func (m *MockPlayer) Play(soundFile string, volume float64) error {
	m.PlayCalls = append(m.PlayCalls, PlayCall{File: soundFile, Volume: volume})
	return m.PlayErr
}

func (m *MockPlayer) KillPrevious() error {
	m.KillCalls++
	return nil
}
