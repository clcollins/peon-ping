# peon-ping (personal fork)

Warcraft III Peon voice notifications for Claude Code hooks.
Stripped-down, Linux-only Go rewrite with full Fedora Toolbox support.

> **Attribution**: This is a personal fork of
> [PeonPing/peon-ping](https://github.com/PeonPing/peon-ping) by Tony Sheng.
> Original project licensed under MIT.
> This fork removes all non-Linux platforms, non-Claude-Code IDE integrations,
> and Python/bash dependencies in favor of a minimal Go binary targeting
> Linux with PipeWire audio. Fully tested and supported inside Fedora Toolbox
> containers, but works on any Linux system with PipeWire.

## What it does

Plays sound effects and sends desktop notifications when Claude Code reaches
certain states: task complete, needs input, errors, session start, compacting,
and more. Sound packs are downloaded separately from the upstream
[PeonPing/og-packs](https://github.com/PeonPing/og-packs) repository and
played via PipeWire (`pw-play`). When running inside a Fedora Toolbox
container, audio routes through the host's PipeWire socket which is
automatically bind-mounted into the toolbox.

## Requirements

- Linux with PipeWire (works directly on host or inside Fedora Toolbox)
- `pw-play` (from `pipewire-utils`)
- `notify-send` (from `libnotify`)
- Go 1.23+ (for building from source)

### Fedora Toolbox setup (optional)

If running inside a Fedora Toolbox container, `pipewire-utils` must be
installed in the toolbox image to provide `pw-play`. The host's PipeWire
socket is automatically bind-mounted, so audio plays through the host's
speakers with no extra configuration.

If using
[toolbox-devtools](https://github.com/clcollins/toolbox-devtools), this is
included ([PR #7](https://github.com/clcollins/toolbox-devtools/pull/7)).
For manual installation inside a running toolbox:

```bash
sudo dnf install pipewire-utils
```

## Installation

Build and install the binary:

```bash
make install
```

This builds `bin/peon` and copies it to `~/.local/bin/peon`.

Set up the peon-ping data directory and copy config:

```bash
mkdir -p ~/.claude/hooks/peon-ping
cp config.json ~/.claude/hooks/peon-ping/config.json
```

Download sound packs from the upstream
[PeonPing/og-packs](https://github.com/PeonPing/og-packs) repository
(licensed CC-BY-NC-4.0 by their respective authors):

```bash
git clone --depth=1 --branch=v1.1.0 \
  https://github.com/PeonPing/og-packs.git /tmp/og-packs

for pack in peon peasant sc_kerrigan sc_battlecruiser glados; do
  cp -r "/tmp/og-packs/$pack" ~/.claude/hooks/peon-ping/packs/"$pack"
  mv ~/.claude/hooks/peon-ping/packs/"$pack"/openpeon.json \
     ~/.claude/hooks/peon-ping/packs/"$pack"/manifest.json
done

rm -rf /tmp/og-packs
```

Register Claude Code hooks in `~/.claude/settings.json`. Add the following
entries to the `hooks` object. All events use async execution except
`SessionStart` (sync for immediate feedback):

```json
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "/home/<user>/.local/bin/peon",
            "timeout": 10
          }
        ]
      }
    ],
    "Stop": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "/home/<user>/.local/bin/peon",
            "timeout": 10,
            "async": true
          }
        ]
      }
    ],
    "Notification": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "/home/<user>/.local/bin/peon",
            "timeout": 10,
            "async": true
          }
        ]
      }
    ],
    "PermissionRequest": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "/home/<user>/.local/bin/peon",
            "timeout": 10,
            "async": true
          }
        ]
      }
    ],
    "UserPromptSubmit": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "/home/<user>/.local/bin/peon",
            "timeout": 10,
            "async": true
          }
        ]
      }
    ],
    "PostToolUseFailure": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "/home/<user>/.local/bin/peon",
            "timeout": 10,
            "async": true
          }
        ]
      }
    ],
    "PreCompact": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "/home/<user>/.local/bin/peon",
            "timeout": 10,
            "async": true
          }
        ]
      }
    ],
    "SessionEnd": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "/home/<user>/.local/bin/peon",
            "timeout": 10,
            "async": true
          }
        ]
      }
    ]
  }
}
```

Replace `<user>` with your username.

## Uninstallation

Remove hook entries from `~/.claude/settings.json` (delete each event key that
references the peon binary).

```bash
rm ~/.local/bin/peon
rm -rf ~/.claude/hooks/peon-ping/
```

## Usage

### Automatic (via hooks)

Once installed, sounds play automatically when Claude Code fires hook events.
No manual intervention needed.

### CLI commands

```bash
peon toggle          # Toggle sounds on/off
peon status          # Show current configuration
peon use <pack>      # Switch sound pack
peon list            # List available packs
peon volume <0-1>    # Set volume (0.0 to 1.0)
peon help            # Show usage
```

## Sound Packs

Sound packs are downloaded from
[PeonPing/og-packs](https://github.com/PeonPing/og-packs) (CC-BY-NC-4.0).
They are not included in this repository. Default packs:

| Pack | Description |
|------|-------------|
| peon | Warcraft III Orc Peon |
| peasant | Warcraft III Human Peasant |
| sc_kerrigan | StarCraft Kerrigan |
| sc_battlecruiser | StarCraft Battlecruiser |
| glados | Portal GLaDOS |

## Configuration

Edit `~/.claude/hooks/peon-ping/config.json`:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `enabled` | bool | `true` | Master on/off |
| `default_pack` | string | `"peon"` | Active sound pack |
| `volume` | float | `0.5` | Volume (0.0-1.0) |
| `desktop_notifications` | bool | `true` | Enable notify-send |
| `categories` | object | see below | Per-category toggles |
| `annoyed_threshold` | int | `3` | Spam detection threshold |
| `annoyed_window_seconds` | int | `10` | Spam detection window |
| `session_start_cooldown_seconds` | int | `30` | SessionStart debounce |

### Categories

| Category | Default | Trigger |
|----------|---------|---------|
| `session.start` | on | New session starts |
| `task.acknowledge` | off | User submits a prompt |
| `task.complete` | on | Task finishes |
| `task.error` | on | Bash tool errors |
| `input.required` | on | Permission or input needed |
| `resource.limit` | on | Context compaction |
| `user.spam` | on | Rapid prompt spam |

## Development

```bash
make test            # Run tests with race detector
make test-verbose    # Verbose test output
make cover           # Tests with coverage profile
make lint            # Run golangci-lint
make fmt             # Check gofmt
make vet             # Run go vet
make build           # Build binary
make ci              # Run all CI checks locally
make ci-all          # Run CI inside container (podman)
make help            # Show all targets
```

## License

MIT (see [LICENSE](LICENSE))

---

*100% Vibe Coded* 🤖
