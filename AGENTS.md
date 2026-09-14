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
├── skills/references/   # API quick ref and language guide
├── bin/                 # Built binaries
└── go.{mod,sum}         # Dependencies (DO NOT MODIFY)
```

## Key Conventions

### Package Layout
- `cmd/<name>` for CLI tools
- `main.go` for primary binary
- `docs/` for reference documentation
- `skills/references/` for skill agent documentation

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

# ✅ WORKS - chain directly
print(http.get("https://api.example.com").body)
```

### No Loops
```rison
# ❌ DOESN'T WORK
for item in items { print(item) }

# ✅ WORKS - use list functions
list.each(items, fn(item) { print(item) })
```

### No try/catch
Errors propagate up — design scripts to fail fast and report errors clearly.

### When NOT to Use Risor
- Complex control flow (nested conditionals, branching logic)
- Stateful operations requiring variable persistence
- Multi-step workflows with error recovery
- Performance-critical code
- If a script gets complex, write it in Go/Python instead

## Data Types

| Type | Examples |
|------|----------|
| `int` | `1`, `42`, `-10` |
| `float` | `3.14`, `-0.5` |
| `string` | `"hello"`, `'world'` |
| `bool` | `true`, `false` |
| `list` | `[1, 2, 3]`, `["a", "b"]` |
| `map` | `{"key": "value", "n": 1}` |
| `null` | `nil` |

## Built-in Functions

```rison
print(value)           # Print to stdout
len(collection)        # Length of string or list
type(value)           # Type name as string
range(n)              # Integer iterator
```

## Operators

```rison
# Arithmetic
1 + 1, 5 - 3, 4 * 2, 10 / 2

# Comparison
1 == 1, 2 != 3, 5 > 3, 5 < 3, 5 >= 5, 3 <= 5

# Logical
true and false, true or false, not true

# String concatenation
"hello" + " " + "world"  # "hello world"
```

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

## Utility Globals (no namespace)

- `env_get("KEY")`, `env_set("KEY", "val")`, `env_vars()`
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

### Template rendering
```rison
body = template_render("Hi {{name}}, your order #{{order}} is ready.", {"name": "Eli", "order": "12345"})
http.post("https://api.example.com/send", body)
```

### Conditional output
```rison
result = http.get("https://api.example.com/status")
print(result.status == 200, "OK", "FAILED")
```

### Environment variables
```rison
home = env_get("HOME")
all_vars = env_vars()
```

### System info
```rison
os_name = sys.os_name()
hostname = sys.hostname()
uuid = sys.uuid()
```

### List operations
```rison
first = list.first([1,2,3])
last = list.last([1,2,3])
unique = list.unique([1,2,2,3])
sorted = list.sort([3,1,2])
```

### JSON/YAML conversion
```rison
json_to_yaml = json.to_yaml({"x": 1})
```

### Math operations
```rison
abs_val = math.abs(-5)
avg = math.avg([1,2,3])
sum_val = math.sum([1,2,3])
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

## Troubleshooting

### Script fails with "assignment" error
Risor doesn't support variable assignment. Chain operations instead:
```rison
# ❌ x = "hello"
# ✅ print(strings.upper("hello"))
```

### Script fails with "loop" error
Risor doesn't support loops. Use list functions:
```rison
# ❌ for item in items { print(item) }
# ✅ list.each(items, fn(item) { print(item) })
```

### Need error handling
Risor has no try/catch. Design scripts to fail fast:
```rison
# Check status before proceeding
status = http.get(url).status
if status != 200 { print("Failed") }
```

## Tips for Writing Risor Scripts

1. **Think in pipelines** — data flows through functions left-to-right or nested
2. **Test inline first** — use `./risor-runner 'expression'` to verify before writing a file
3. **Keep it simple** — if a script gets complex, consider writing it in Go instead
4. **Use template_render** for string interpolation with dynamic data
5. **Access nested data directly** — no need to assign intermediate values
6. **Use `--clean` for piping** — raw output works better in shell pipelines

## Documentation References

- **Full API**: `docs/risor.md`
- **Skill Template**: `skills/SKILL.md`
- **API Quick Ref**: `skills/references/api-quick-ref.md`
- **Language Guide**: `skills/references/risor-lang.md`
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

### String Operations
```rison
strings.upper("hello")
strings.lower("HELLO")
strings.trim("  hi  ")
strings.split("a,b,c", ",")
strings.join(["a", "b"], "-")
strings.replace("hello", "l", "x")
strings.contains("hello", "ell")
```

### Crypto
```rison
crypto.md5("hello")
crypto.sha256("hello")
crypto.base64_enc("hello")
crypto.base64_dec("aGVsbG8=")
```

### Time
```rison
time.now()
time.timestamp()
time.format(time.now(), "2006-01-02")
```

### System
```rison
sys.os_name()
sys.hostname()
sys.uuid()
sys.random_choice(["a","b"])
```

### Encoding
```rison
encoding.url_encode("hello world")
encoding.url_decode("hello+world")
```

## Project Location

Risor project: `/home/openbot/projects/risor`

Binary: `/home/openbot/projects/risor/risor-runner`

Full API docs: `/home/openbot/projects/risor/docs/risor.md`