---
name: risor-scripting
description: Use when writing or executing risor scripts for OpenClaw automation. Triggers on: "write a risor script", "run a risor script", "create a skill script", "automate with risor", or when asked to perform file operations, HTTP calls, JSON processing, or templating that would benefit from risor's scripting capabilities.
---

# Risor Scripting

## Overview

Enable OpenClaw to write and execute risor scripts for automation tasks. Risor is a embedded scripting language with a namespaced stdlib for HTTP, file I/O, JSON, strings, math, time, crypto, and more.

## Quick Start

```bash
# Run inline script
./risor-runner 'strings.upper("hello")'

# Run script file
./risor-runner -f script.risor

# Pass arguments
./risor-runner -f script.risor name=World city=London

# Clean output (for piping)
./risor-runner --clean 'strings.upper("hello")'
```

## Core Constraints

**Critical: No variable assignment or loops exist in risor.**

```rison
# ❌ DOESN'T WORK - no variables
result = http.get("https://api.example.com")
print(result.body)

# ✅ WORKS - chain directly
print(http.get("https://api.example.com").body)

# ❌ DOESN'T WORK - no loops
for item in items { print(item) }

# ✅ WORKS - use list functions
list.each(items, fn(item) { print(item) })
```

If you need complex control flow, write a script in Go/Python instead.

## Namespaced API

All functions are namespaced. Call as `namespace.function(args)`:

| Namespace | Purpose | Example |
|-----------|---------|---------|
| `strings` | String manipulation | `strings.upper("hi")` |
| `json` | Parse/stringify | `json.parse('{"a":1}').a` |
| `file` | File operations | `file.read("data.json")` |
| `http` | HTTP client | `http.get(url).body` |
| `math` | Math functions | `math.abs(-5)` |
| `time` | Time/date | `time.now()` |
| `crypto` | Hashing | `crypto.md5("data")` |
| `encoding` | URL encoding | `encoding.url_encode("hi there")` |
| `list` | List ops | `list.first([1,2,3])` |
| `sys` | System info | `sys.hostname()` |

Utility globals (no namespace):
- `env_get("KEY")`, `env_set("KEY", "val")`
- `log_info("msg")`, `log_debug("msg")`, `log_warn("msg")`, `log_error("msg")`
- `exec_cmd("ls", "-la").output`
- `template_render("Hello {{name}}", {"name": "World"})`
- `args.name`, `args.city` (from key=value args)
- `skill_validate(frontmatter_string)`

## Common Patterns

### HTTP + JSON pipeline
```rison
data = json.parse(http.get("https://api.example.com/data").body)
print(data.items[0].name)
```

### File read + process + write
```rison
content = file.read("input.txt")
processed = strings.replace(content, "{{version}}", "1.0")
file.write("output.txt", processed)
```

### Template email
```rison
body = template_render("Hi {{name}}, your order #{{order}} is ready.", {"name": "Eli", "order": "12345"})
http.post("https://api.example.com/send", body)
```

### Conditional output
```rison
result = http.get("https://api.example.com/status")
print(result.status == 200, "OK", "FAILED")
```

## Argument Passing

```bash
./risor-runner -f script.risor name=World version=1.0
```

```rison
print("Hello " + args.name + ", version " + args.version)
print(args.argv)  # ["name=World", "version=1.0"]
```

## Output Formats

- Default: `{"result": value}`
- `--clean`: raw value only (for piping)
- `-pretty`: formatted JSON

## Path Security

File paths are relative to working directory. Absolute paths and `..` traversal are blocked for security. Sensitive directories (`/etc`, `/root`, `/home`, `/sys`, `/proc`) are off-limits.

## Testing Scripts

```bash
# Test inline
./risor-runner 'strings.upper("test")'

# Test file
echo 'strings.upper("test")' > /tmp/test.risor
./risor-runner -f /tmp/test.risor

# Check syntax (look for errors)
./risor-runner -f script.risor 2>&1
```

## Project Location

Risor project: `/home/openbot/projects/risor`

Binary: `/home/openbot/projects/risor/risor-runner`

Full API docs: `/home/openbot/projects/risor/docs/risor.md`

## Resources

### references/
- `api-quick-ref.md` - Concise API reference for common operations
- `risor-lang.md` - Risor language basics and constraints
