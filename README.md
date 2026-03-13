# Risor Runner

Custom Risor interpreter with extended builtins for OpenClaw.

## Binaries

| Binary | Command | Purpose |
|--------|---------|---------|
| `risor-runner` | `./risor-runner` | Scripting runtime |
| `skill-pack` | `./bin/skill-pack` | Package skills |
| `skill-validate` | `./bin/skill-validate` | Validate skills |

## Quick Start

```bash
# Run a script
./risor-runner 'print(1 + 1)'

# Run with pretty output
./risor-runner -pretty 'print(http_get("https://httpbin.org/get").body)'
```

See [docs/risor.md](./docs/risor.md) for full API reference.

## Build

```bash
go build -o risor-runner .
go build -o bin/skill-pack ./cmd/pack
go build -o bin/skill-validate ./cmd/validate
```
