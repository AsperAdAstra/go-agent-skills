# Risor Runner

Custom Risor interpreter with extended builtins for OpenClaw.

## Binaries

Currently MacOS ARM64

| Binary           | Command                | Purpose           |
| ---------------- | ---------------------- | ----------------- |
| `risor-runner`   | `./risor-runner`       | Scripting runtime |
| `skill-pack`     | `./bin/skill-pack`     | Package skills    |
| `skill-validate` | `./bin/skill-validate` | Validate skills   |

## Quick Start

```bash
# Run a script
./risor-runner 'print(1 + 1)'

# Run with pretty output
./risor-runner -pretty 'print(http_get("https://httpbin.org/get").body)'
```

See [docs/risor.md](./docs/risor.md) for full API reference.

## Build your own binaries

To compile the binaries yourself, you need **Go 1.23+** installed and your `GOBIN`/`GOPATH` set up as usual.

From the repository root, run:

```bash
# Build the main runner
go build -o risor-runner .

# Build the skill packer
go build -o bin/skill-pack ./cmd/pack

# Build the skill validator
go build -o bin/skill-validate ./cmd/validate
```

You can then run the compiled binaries directly from this directory as shown in the **Quick Start** section.
