# Periscope 🔭

Periscope is a small CLI tool for quick directories inspections: it recursively walks through a folder and either prints its structure or outputs the contents of all text files in a clean, readable format.

## Why Periscope

- Quickly understand what's inside a project without opening files manually.
- Inspect Git repositories.
- Produce a clean, consolidated text snapshot of a codebase — handy for reviews, documentation, sharing, or feeding into an LLM.
- Strip comments or mask domains in URLs to avoid leaking sensitive information.
- Exclude unwanted folders or files (e.g., `vendor/`, `node_modules/`, build artifacts, logs).

## Installation

### Using Go

```bash
go install github.com/sangrita-tech/periscope@latest
```

### Downloading a release

Prebuilt binaries for macOS, Linux, and Windows are available in the Releases section of the repository.
Just download the appropriate archive, unpack it, and place periscope in your PATH.

## Quick Start

### Show directory structure

```bash
periscope tree dir [path]
```

### Show contents of all text files

```bash
periscope view dir [path]
```

`[path]` is optional — defaults to the current directory.

## Flags

### Common flags (tree and view)

- `-c, --copy` — copy the result to clipboard
- `-i, --ignore-path` "mask" — ignore paths matching a mask (\* supported), can be used multiple times
- `-I, --ignore-content` "mask" — skip a file if any line inside matches the mask

### View-only flags

- `-z, --strip-comments` — remove comment-only lines (//, #, --, ;)
- `-m, --mask-url` — mask domain names inside URLs (replaced with deterministic fake domains)

## Examples

Ignore vendor and node_modules:

```bash
periscope view dir . -i "vendor" -i "node_modules"
```

Strip comments and mask URLs:

```bash
periscope view dir -z -m
```

Copy a directory tree directly to clipboard:

```bash
periscope tree dir src -c
```

## Updating

If installed via Go:

```bash
go install github.com/sangrita-tech/periscope@latest
```

If downloaded from Releases, simply replace the binary with a newer one.
