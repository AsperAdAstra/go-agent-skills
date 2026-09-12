# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Custom Risor scripting runtime ("risor-runner") with 114 extended built-in functions, built for OpenClaw. Includes tooling for packaging and validating reusable "skills" (script bundles with SKILL.md metadata).

## Build

```bash
go build -o risor-runner .
go build -o bin/skill-pack ./cmd/pack
go build -o bin/skill-validate ./cmd/validate
```

No Makefile — use `go build` directly. Pre-compiled binaries exist at `risor-runner`, `bin/skill-pack`, and `bin/skill-validate`.

## Testing

No test files exist yet. Standard `go test ./...` would be the convention.

## Architecture

**main.go** — Single-file runtime (~710 lines). Registers 114 built-in functions via `risor.WithGlobal()`, then evaluates a script passed as a CLI argument. All output is JSON (`{"result": ...}` or `{"status": "error", "error": "..."}`). Flag `-pretty` enables formatted output.

Built-in function categories: HTTP (get/post/put/delete), File I/O, Exec/Env, JSON/YAML, String (including regex), List operations, Math, Time, Crypto/Hashing, Encoding (URL/HTML/Base64), System info, Random/UUID.

Each built-in converts between Go native types and Risor `object.Object` types using helpers like `goToRisor()` and type assertion on `args`.

**cmd/pack/main.go** — Packages a skill directory into a zip, skipping `.git`, `node_modules`, etc.

**cmd/validate/main.go** — Validates skill structure: checks for required `SKILL.md` with YAML frontmatter metadata.

**docs/risor.md** — Full API reference for all built-in functions.

## Key Constraints

- Go 1.23.0 with a `replace` directive in go.mod mapping `deepnoodle-ai/risor/v2` to `risor-io/risor v1.8.1`
- Risor scripts have no variable assignment or loops — uses functional/chaining style
- Module: `risor-runner` (not a library, standalone binary)
