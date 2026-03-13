# Risor Runner - Documentation

Risor is an embedded scripting language for Go. This runner includes custom functions for HTTP, file I/O, and more.

## Binary Location
```
~/go/bin/risor-runner
```

## Usage
```bash
# Run a one-liner
~/go/bin/risor-runner 'print(1 + 1)'

# Run a script file
~/go/bin/risor-runner -f script.risor

# Pretty print JSON output
~/go/bin/risor-runner -pretty '...'
```

## Important Constraints

### No Variable Assignment
You cannot assign to variables. Use direct chaining instead:
```rison
# ❌ Doesn't work
r = http_get("https://example.com")
print(r.body)

# ✅ Works
print(http_get("https://example.com").body)
```

### No Loops
For-loops and list comprehensions are not available. Work with collections directly.

---

## Available Functions

### HTTP Functions

#### `http_get(url string) -> response`
Make a GET request. Returns a response object with:
- `.body` (string) - response body
- `.headers` (map) - response headers
- `.status` (int) - status code

```risor
print(http_get("https://httpbin.org/get").body)
```

#### `http_post(url string, body string) -> response`
Make a POST request with a body string.

```risor
print(http_post("https://httpbin.org/post", "{\"test\":1}").status)
```

**Note:** URL must include `https://` prefix.

---

### File Functions

#### `file_exists(path string) -> bool`
Check if a file exists.

```risor
print(file_exists("/tmp/test.txt"))
```

#### `file_read(path string) -> string`
Read file contents.

```risor
print(file_read("/etc/hostname"))
```

#### `file_write(path string, content string)`
Write content to a file.

```risor
file_write("/tmp/output.txt", "hello world")
print(file_read("/tmp/output.txt"))
```

---

### Environment Functions

#### `env_get(key string) -> string`
Get an environment variable.

```risor
print(env_get("HOME"))
```

---

### JSON Functions

#### `json_parse(str string) -> object`
Parse JSON string to Risor object.

```risor
print(json_parse("{\"a\":1}").a)  # 1
```

#### `json.marshal(obj) -> string`
Marshal Risor object to JSON string.

```risor
print(json.marshal({"a": 1, "b": "hello"}))
```

---

### Time Functions

#### `time.now() -> time`
Get current time.

```risor
print(time.now())  # 2026-03-13T05:16:35Z
```

#### `time.now().format(layout string) -> string`
Format time using Go layout (2006=year, 01=month, 02=day, 15=hour, 04=minute, 05=second).

```risor
print(time.now().format("2006-01-02"))        # 2026-03-13
print(time.now().format("15:04:05"))           # 05:16:35
print(time.now().format("Jan 2, 2006"))        # Mar 13, 2026
```

---

### String Functions (strings module)

#### `strings.contains(s string, substr string) -> bool`

```risor
print(strings.contains("hello world", "world"))  # true
```

#### `strings.split(s string, sep string) -> list`

```risor
print(strings.split("a,b,c", ","))  # ["a", "b", "c"]
```

---

### Math Functions (math module)

```risor
print(math.abs(-5))     # 5
print(math.floor(5.9))  # 5
print(math.ceil(5.1))   # 6
print(math.round(5.5))  # 6
```

---

### Other Functions

#### `print(...)`
Print values to stdout.

```risor
print("hello")
print(123)
print([1, 2, 3])
```

#### `len(obj) -> int`
Get length of string or list.

```risor
print(len("hello"))    # 5
print(len([1, 2, 3]))  # 3
```

#### `type(obj) -> string`
Get type of value.

```risor
print(type(1))         # int
print(type("hello"))   # string
print(type(true))      # bool
print(type(http_get))  # builtin
```

#### `range(n int) -> int_iter`
Create an integer iterator (for use with comprehensions if available).

```risor
print(range(5))  # int_iter(5)
```

---

## Data Types

- **int** - integers (1, 42, -10)
- **float** - floating point (3.14)
- **string** - strings ("hello", 'world')
- **bool** - boolean (true, false)
- **list** - arrays ([1, 2, 3])
- **map** - objects ({"key": "value"})
- **null** - nil value

---

## Operators

```risor
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
true and false  # false
true or false   # true
not true        # false

# String
"hello" + " " + "world"  # "hello world"
```

---

## Examples

### Weather API
```risor
print(json_parse(http_get("https://wttr.in/London?format=j1").body).current_condition[0].temp_C)
```

### GitHub API
```risor
print(http_get("https://api.github.com/repos/risor-io/risor").status)
```

### File Operations
```risor
file_write("/tmp/test.json", "{\"data\": 123}")
print(json_parse(file_read("/tmp/test.json")).data)
```

---

## Risor Generator

There's also a code generator at `~/go/bin/risor-gen` that can generate Risor scripts from descriptions:

```bash
~/go/bin/risor-gen --desc "add two numbers"
~/go/bin/risor-gen --desc "get weather for a city" --city "London"
```

---

## Notes

- The runner outputs both the result and `{"result": null}` - this is normal
- Use `-pretty` flag for pretty-printed JSON output
- All URLs must include the protocol (https://)
