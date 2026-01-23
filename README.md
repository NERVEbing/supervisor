# supervisor

A CLI tool for tracking project changes across Git repositories, JIRA, and Confluence, generating structured data for AI-powered analysis and reporting.

## Features

### Phase 1: Git Version Diff (✅ Current)

Generate structured JSON diffs between two Git references with:
- **Correct ancestry handling** via automatic merge-base calculation
- Pure Go implementation (no CGO, static binaries)
- Configurable file filtering (suffixes, paths)
- Schema v1.0 compliant JSON output

### Upcoming

- Phase 2: JIRA issue tracking integration
- Phase 3: Multi-source data aggregation and AI-ready reports

---

## Installation

### Build from source

```bash
git clone https://github.com/NERVEbing/supervisor.git
cd supervisor
CGO_ENABLED=0 go build -o bin/supervisor ./cmd/supervisor
```

### Verify installation

```bash
./bin/supervisor diff --help
```

---

## Usage

### Git Version Diff

```bash
supervisor diff --repo /path/to/repo --from v1.0.0 --to v1.1.0
```

#### Flags

- `--repo`: Path to local git repository (default: current directory)
- `--from`: Starting reference (tag/branch/commit) **[required]**
- `--to`: Target reference (tag/branch/commit) **[required]**
- `--exclude-suffix`: File suffixes to exclude (e.g., `--exclude-suffix .png --exclude-suffix .wasm`)
- `--exclude-path`: Path prefixes to exclude (e.g., `--exclude-path vendor/ --exclude-path dist/`)

#### Environment Variables

```bash
export SUPERVISOR_REPO_PATH=/default/repo/path
export SUPERVISOR_EXCLUDE_SUFFIXES=.png,.wasm,.gz
export SUPERVISOR_EXCLUDE_PATHS=vendor/,third_party/
```

Command-line flags override environment variables.

---

## Critical Design Principles

### 1. Merge-Base Ancestry Handling

**The tool AUTOMATICALLY uses merge-base when comparing non-linear histories.**

Example scenario:
```
       C---D---E  (feature branch, tag v1.1.0)
      /
 A---B---F---G    (main branch, tag v1.0.0)
```

When comparing `v1.0.0` to `v1.1.0`:
- Naive diff would include commits F and G (incorrect!)
- **This tool automatically finds merge-base (B) and diffs B→E** (correct!)

This ensures you only see changes actually introduced in the target reference, not unrelated commits.

### 2. Tree Diff is Truth

The tool reports **actual file state changes** (tree diff), NOT accumulated commit messages.
- Reverted commits don't appear in output
- Squashed commits show final result
- Cherry-picks don't cause duplicates

### 3. Local-Only, Read-Only

- No network calls (GitHub API not used in Phase 1)
- No write operations to repository
- Safe to run on production repos

---

## Output Schema

The tool outputs JSON conforming to **Schema v1.0**.

Example output structure:
```json
{
  "schema_version": "1.0",
  "repository": { "name": "...", "vcs": "git" },
  "resolution": {
    "from": { "ref": "v1.0.0", "type": "tag", "commit": "abc123..." },
    "to": { "ref": "v1.1.0", "type": "tag", "commit": "def456..." }
  },
  "baseline": {
    "strategy": "merge-base",
    "base_commit": "xyz789...",
    "ancestry": { "is_linear": false, "relationship": "branched" }
  },
  "tree_diff": {
    "summary": {
      "files": { "added": 3, "modified": 12, "deleted": 1 },
      "lines": { "added": 1210, "deleted": 84, "net": 1126 }
    },
    "files": [ /* detailed file changes */ ]
  }
}
```

---

## Development

### Prerequisites

- Go 1.25+
- golangci-lint (optional)

### Running tests

```bash
go test ./... -v
```

### Linting

```bash
golangci-lint run ./...
```

### Building

```bash
CGO_ENABLED=0 go build -o bin/supervisor ./cmd/supervisor
```

---

## License

MIT License - see [LICENSE](LICENSE)
