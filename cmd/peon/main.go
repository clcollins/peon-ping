package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/clcollins/peon-ping/internal/cli"
	"github.com/clcollins/peon-ping/internal/config"
	"github.com/clcollins/peon-ping/internal/event"
	"github.com/clcollins/peon-ping/internal/notifier"
	"github.com/clcollins/peon-ping/internal/player"
	"github.com/clcollins/peon-ping/internal/sound"
	"github.com/clcollins/peon-ping/internal/state"
)

var version = "dev"

func fixedTime() time.Time {
	return time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)
}

func peonDir() string {
	if d := os.Getenv("PEON_DIR"); d != "" {
		return d
	}
	if d := os.Getenv("CLAUDE_PEON_DIR"); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "hooks", "peon-ping")
}

func processHook(ev *event.HookEvent, dir string, p player.Player, n notifier.Notifier, clock event.Clock) error {
	pausedPath := filepath.Join(dir, ".paused")
	if _, err := os.Stat(pausedPath); err == nil {
		return nil
	}

	category := event.MapEvent(ev)
	if category == event.None {
		return nil
	}

	cfgPath := filepath.Join(dir, "config.json")
	cfg, _ := config.Load(cfgPath)

	if !cfg.Enabled {
		return nil
	}

	if !cfg.IsCategoryEnabled(string(category)) {
		return nil
	}

	statePath := filepath.Join(dir, ".state.json")
	s, _ := state.Load(statePath)

	if category == event.TaskAcknowledge {
		if event.CheckSpam(s, ev.SessionID, cfg.AnnoyedThreshold, cfg.AnnoyedWindowSeconds, clock) {
			category = event.UserSpam
		}
	}

	if ev.EventName == "SessionStart" {
		if event.CheckCooldown(s, ev.SessionID, "SessionStart", cfg.SessionStartCooldownSeconds, clock) {
			s.Save(statePath)
			return nil
		}
	}

	packsDir := filepath.Join(dir, "packs")
	packName := sound.ResolvePack(cfg, s, ev.SessionID, packsDir)
	manifestPath := filepath.Join(packsDir, packName, "manifest.json")

	manifest, err := sound.LoadManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("main: load manifest: %w", err)
	}

	lastPlayed := s.LastPlayed[packName]
	rng := rand.New(rand.NewSource(clock.Now().UnixNano()))

	soundFile, err := sound.PickSound(manifest, string(category), lastPlayed, rng)
	if err != nil {
		s.Save(statePath)
		return nil
	}

	fullPath := filepath.Join(packsDir, packName, soundFile)
	p.Play(fullPath, cfg.Volume)

	s.LastPlayed[packName] = soundFile
	s.Save(statePath)

	if cfg.DesktopNotifications && (category == event.InputRequired || category == event.TaskError) {
		title := "peon-ping"
		message := fmt.Sprintf("[%s] %s", packName, string(category))
		n.Send(title, message, "normal")
	}

	return nil
}

func main() {
	if len(os.Args) > 1 {
		runCLI(os.Args[1:])
		return
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil || len(data) == 0 {
		os.Exit(0)
	}

	var ev event.HookEvent
	if err := json.Unmarshal(data, &ev); err != nil {
		os.Exit(0)
	}

	dir := peonDir()
	p := &player.PipeWirePlayer{}
	n := &notifier.DesktopNotifier{}
	clock := event.RealClock{}

	if err := processHook(&ev, dir, p, n, clock); err != nil {
		fmt.Fprintf(os.Stderr, "peon-ping: %v\n", err)
	}
}

func runCLI(args []string) {
	dir := peonDir()
	cfgPath := filepath.Join(dir, "config.json")
	packsDir := filepath.Join(dir, "packs")

	switch args[0] {
	case "toggle":
		msg, err := cli.Toggle(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "peon: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(msg)

	case "status":
		fmt.Print(cli.Status(cfgPath, dir))

	case "use":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: peon use <pack>")
			os.Exit(1)
		}
		if err := cli.Use(args[1], cfgPath, packsDir); err != nil {
			fmt.Fprintf(os.Stderr, "peon: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Switched to pack: %s\n", args[1])

	case "list":
		packs, _ := cli.List(packsDir)
		for _, p := range packs {
			fmt.Println(p)
		}

	case "volume":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: peon volume <0.0-1.0>")
			os.Exit(1)
		}
		vol, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "peon: invalid volume: %v\n", err)
			os.Exit(1)
		}
		if err := cli.Volume(vol, cfgPath); err != nil {
			fmt.Fprintf(os.Stderr, "peon: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Volume set to %.1f\n", vol)

	case "version":
		fmt.Printf("peon-ping %s\n", version)

	case "help", "--help", "-h":
		fmt.Println(cli.Help())

	default:
		fmt.Fprintf(os.Stderr, "peon: unknown command %q\n", args[0])
		fmt.Fprintln(os.Stderr, cli.Help())
		os.Exit(1)
	}
}
