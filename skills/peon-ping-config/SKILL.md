# peon-ping-config

View or modify peon-ping configuration.

## Instructions

The config file is at `~/.claude/hooks/peon-ping/config.json`.

### View config

```bash
peon status
```

### Change volume

```bash
peon volume 0.3
```

### Switch sound pack

```bash
peon use <pack_name>
```

### Available settings

Edit `~/.claude/hooks/peon-ping/config.json` directly for:

- `enabled` (bool): master on/off
- `default_pack` (string): active sound pack
- `volume` (float 0.0-1.0): playback volume
- `desktop_notifications` (bool): enable notify-send popups
- `categories` (object): per-category toggles
- `annoyed_threshold` (int): spam detection prompt count
- `annoyed_window_seconds` (int): spam detection time window
- `session_start_cooldown_seconds` (int): SessionStart debounce
