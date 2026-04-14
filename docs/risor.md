# Risor Runner - Documentation

Risor is an embedded scripting language for Go. This runner includes a namespaced standard library for HTTP, file I/O, JSON, and more.

## Binary Location
```
~/go/bin/risor-runner
```

## Usage
```bash
# Run a one-liner
./risor-runner 'strings.upper("hello")'

# Run a script file
./risor-runner -f script.risor

# Run with script arguments (key=value pairs)
./risor-runner -f script.risor name=World

# Pretty print JSON output
./risor-runner -pretty 'math.abs(-5)'

# Clean output (no JSON wrapper)
./risor-runner --clean 'strings.upper("hello")'
```

## Output Formats

By default, output is wrapped: `{"result": value}`

- `-pretty` — pretty-prints the JSON wrapper
- `--clean` — outputs the raw value without wrapper (for piping)

## Important Constraints

### No Variable Assignment
You cannot assign to variables. Use direct chaining:
```rison
# ❌ Doesn't work
result = http.get("https://example.com")
print(result.body)

# ✅ Works
print(http.get("https://example.com").body)
```

### No Loops
For-loops and list comprehensions are not available. Work with collections directly.

---

## Standard Library

All functions are organized into namespaces. Call them as `namespace.function(args)`.

### strings

String manipulation functions.

```rison
strings.upper("hello")              # "HELLO"
strings.lower("HELLO")              # "hello"
strings.trim("  hello  ")          # "hello" (trims spaces)
strings.split("a,b,c", ",")         # ["a", "b", "c"]
strings.join(["a", "b", "c"], "-") # "a-b-c"
strings.replace("hello world", "world", "go")  # "hello go"
strings.contains("hello", "ell")    # true
strings.starts_with("hello", "hel") # true
strings.ends_with("hello", "llo")   # true
strings.regex_match("hello", "ello") # true
strings.regex_replace("hello", "ello", "i")  # "hi"
strings.html_encode("<div>&</div>") # "&lt;div&gt;&amp;&lt;/div&gt;"
```

### json

JSON parsing and serialization.

```rison
json.parse("{\"a\": 1}").a          # 1
json.stringify({"b": 2})             # "{\"b\":2}"
json.to_yaml({"x": 1})              # "x: 1"
```

### file

File system operations. Paths are relative to working directory; absolute paths and `..` traversal are blocked.

```rison
file.exists("main.go")              # true
file.read("main.go")                # file contents as string
file.write("output.txt", "hello")   # true
file.delete("output.txt")           # true
file.list(".")                      # ["file1.go", "file2.go", ...]
file.list_recursive(".")            # [{name, path, is_file, is_dir, size}, ...]
```

### http

HTTP client functions. All return a response object with `.status`, `.body`, `.headers`.

```rison
http.get("https://httpbin.org/get").status        # 200
http.post("https://httpbin.org/post", "data").body
http.put("https://httpbin.org/put", "data").status
http.delete("https://httpbin.org/delete").status
http.headers("https://httpbin.org/get").headers
```

### math

Mathematical functions.

```rison
math.abs(-5)        # 5
math.floor(5.9)    # 5
math.ceil(5.1)     # 6
math.round(5.5)    # 6
math.min(1, 5, 3)  # 1
math.max(1, 5, 3)  # 5
math.sum([1, 2, 3]) # 6
math.avg([1, 2, 3]) # 2
math.random_int(1, 100)  # random int between 1 and 100
```

### time

Time functions.

```rison
time.now()                  # "2026-04-14T05:30:00Z" (RFC3339)
time.timestamp()            # 1776144000 (Unix timestamp)
time.format(time.now(), "2006-01-02")  # "2026-04-14"
time.parse("2026-04-14", "2006-01-02").timestamp  # Unix timestamp
```

### crypto

Hashing and encoding.

```rison
crypto.md5("hello")     # "5d41402abc4b2a76b9719d911017c592"
crypto.sha256("hello")  # "2cf24dba5fb0..."
crypto.base64_enc("hello")  # "aGVsbG8="
crypto.base64_dec("aGVsbG8=")  # "hello"
```

### encoding

URL encoding/decoding.

```rison
encoding.url_encode("hello world")   # "hello+world"
encoding.url_decode("hello+world")   # "hello world"
```

### list

List/collection operations.

```rison
list.first([1, 2, 3])       # 1
list.last([1, 2, 3])        # 3
list.reverse([1, 2, 3])     # [3, 2, 1]
list.unique([1, 2, 2, 3])  # [1, 2, 3]
list.flatten([[1, 2], [3]]) # [1, 2, 3]
list.sort([3, 1, 2])       # [1, 2, 3]
```

### sys

System utilities.

```rison
sys.os_name()      # "linux", "darwin", "windows"
sys.hostname()     # "myserver"
sys.uuid()         # "550e8400-e29b-..."
sys.random_choice(["a", "b", "c"])  # random element
```

---

## Utility Functions (Global)

These are available directly without a namespace prefix.

### env

```rison
env_get("HOME")       # Get env var
env_set("KEY", "val") # Set env var (for subprocesses)
env_vars()            # All env vars as list
env_var("HOME")        # {"value": "...", "exists": true}
```

### log

Structured logging (outputs to stderr with prefixes).

```rison
log_debug("debug message")
log_info("info message")
log_warn("warning message")
log_error("error message")
```

### exec_cmd

Execute shell commands with 30s timeout.

```rison
exec_cmd("ls", "-la").output  # command output
exec_cmd("ls", "-la").error   # error string (empty if success)
```

### template_render

Simple `{{placeholder}}` interpolation.

```rison
template_render("Hello {{name}}", {"name": "World"})  # "Hello World"
```

### args

Script arguments passed as `key=value` pairs.

```rison
# ./risor-runner -f script.risor name=World city=London
args.name    # "World"
args.city    # "London"
args.argv    # ["name=World", "city=London"]
```

### skill_validate

Validate SKILL.md frontmatter format.

```rison
skill_validate(frontmatter_string)  # {"valid": true, "errors": null, "parsed": {...}}
```

---

## Built-in Functions

These are always available in Risor:

```rison
print("hello")           # Print to stdout
len([1, 2, 3])           # 3
type(1)                  # "int"
range(5)                 # int_iter
true and false           # logical operators
1 + 1                    # arithmetic
"hello" + " " + "world"  # string concatenation
```

---

## Data Types

- **int** — integers (1, 42, -10)
- **float** — floating point (3.14)
- **string** — strings ("hello", 'world')
- **bool** — boolean (true, false)
- **list** — arrays ([1, 2, 3])
- **map** — objects ({"key": "value"})
- **null** — nil value

---

## Operators

```rison
# Arithmetic
1 + 1       # 2
5 - 3       # 2
4 * 2       # 8
10 / 2      # 5

# Comparison
1 == 1      # true
2 != 3      # true
5 > 3       # true
5 < 3       # false

# Logical
true and false   # false
true or false    # true
not true         # false

# String
"hello" + " " + "world"  # "hello world"
```

---

## Examples

### Weather API
```rison
json.parse(http.get("https://wttr.in/London?format=j1").body).current_condition[0].temp_C
```

### File + JSON
```rison
file.write("/tmp/test.json", "{\"data\": 123}")
json.parse(file.read("/tmp/test.json")).data
```

### GitHub API
```rison
http.get("https://api.github.com/repos/risor-io/risor").status
```

### Template Email
```rison
template_render("Hi {{name}}, your order #{{order}} is ready.", {"name": "Eli", "order": "12345"})
```

---

## Risor Generator

A code generator at `~/go/bin/risor-gen` can generate Risor scripts from descriptions:

```bash
~/go/bin/risor-gen --desc "add two numbers"
~/go/bin/risor-gen --desc "get weather for a city" --city "London"
```

---

## Notes

- Output format: `{"result": value}` by default, `{"status": "error", "error": "..."}` on error
- Use `--clean` flag for raw value output (useful for piping)
- All URLs must include the protocol (https://)
- File paths are relative to working directory; `..` and absolute paths are blocked
