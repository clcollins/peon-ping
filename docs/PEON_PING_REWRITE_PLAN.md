# Plan -- Peon-Ping Go Rewrite

- **Status**: In Progress
- **Created**: 2026-05-14
- **Completed**: -
- **Depends On**: toolbox-devtools PR (pipewire-utils)

## Problem

PeonPing/peon-ping is a 6400-line bash script supporting 7 platforms and 17 IDEs.
We need sound notifications for Claude Code on Fedora Toolbox only. The upstream
is unmaintainable for our single use case.

## Proposed Solution

Rewrite as a Go binary (~500-700 lines across packages) with:

- Interface-based DI for testability (Player, Notifier, Clock, Parser)
- Native JSON via encoding/json (no jq/Python dependency)
- Single audio backend: pw-play via PipeWire socket
- Single IDE: Claude Code hooks (with Parser interface for future IDE support)

## Architecture

```text
stdin (Claude Code hook JSON)
       |
       v
  Parser interface  -->  ClaudeCodeParser (converts IDE-specific JSON to HookEvent)
       |
       v
  event.MapEvent()  -->  CESP Category (session.start, task.complete, etc.)
       |
       v
  Suppression checks (spam, cooldown, replay)
       |
       v
  sound.PickSound()  -->  random WAV from pack manifest
       |
       v
  Player interface  -->  PipeWirePlayer (pw-play)
  Notifier interface  -->  DesktopNotifier (notify-send)
       |
       v
  state.Save()  -->  atomic write to .state.json
```

## Implementation Phases

1. Scaffold (go.mod, Makefile, docs, CLAUDE.md, AGENTS.md, README)
2. Config package (TDD)
3. State package (TDD)
4. Event mapping + suppression (TDD)
5. Sound selection (TDD)
6. Player + Notifier interfaces (TDD)
7. CLI subcommands (TDD)
8. Integration + main.go wiring (TDD)
9. Skills, sound packs, final docs

## Pre-Deployment Checklist

- [ ] All go tests pass with race detector
- [ ] golangci-lint passes
- [ ] gofmt clean
- [ ] go vet clean
- [ ] pw-play works from toolbox (requires pipewire-utils)
- [ ] notify-send delivers desktop notifications from toolbox
- [ ] Hooks registered in settings.json fire correctly

## Post-Deployment Verification

- [ ] Start Claude Code session -- session.start sound plays
- [ ] Complete a task -- task.complete sound plays
- [ ] Permission prompt -- input.required sound + notification
- [ ] peon toggle / status / use / list work
- [ ] State persists across sessions

## Success Criteria

- All CESP categories mapped and playing correct sounds
- All tests pass in CI with race detector
- Zero external runtime deps beyond pw-play and notify-send
- go vet, golangci-lint, gofmt all clean
