# Risor API Quick Reference

## strings
```rison
strings.upper("hello")              # "HELLO"
strings.lower("HELLO")              # "hello"
strings.trim("  hi  ")             # "hi"
strings.split("a,b,c", ",")         # ["a", "b", "c"]
strings.join(["a", "b"], "-")       # "a-b"
strings.replace("hello", "l", "x")  # "hexxo"
strings.contains("hello", "ell")    # true
strings.starts_with("hello", "he")  # true
strings.ends_with("hello", "lo")    # true
strings.regex_match("hello", "ello") # true
strings.regex_replace("hi", "i", "ello")  # "hello"
strings.html_encode("<div>")         # "&lt;div&gt;"
```

## json
```rison
json.parse('{"a": 1}').a            # 1
json.stringify({"b": 2})             # "{\"b\":2}"
json.to_yaml({"x": 1})              # "x: 1"
```

## file
```rison
file.exists("main.go")              # true/false
file.read("main.go")                # string contents
file.write("out.txt", "data")       # true
file.delete("out.txt")              # true
file.list(".")                      # ["file1", "file2"]
file.list_recursive(".")             # [{name, path, is_file, is_dir, size}, ...]
```

## http
```rison
http.get(url).status                # 200
http.get(url).body                  # response body
http.get(url).headers               # response headers
http.post(url, body).status        # 201
http.put(url, body).status         # 200
http.delete(url).status             # 204
http.headers(url).headers           # HEAD request
```

## math
```rison
math.abs(-5)        # 5
math.floor(5.9)  # 5
math.ceil(5.1)   # 6
math.round(5.5)  # 6
math.min(1, 5, 3)  # 1
math.max(1, 5, 3)  # 5
math.sum([1,2,3])  # 6
math.avg([1,2,3])  # 2
math.random_int(1, 100)  # random int
```

## time
```rison
time.now()                  # "2026-04-14T09:00:00Z"
time.timestamp()            # 1776144000
time.format(time.now(), "2006-01-02")  # "2026-04-14"
```

## crypto
```rison
crypto.md5("hello")     # "5d41402abc4b2a76b..."
crypto.sha256("hello")  # "2cf24dba5fb0a..."
crypto.base64_enc("hello")  # "aGVsbG8="
crypto.base64_dec("aGVsbG8=")  # "hello"
```

## encoding
```rison
encoding.url_encode("hello world")   # "hello+world"
encoding.url_decode("hello+world")   # "hello world"
```

## list
```rison
list.first([1,2,3])       # 1
list.last([1,2,3])        # 3
list.reverse([1,2,3])     # [3,2,1]
list.unique([1,2,2,3])   # [1,2,3]
list.flatten([[1,2],[3]]) # [1,2,3]
list.sort([3,1,2])       # [1,2,3]
```

## sys
```rison
sys.os_name()        # "linux"
sys.hostname()       # "myserver"
sys.uuid()           # "550e8400-e29b-..."
sys.random_choice(["a","b"])  # random element
```

## Utility (no namespace)
```rison
env_get("HOME")
env_set("KEY", "value")
env_vars()
log_info("msg")
exec_cmd("ls").output
template_render("Hi {{name}}", {"name": "X"})
args.name
skill_validate(frontmatter_string)
```
