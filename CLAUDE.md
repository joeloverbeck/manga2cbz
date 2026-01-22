# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

manga2cbz is a Go CLI tool that converts manga chapter folders into CBZ (Comic Book ZIP) archives for VR reading applications. It processes a directory containing chapter subfolders, each with image files, and outputs properly sorted CBZ archives.

## Environment Setup

**Go executable location**: `/usr/local/go/bin/go`

The `go` command is NOT in the default PATH in this environment. Always use the full path:

```bash
# Instead of: go test ./...
# Use:
/usr/local/go/bin/go test ./...

# Instead of: go build ./cmd/manga2cbz
# Use:
/usr/local/go/bin/go build ./cmd/manga2cbz
```

All build commands in this document assume you prefix with `/usr/local/go/bin/` or have added it to your PATH.

## Build Commands

```bash
# Development build
go build -o manga2cbz ./cmd/manga2cbz

# Run all tests with race detection
go test -v -race -cover ./...

# Run tests for a specific package
go test -v ./internal/sort/...

# Production build (optimized)
go build -trimpath -ldflags "-s -w" -o manga2cbz ./cmd/manga2cbz

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o manga2cbz.exe ./cmd/manga2cbz
```

## Quick Reference for Non-Go Developers

| JavaScript/npm | Go Equivalent |
|---------------|---------------|
| `npm test` / `npm run test:unit` | `go test ./...` |
| `npm run build` | `go build ./cmd/manga2cbz` |
| `package.json` | `go.mod` |
| `node_modules/` | Go modules cached in `$GOPATH/pkg/mod` |
| `src/` | `cmd/` (executables), `internal/` (libraries) |

## Architecture

```
cmd/manga2cbz/main.go     # CLI entry point, flag parsing, orchestration
internal/
  sort/natural.go         # Natural sort: "page2" < "page10" (not lexicographic)
  chapter/
    chapter.go            # Chapter directory discovery (flat or recursive)
    images.go             # Image file collection and filtering by extension
  cbz/writer.go           # CBZ archive creation (ZIP with Store method)
```

**Data Flow**: Input directory → Discover chapters → Collect images per chapter → Natural sort → Create CBZ

## Chapter Discovery (internal/chapter/chapter.go)

**Flat Mode (default)**:
```bash
manga2cbz /path/to/manga
```
- Discovers only direct subdirectories
- Input: `/manga/` with `Chapter1/`, `Chapter2/`
- Output: Processes `Chapter1/`, `Chapter2/`

**Recursive Mode (`-r` flag)**:
```bash
manga2cbz /path/to/manga -r
```
- Discovers chapters at any depth
- Input: `/manga/` with `Vol1/Chapter1/`, `Vol2/Chapter2/`
- Output: Processes both nested chapters

**Implementation**: Two internal functions:
- `discoverFlat()` - single ReadDir call
- `discoverRecursive()` - filepath.Walk for full tree

## Natural Sort (internal/sort/natural.go)

**Problem**: Standard `sort.Strings()` gives lexicographic order:
```
["page1", "page10", "page2"]  // WRONG: 1, 10, 2
```

**Solution**: Natural sort compares numeric chunks numerically:
```
["page1", "page2", "page10"]  // CORRECT: 1, 2, 10
```

**Algorithm**:
1. Split string into chunks: `["page", "1"]`, `["page", "10"]`
2. Compare chunk-by-chunk
3. If both chunks are numeric, compare as integers
4. Otherwise, compare as strings

**Why It Matters**: Manga page ordering is critical for readability.
`page9.jpg` MUST come before `page10.jpg`.

## Key Implementation Details

- **CBZ files** use ZIP Store method (no compression) since images are already compressed
- **Images stored at archive root** (no nested folders inside CBZ)
- **Stream files** via `io.Copy` to avoid loading entire images into memory
- **Default image extensions**: jpg, jpeg, png, gif, bmp, webp

## Go Language Gotchas (Project-Specific)

### Extension Matching is Case-Sensitive
```go
// filepath.Ext() preserves case - ".JPG" != ".jpg"
ext := strings.ToLower(filepath.Ext(filename))  // Always normalize
```

### Path Handling Must Use filepath Package
```go
// CORRECT - works on Windows and Unix
path := filepath.Join(dir, chapter, image)

// WRONG - breaks on Windows (uses "/" instead of "\")
path := dir + "/" + chapter + "/" + image
```

### File Cleanup with Defer
```go
file, err := os.Open(path)
if err != nil {
    return err
}
defer file.Close()  // Ensures cleanup even on error
```

## Test Patterns

**Test File Location**: Same directory as code, suffix `_test.go`
```
internal/sort/natural.go       # Implementation
internal/sort/natural_test.go  # Tests
```

**Table-Driven Tests** (preferred pattern for this project):
```go
func TestNaturalSort(t *testing.T) {
    tests := []struct {
        name     string
        input    []string
        expected []string
    }{
        {"empty", []string{}, []string{}},
        {"numeric", []string{"2", "10", "1"}, []string{"1", "2", "10"}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := NaturalSort(tt.input)
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Exit Codes

- `0`: All chapters processed successfully
- `1`: Partial failure (some chapters failed)
- `2`: Total failure (invalid args, no chapters found)

## Specification and Tickets

The full technical specification is in `specs/manga2cbz-spec.md`. Implementation tickets are in `tickets/` with prefix `MAN2CBZSPE-*`, defining acceptance criteria and test cases for each component.

## Session Workflow

**Before implementing any ticket:**
1. If the user has told you to rely on a file or files for reference, read them for context
2. Read the specific ticket provided
3. Verify ticket assumptions against current codebase
4. Correct ticket if assumptions are outdated
5. Implement the minimal required changes
6. Run `go test -v -race -cover ./...`
7. Archive completed ticket to `archive/` with Outcome section
