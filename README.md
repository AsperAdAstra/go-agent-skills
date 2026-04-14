# Risor Runner

Custom Risor interpreter with a namespaced stdlib for OpenClaw skill scripting.

## Binaries

| Binary | Command | Purpose |
|--------|---------|---------|
| `risor-runner` | `./risor-runner` | Scripting runtime |
| `skill-pack` | `./bin/skill-pack` | Package skills |
| `skill-validate` | `./bin/skill-validate` | Validate skills |

## Quick Start

```bash
# Run a script
./risor-runner 'strings.upper("hello")'

# Run a script file
./risor-runner -f script.risor

# Script arguments
./risor-runner -f script.risor name=World

# Pretty output
./risor-runner -pretty 'json.stringify({"a": 1})'

# Clean output (for piping)
./risor-runner --clean 'strings.upper("hello")'
```

## Namespaced Stdlib

```rison
strings.upper("hello")      # "HELLO"
json.parse('{"a": 1}').a    # 1
file.exists("main.go")      # true
http.get("https://...").body
math.abs(-5)                # 5
time.now()                   # "2026-04-14T..."
crypto.md5("hello")
list.first([1, 2, 3])       # 1
sys.hostname()              # "myserver"
```

See [docs/risor.md](./docs/risor.md) for full API reference.

## Build

```bash
go build -o risor-runner .
go build -o bin/skill-pack ./cmd/pack
go build -o bin/skill-validate ./cmd/validate
```
