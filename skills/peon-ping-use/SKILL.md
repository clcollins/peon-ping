# peon-ping-use

Switch the active peon-ping sound pack.

## Instructions

List available packs:

```bash
peon list
```

Switch to a pack:

```bash
peon use <pack_name>
```

This updates `default_pack` in `~/.claude/hooks/peon-ping/config.json`.
The change takes effect on the next hook event.
