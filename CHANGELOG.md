# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [2026-03-13] - 2026-03-13

### Added

- Risor scripting runtime with 114 extended built-in functions (HTTP, File I/O, String, List, Math, Time, Crypto, System, JSON)
- `--file` flag for reading scripts from files instead of CLI arguments
- JSON output by default for LLM compatibility (`{"result": ...}` / `{"status": "error", ...}`)
- `-pretty` flag for formatted output
- `skill-validate` tool for validating skill structure (SKILL.md, scripts/, metadata)
- `skill-pack` tool for packaging skills into distributable zip files
- Random/ID functions: `uuid`, `random_int`, `random_choice`
- Encoding functions: `url_encode`, `url_decode`, `html_encode`
- Full API reference documentation (`docs/risor.md`)
- Pre-compiled binaries for `risor-runner`, `skill-pack`, and `skill-validate`

### Fixed

- Script argument parsing
- Nil pointer dereference in HTTP functions
- `rand.Seed` deprecation warning
- `html_encode` to properly escape HTML entities
- Removed ignored errors and dead code

### Changed

- Migrated to Risor v2
- Removed unused code generator (`cmd/gen`)
