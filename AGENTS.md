# Peon-Ping Agent

Sound notifications for Claude Code hooks via PipeWire (pw-play)
on Fedora Toolbox. Go binary with interface-based dependency injection.

## Role

Assist with developing and maintaining the peon-ping Go binary
and its test suite.

## Scope

- cmd/peon/ (entry point)
- internal/ (all packages: config, state, event, sound, player, notifier, cli)
- packs/ (vendored sound files -- assets, not generated)
- skills/ (Claude Code skill definitions)
- docs/ (plan documents)

## Capabilities

- Modify Go source and tests
- Update documentation and plans
- Add new sound packs
- Debug audio/notification issues

## Boundaries

- Do not execute the peon binary's play functions outside of tests
- Do not modify ~/.claude/settings.json hooks directly (document for user)
- Do not add audio players beyond pw-play
- Sound files in packs/ are vendored assets; do not regenerate
- Never push to upstream PeonPing/peon-ping

## Conventions

- Language: Go (module github.com/clcollins/peon-ping)
- Go version: match go.mod
- Testing: go test with race detector, table-driven tests, interface mocking
- TDD: write tests first, then implement
- JSON: encoding/json (no jq, no Python)
- Audio: pw-play via Player interface
- Notifications: notify-send via Notifier interface
- Container engine: podman, not docker
- File format: Unix line endings only
- Error handling: wrap with context (fmt.Errorf("pkg: op: %w", err))
- Logging: log/slog structured logging

## Plan/Document/Review Cycle

Same 5-phase workflow as cluster-config and gort:

1. Research -- explore codebase, understand problem
2. Plan -- write docs/PLAN.md with problem, solution, architecture
3. Review -- validate plan against conventions
4. Implement -- TDD, tests first
5. Update Documentation -- append PMR notes, lessons learned

Plans in docs/ are never overwritten or deleted -- only appended with
post-mortem notes, lessons learned, or lint fixes.
