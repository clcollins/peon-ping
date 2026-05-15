# peon-ping-toggle

Toggle peon-ping sound notifications on or off.

## Instructions

Run the peon binary with the `toggle` subcommand:

```bash
peon toggle
```

This creates or removes `~/.claude/hooks/peon-ping/.paused`.
When paused, no sounds play for any hook event.

Report the output to the user.
