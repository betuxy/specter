# specter

An interactive TUI for inspecting JSON and YAML in a collapsible tree view.

## Install

```sh
go install github.com/betuxy/specter@latest
```

## Usage

```sh
# Open a file
specter -f file.json

# Pipe from stdin
cat file.json | specter
curl -s https://api.example.com/data | specter

# Start fully expanded
specter -f file.json --expanded

# Explicit subcommand
specter view -f file.json
```

## Commands

| Command | Description |
|---|---|
| `specter` / `specter view` | Open interactive TUI (default) |
| `specter fmt` | Pretty-print JSON or YAML to stdout |
| `specter convert --to yaml` | Convert between JSON and YAML |
| `specter version` | Print version |

### fmt

```sh
specter fmt -f file.json          # pretty-print JSON
specter fmt -f file.yaml          # pretty-print YAML
specter fmt -f file.json -i 4     # 4-space indent
cat file.json | specter fmt
```

### convert

```sh
specter convert -f file.json --to yaml
specter convert -f file.yaml --to json
cat file.yaml | specter convert --to json
```

## Keybindings

| Key | Action |
|---|---|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `space` / `enter` | Expand / collapse node |
| `e` | Expand / collapse entire subtree |
| `q` | Quit |

### Custom keybindings

Create `~/.config/specter/config.yaml` to override any key:

```yaml
keybindings:
  up: "k"
  down: "j"
  expand_collapse: " "
  expand_collapse_all: "e"
  quit: "q"
```

Omitted keys fall back to their defaults.

## Flags

| Flag | Description |
|---|---|
| `-f, --file` | Input file (reads stdin if omitted) |
| `-e, --expanded` | Start with all nodes expanded |
| `-c, --config` | Path to config file |
