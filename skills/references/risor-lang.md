# Risor Language Basics

## What Risor Is

Risor is an embedded scripting language for Go. It has a small, clean feature set focused on data manipulation and automation.

## Critical Constraints

### No Variable Assignment
You CANNOT assign to variables:
```rison
# ❌ This fails
x = 5
name = "World"

# ✅ Chain operations instead
print(strings.replace("hello world", "world", "go"))

# ✅ Nested access
data = json.parse(http.get(url).body)
print(data.items[0].name)  # This works because we're chaining on the result
```

### No Loops
There are no for/while loops:
```rison
# ❌ No loops
for i in range(10) { print(i) }
while true { print("x") }

# ✅ Use list functions
list.each(items, fn(item) { print(item) })  # if available

# ✅ Work with whole collections
list.map([1,2,3], fn(x) { x * 2 })  # if available
```

### No try/catch or exceptions
Errors propagate up but there's no try/catch. Design scripts to fail fast and report errors clearly.

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

## Accessing Map/List Values

```rison
# Map access with dot notation
data = {"name": "Eli", "age": 30}
print(data.name)    # "Eli"
print(data["name"]) # "Eli" - bracket notation also works

# List access
items = ["first", "second", "third"]
items[0]    # "first"
items[1]    # "second"
items[-1]   # "third" (last element)
```

## Functions

Risor uses a functional style. Functions are first-class values.

```rison
# Built-in print
print("hello")           # prints to stdout
print(1, 2, "three")     # multiple args

# len works on strings and lists
len("hello")    # 5
len([1,2,3])  # 3

# type returns the type name
type(1)        # "int"
type("hi")     # "string"
```

## Chaining Pattern

Without variables, chain operations:

```rison
# HTTP → JSON → extract field
user = json.parse(http.get("https://api.example.com/user/1").body)
print(user.name)

# File → process → write
content = file.read("input.txt")
upper = strings.upper(content)
file.write("output.txt", upper)

# String operations chain
result = strings.trim(strings.replace("  hello  ", "hello", "world"))
```

## Built-in Functions

```rison
print(value)           # Print to stdout
len(collection)        # Length of string or list
type(value)           # Type name as string
range(n)              # Integer iterator
```

## Tips for Writing Risor Scripts

1. **Think in pipelines** — data flows through functions left-to-right or nested
2. **Test inline first** — use `./risor-runner 'expression'` to verify before writing a file
3. **Keep it simple** — if a script gets complex, consider writing it in Go instead
4. **Use template_render** for string interpolation with dynamic data
5. **Access nested data directly** — no need to assign intermediate values

## Example: Complete Workflow

```rison
# Fetch data from API
response = http.get("https://api.example.com/data")

# Parse JSON
data = json.parse(response.body)

# Extract and transform
name = data.items[0].name
timestamp = time.now()

# Generate output using template
output = template_render("Item: {{name}}\nFetched: {{time}}\n", {
  "name": name,
  "time": timestamp
})

# Write to file
file.write("report.txt", output)

# Log success
log_info("Report generated: " + name)
```
