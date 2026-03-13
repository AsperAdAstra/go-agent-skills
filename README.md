# OpenClaw Workspace

Personal AI assistant workspace running on OpenClaw.

## Project Overview

This workspace hosts the configuration, skills, and tools for an AI assistant (named "Stew") that helps with software engineering, automation, and daily tasks.

## Documentation

- **[Risor Scripting](./docs/risor.md)** — Custom scripting runtime with HTTP, file I/O, JSON, and time functions

## Available Tools

### Skills (Installed)

| Skill | Description |
|-------|-------------|
| `notion` | Notion API integration |
| `weather` | Weather via wttr.in |
| `tmux` | Terminal session control |
| `blogwatcher` | RSS/blog monitoring |
| `clawhub` | Skill marketplace |
| `skill-creator` | Create/audit skills |
| `healthcheck` | Security hardening |

### Binaries

| Binary | Path | Purpose |
|--------|------|---------|
| `risor-runner` | `~/go/bin/risor-runner` | Scripting runtime |
| `risor-gen` | `~/go/bin/risor-gen` | Code generator |
| `blogwatcher` | `~/go/bin/blogwatcher` | RSS monitor |
| `skill-pack` | `~/go/bin/skill-pack` | Package skills |
| `skill-validate` | `~/go/bin/skill-validate` | Validate skills |

## Quick Start

### Run a Risor Script

```bash
~/go/bin/risor-runner 'print(1 + 1)'
~/go/bin/risor-runner 'print(http_get("https://httpbin.org/get").status)'
```

See [docs/risor.md](./docs/risor.md) for full API reference.

### Check Weather

```bash
curl -s "wttr.in/London?format=%l:+%c+%t"
```

### Use Skills

Skills are automatically loaded based on context. Common skills:
- `openclaw weather "London"` — Get weather
- `openclaw skills check` — Audit installed skills

## Workspace Structure

```
.openclaw/
├── AGENTS.md          # Agent instructions
├── SOUL.md            # Persona definition
├── USER.md            # User profile
├── TOOLS.md           # Tool configuration
├── HEARTBEAT.md       # Periodic tasks
├── docs/              # Documentation
│   └── risor.md       # Risor API reference
├── memory/            # Daily notes
└── skills/            # Installed skills
```

## Configuration

- **Model**: minimax-portal/MiniMax-M2.5
- **Channel**: Telegram
- **Timezone**: UTC

## Resources

- [OpenClaw Docs](https://docs.openclaw.ai)
- [GitHub](https://github.com/openclaw/openclaw)
- [Discord](https://discord.com/invite/clawd)
- [ClawHub](https://clawhub.com) — Skill marketplace
