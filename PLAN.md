# Risor Runner — Improvement Plan

Grounded assessment from static analysis of `main.go`, `cmd/pack/main.go`,
`cmd/validate/main.go`, and the docs (`docs/risor.md`, `skills/SKILL.md`,
`skills/references/*`).

> **Environment note:** the committed binaries (`risor-runner`, `bin/skill-pack`,
> `bin/skill-validate`) are Linux x86-64 ELF and cannot run on this macOS/arm64
> host. Building from source is currently blocked (no network + empty module
> cache), so all findings below are from code + doc inspection, not execution.
> Verification must be done once a build is possible.

---

## Priority 1 — Fix the type-conversion layer (root cause, breaks many functions)

**File:** `main.go`

`toObject` (lines 1106–1135) is the Go-value → Risor-object bridge used by every
builtin return (`wrapFunc`, line 293). It only handles `string`, `float64`,
`bool`, `nil`, `int`, `int64`, `[]interface{}`, and `map[string]interface{}`.
Any other concrete Go type falls into the `default` branch (line 1132) and is
turned into a **string** via `fmt.Sprintf("%v", v)`.

That silently corrupts return values from several builtins:

| Function | Line | Returns | Actual result |
|---|---|---|---|
| `split` | 772–774 | `[]string` | stringified `[a b c]` |
| `fileList` | 452–469 | `[]string` | stringified |
| `fileListRecursive` | 473–510 | `[]map[string]interface{}` | stringified |
| `envVars` | 542–544 | `[]string` | stringified |
| `httpGet`/`httpHeaders` `.headers` | 331, 388 | `http.Header` (`map[string][]string`) | stringified |
| `skillValidate` `.parsed` | 676, 690 | `map[string]string` | stringified |

Consequence: chaining breaks. `list.first(strings.split("a,b", ","))` panics
because `first` (line 828) does `args[0].([]interface{})` on what is now a string.
`http.get(url).headers` and `file.list(".")` are unusable as structured data.

**Fix:**
1. Extend `toObject` to handle `[]string`, `[]map[string]interface{}`,
   `map[string]string`, `map[string]any`/`map[string][]string` (or convert
   `http.Header` before returning it from the HTTP funcs).
2. Better: centralize conversion so every builtin returns normalized Go types
   (`[]interface{}`, `map[string]interface{}`) and `toObject` needs only the
   canonical set.
3. Add a unit test that round-trips each builtin's return type through
   `toObject` → `toGoValue` and asserts the type survives.

**Related:** `toGoValue` (lines 1073–1104) also has a narrow type set and falls
back to `result.Inspect()` (line 1102) for `object.ERROR`, `object.BYTES`, etc.
Handle `ERROR` explicitly (return its message) to avoid leaking `Inspect()`
strings into results.

---

## Priority 2 — Doc / implementation mismatches

These are cases where the documented API does not match the code (or vice versa).

1. **`time.format`** — `docs/risor.md:134` and `api-quick-ref.md:64` show
   `time.format(time.now(), "2006-01-02")`, but `formatTime` (lines 965–971)
   expects a numeric Unix timestamp (`toFloat(args[0])`). Passing the RFC3339
   string from `time.now()` yields `toFloat == 0` → formats 1970-01-01.
   → Fix code to accept both a timestamp and an RFC3339 string, or fix the docs.

2. **`time.parse`** — `docs/risor.md:135` shows `time.parse(...).timestamp`, but
   `parseTime` (lines 973–983) returns the bare `int64` timestamp (line 982),
   not an object with a `.timestamp` field.
   → Align docs to the real return type, or return a map.

3. **`list.each` / `list.map`** — referenced in `AGENTS.md`, `skills/SKILL.md:44`,
   and `skills/references/risor-lang.md:32,35` (as "if available"), but the
   `list` module (lines 163–170) only registers `first`, `last`, `reverse`,
   `unique`, `flatten`, `sort`.
   → Either implement them (requires closure/`fn(...)` support in the Risor
   env — confirm it exists) or delete the references.

4. **`--file` flag** — `CHANGELOG.md:12` claims a `--file` flag, but the code
   only registers `-f` (line 42). `docs/risor.md` and `README.md` correctly use
   `-f`.
   → Fix the changelog, or add `--file` as an alias.

5. **`math.round/floor/ceil/abs`** return `float64` (`math.Round` etc.), so
   `math.abs(-5)` prints `5` rather than `5` — cosmetic but confusing next to
   docs that show integer results. Decide int vs float and be consistent.

---

## Priority 3 — `skill-pack` and `skill-validate` bugs

**File:** `cmd/pack/main.go`

1. **Dead `skipDirs` logic.** Directories are returned early at lines 65–67
   (`if info.IsDir() { return nil }`), so the `skipDirs` loop (lines 70–78)
   never sees a directory — its `filepath.SkipDir` branch is unreachable.
   Result: files inside `.git`, `node_modules`, `.npm`, `__pycache__` **are**
   packed. The `if info.IsDir()` checks at lines 91 and 100 are also dead code.
   → Skip directory *subtrees* (return `filepath.SkipDir` on the directory
   entry, before the early return).

2. **"Skip hidden files" is not implemented.** The comment at line 58 promises
   skipping hidden files (except `.env.example`), but there is no such logic.
   → Add it, or remove the comment.

**File:** `cmd/validate/main.go`

3. **`requires` is never parsed.** `MetadataRequires` (lines 50–53) is declared
   but `parseMetadata` (lines 119–145) never reads `requires.bins` / `requires.env`.
   → Populate it from YAML, and emit warnings/errors for missing declared bins/env.

4. **Two divergent frontmatter parsers.** The `skill-validate` binary uses
   `gopkg.in/yaml.v3` (line 130), while the in-runner `skill_validate` builtin
   (lines 637–721) hand-rolls a naive line splitter. They can disagree on the
   same input.
   → Consolidate on the YAML parser (shared helper) so validation is identical.

5. **`strings.SplitN(content, "---", 3)`** (line 124) breaks if `---` appears
   inside a description/body. Prefer splitting only on the leading delimiter.

---

## Priority 4 — Robustness & cleanup

1. **Unchecked type assertions that panic** instead of returning an error:
   `join` (777), `first` (828), `last` (836), `reverseList` (844), `unique`
   (853), `flatten` (867), `sortList` (883), `sumVals` (920), `avgVals` (929),
   `randomChoice` (1033), `fileList` (455), `fileListRecursive` (476),
   `formatTime` (968).
   → Add safe accessors / bounds + type checks that return a clear error.

2. **`rand.Seed` still used** (lines 1028, 1037) despite `CHANGELOG.md` claiming
   it was fixed, and it's deprecated + global-state. Use a local
   `rand.New(rand.NewSource(...))` (or `math/rand/v2`), which is also safe for
   the builtin's 30s `exec_cmd` timeout context.

3. **`jsonToYaml`** (lines 739–747) assumes a top-level object; it silently
   fails for arrays/scalars and its hand-rolled `toYaml` (749–769) emits
   non-compliant YAML for nested maps. Consider delegating to `yaml.v3`
   (already an indirect dep) or document the limitation.

4. **`jsonStringify`/`jsonParse`** (724–737) and other funcs index `args[0]`
   without bounds checks — same class of panic risk as item 1.

5. **`sys.random_int` duplicates `math.random_int`** (both `randomInt`, lines 177
   and 136). Fine to keep, but `docs/risor.md` only documents `math.random_int`;
   decide where it canonically lives.

---

## Suggested sequencing

1. **P1** type-conversion layer + round-trip tests (unblocks the most behavior).
2. **P3** `skill-pack` directory-skip fix (small, high user impact).
3. **P2** doc/impl alignment (`time.format`, `time.parse`, `list.each/map`).
4. **P4** panic hardening + `rand` cleanup.

## Verification

- `go test ./...` — add tests for `toObject`/`toGoValue` round-trips and the
  previously-panicking list/string paths.
- `go vet ./...`
- Manual smoke tests once a build is possible:
  - `./risor-runner 'list.first(strings.split("a,b,c", ","))'` → `a`
  - `./risor-runner --clean 'http.get("https://example.com").status'`
  - `./bin/skill-pack -path <dir>` → verify `node_modules` is excluded
  - `./bin/skill-validate -path <dir> -json` → verify `requires` populated

## Out of scope / blocked

- Rebuilding binaries and any execution-based verification are blocked until a
  build environment (network + module cache) is available.
- `go.mod` / `go.sum` are not to be modified per project convention.
