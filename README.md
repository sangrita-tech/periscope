# 🔭 Periscope

Periscope prints a text snapshot of a local path or a GitHub/GitLab repository.

## Install

```bash
go install github.com/sangrita-tech/periscope/cmd/periscope@latest
```

## Usage

```bash
periscope <path> -t -c -i "vendor" -i "*.log"
```

`path` can be a file, a directory, a path inside the current repository, or a GitHub/GitLab repository URL.

Flags:

- `-t, --tree` prints only the file tree.
- `-c, --copy` copies the result to the clipboard.
- `-i, --ignore` ignores a file or directory pattern. This flag can be repeated.

## Examples

```bash
periscope .
```

```bash
periscope . -t
```

```bash
periscope https://github.com/sangrita-tech/periscope -i ".git" -i "coverage"
```

Output file blocks start with `[FILE]` and use paths relative to the directory where the command is called.
