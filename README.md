# OpenClaw Workspace

Personal AI assistant workspace running on OpenClaw.

## Project Overview

This workspace hosts the configuration, skills, and tools for an AI assistant (named "Stew") that helps with software engineering, automation, and daily tasks.

## Documentation

- **[Risor Scripting](./docs/risor.md)** — Custom scripting runtime with HTTP, file I/O, JSON, and time functions

## Available Tools

### Binaries

| Binary | Path | Purpose |
|--------|------|---------|
| `risor-runner` | `~/go/bin/risor-runner` | Scripting runtime |
| `risor-gen` | `~/go/bin/risor-gen` | Code generator |
| `skill-pack` | `~/go/bin/skill-pack` | Package skills |
| `skill-validate` | `~/go/bin/skill-validate` | Validate skills |

## Quick Start

### Run a Risor Script

```bash
~/go/bin/risor-runner 'print(1 + 1)'
~/go/bin/risor-runner 'print(http_get("https://httpbin.org/get").status)'
```

See [docs/risor.md](./docs/risor.md) for full API reference.

