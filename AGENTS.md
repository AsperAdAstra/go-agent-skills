# Risor Runner - Agent Context

A custom Risor interpreter with a namespaced stdlib for OpenClaw skill scripting.

## Project Overview

- **Module**: `risor-runner`
- **Go version**: 1.25
- **Dependencies**: risor/v2, uuid, wonton, yaml.v3
- **Binaries**: `risor-runner`, `skill-pack`, `skill-validate`
- **Primary use**: Scripting automation with HTTP, file I/O, JSON, and more

## Build Commands

```bash
go build -o risor-runner .
go build -o bin/skill-pack ./cmd/pack
go build -o bin/skill-validate ./cmd/validate
```

## Test Commands

```bash
go test ./...
go vet ./...
```

## Repository Structure

```
.
├── main.go              # risor-runner entry point
├── cmd/
│   ├── pack/main.go     # skill-pack builder
│   └── validate/main.go # skill-validate checker
├── docs/risor.md        # Full API reference
├── skills/SKILL.md      # Skill template format
├── bin/                 # Built binaries
└── go.{mod,sum}         # Dependencies (DO NOT MODIFY)
```

## Key Conventions

### Package Layout
- `cmd/<name>` for CLI tools
- `main.go` for primary binary
- `docs/` for reference documentation

### CLI Flags
- `-f` — Run script file
- `-pretty` — Pretty-print JSON output
- `--clean` — Raw value output (no wrapper)
- `--help` — Usage information

### Security Profiles
- File paths are relative to working directory
- Absolute paths and `..` traversal are blocked
- Sensitive directories (`/etc`, `/root`, `/home`, `/sys`, `/proc`) are off-limits
- HTTP requests require explicit protocol (`https://`)

### JSON Output
- Default: `{"result": value}`
- Use `--clean` for raw value output
- Use `-pretty` for formatted JSON

### Test Patterns
- `go test ./...` — Run all tests
- `go vet ./...` — Static analysis
- Tests should be in `<package>_test.go` files

## What to Avoid

- **DO NOT** modify `go.mod` — dependencies are managed externally
- **DO NOT** edit `go.sum` — checksums are auto-generated
- **DO NOT** hardcode paths — use working directory relative paths
- **DO NOT** bypass security profile guards
- **DO NOT** vendoring — use module dependencies

## Risor Scripting Constraints

### No Variable Assignment
```rison
# ❌ DOESN'T WORK
result = http.get("https://api.example.com")
print(result.body)

# ✅ WORKS
print(http.get("https://api.example.com").body)
```

### No Loops
```rison
# ❌ DOESN'T WORK
for item in items { print(item) }

# ✅ WORKS
list.each(items, fn(item) { print(item) })
```

## Documentation References

- **Full API**: `docs/risor.md`
- **Skill Template**: `skills/SKILL.md`
- **Publish Templates**: `Blog.publish`, `LinkedIn.publish`, `X.publish`

## Common Operations

### HTTP + JSON
```rison
json.parse(http.get("https://api.example.com/data").body).items[0].name
```

### File Processing
```rison
content = file.read("input.txt")
processed = strings.replace(content, "{{version}}", "1.0")
file.write("output.txt", processed)
```

### Environment
```rison
env_get("HOME")
env_vars()
```

### Logging
```rison
log_debug("msg")
log_info("msg")
log_warn("msg")
log_error("msg")
```