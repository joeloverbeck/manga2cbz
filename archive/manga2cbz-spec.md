# manga2cbz - Technical Specification

## Overview

A Go command-line tool that converts manga chapter folders into CBZ (Comic Book ZIP) archives for VR reading applications.

**Language**: Go (for single-binary distribution)
**Target Platforms**: Linux, Windows, macOS

---

## Project Structure

```
manga2cbz/
├── cmd/manga2cbz/
│   └── main.go              # CLI entry point, flag parsing, orchestration
├── internal/
│   ├── cbz/
│   │   └── writer.go        # CBZ archive creation
│   ├── chapter/
│   │   ├── chapter.go       # Chapter directory discovery
│   │   └── images.go        # Image file filtering and collection
│   └── sort/
│       └── natural.go       # Natural/alphanumeric sorting
├── go.mod
├── Makefile
└── tests/
    └── integration_test.go
```

---

## CLI Interface

```
manga2cbz [OPTIONS] <INPUT_DIR>

ARGUMENTS:
    <INPUT_DIR>    Directory containing chapter subfolders

OPTIONS:
    -o, --out <DIR>       Output directory (default: INPUT_DIR)
    -e, --ext <EXTS>      Image extensions (default: jpg,jpeg,png,gif,bmp,webp)
    -f, --force           Overwrite existing CBZ files
    -r, --recursive       Process nested directory structures
    -v, --verbose         Show detailed progress
    -q, --quiet           Suppress non-error output
    -h, --help            Show help
    --version             Show version
```

### Examples
```bash
manga2cbz /path/to/MangaDir
manga2cbz /path/to/MangaDir --out /path/to/output
manga2cbz /path/to/MangaDir -e jpg,png --force -v
```

---

## Core Components

### 1. Chapter Discovery (`internal/chapter/chapter.go`)

**Responsibility**: Find chapter directories in input path

**Functions**:
- `Discover(inputDir string, recursive bool) ([]Chapter, error)`
- `discoverFlat(inputDir string) ([]Chapter, error)` - one level deep
- `discoverRecursive(inputDir string) ([]Chapter, error)` - all depths

**Data Structure**:
```go
type Chapter struct {
    Name string // Directory name (becomes CBZ filename)
    Path string // Full path to directory
}
```

### 2. Image Collection (`internal/chapter/images.go`)

**Responsibility**: Find and filter image files in a chapter directory

**Functions**:
- `CollectImages(dir string, extensions []string) ([]ImageFile, error)`
- `SortImages(images []ImageFile)` - uses natural sort

**Data Structure**:
```go
type ImageFile struct {
    Path string // Full path to file
    Name string // Base filename (for archive entry)
}
```

### 3. Natural Sort (`internal/sort/natural.go`)

**Responsibility**: Sort strings so "10" comes after "9", not after "1"

**Algorithm**:
1. Parse string into alternating numeric/non-numeric chunks
2. Compare chunks: numeric chunks by value, string chunks lexicographically
3. Handle mixed cases like "page10.png" vs "page2.png"

**Function**:
- `Natural(items []string)` - in-place sort

### 4. CBZ Writer (`internal/cbz/writer.go`)

**Responsibility**: Create valid CBZ (ZIP) archives

**Functions**:
- `Create(outputPath string, images []ImageFile, opts CreateOptions) error`
- `Validate(cbzPath string) error` - verify archive integrity

**Key Behaviors**:
- Images stored at archive root (no nested folders)
- Use `zip.Store` method (no compression - images already compressed)
- Stream files via `io.Copy` (no full memory loading)
- Clean up partial files on error

---

## Data Flow

```
Input Directory
      │
      ▼
┌─────────────────────┐
│ 1. Discover Chapters│  → []Chapter
└─────────────────────┘
      │
      ▼ (for each chapter)
┌─────────────────────┐
│ 2. Collect Images   │  → []ImageFile (filtered by extension)
└─────────────────────┘
      │
      ▼
┌─────────────────────┐
│ 3. Natural Sort     │  → []ImageFile (ordered)
└─────────────────────┘
      │
      ▼
┌─────────────────────┐
│ 4. Create CBZ       │  → ChapterName.cbz
└─────────────────────┘
```

---

## Invariants (Must Always Hold)

| # | Invariant | Description |
|---|-----------|-------------|
| 1 | **Complete Coverage** | Every image file included exactly once in CBZ |
| 2 | **Correct Ordering** | Natural sort: 1, 2, 3, ..., 9, 10, 11 (not 1, 10, 11, 2...) |
| 3 | **Proper Naming** | "Chapter 5" folder → "Chapter 5.cbz" |
| 4 | **Valid Archive** | `unzip -t *.cbz` passes for all outputs |
| 5 | **No Internal Folders** | Images at archive root, not nested |
| 6 | **Robustness** | Handle spaces, mixed case extensions, edge cases |

---

## Error Handling

| Category | Behavior | Example |
|----------|----------|---------|
| **Fatal** | Exit immediately | Invalid input directory |
| **Chapter Error** | Log and continue | Write failure on one chapter |
| **Warning** | Log and skip | Empty chapter folder |

### Exit Codes
- `0`: Success (all chapters processed)
- `1`: Partial failure (some chapters failed)
- `2`: Total failure (invalid args, no chapters)

---

## Test Cases

### 1. Basic Multiple Chapters
- Input: TestManga/ with ChapterA/ (3 images), ChapterB/ (2 images)
- Verify: Both CBZ files created, correct image counts, valid archives

### 2. Natural Sort Ordering
- Input: Chapter with files 1.png, 2.png, 10.png, 11.png, 3.png
- Verify: CBZ contains files ordered as 1, 2, 3, 10, 11

### 3. Mixed File Types
- Input: Chapter with 01.jpg, 02.png, 03.gif, notes.txt
- Verify: CBZ contains only image files, txt excluded

### 4. Empty Chapter
- Input: Empty folder in input directory
- Verify: Warning logged, no CBZ created, no crash

### 5. Overwrite Behavior
- Run twice on same input
- Verify: Without --force skips/warns, with --force overwrites

### 6. Large Batch (Performance)
- Input: 50+ chapters with 100+ images each
- Verify: Completes without memory issues, all outputs valid

---

## Build Commands

```makefile
# Development build
build:
    go build -o manga2cbz ./cmd/manga2cbz

# Production build (smaller binary)
release:
    go build -trimpath -ldflags "-s -w" -o manga2cbz ./cmd/manga2cbz

# Cross-compile for Windows
windows:
    GOOS=windows GOARCH=amd64 go build -o manga2cbz.exe ./cmd/manga2cbz

# Run tests
test:
    go test -v -race -cover ./...
```

---

## Implementation Phases

### Phase 1: Foundation
1. Initialize Go module
2. Implement `internal/sort/natural.go` + tests
3. Implement `internal/chapter/images.go` + tests

### Phase 2: Core Logic
4. Implement `internal/cbz/writer.go` + tests
5. Implement `internal/chapter/chapter.go` + tests

### Phase 3: CLI Integration
6. Implement `cmd/manga2cbz/main.go`
7. Write integration tests

### Phase 4: Polish
8. Add verbose/quiet modes
9. Implement --recursive flag
10. Create Makefile

---

## Verification

After implementation, verify using the test data at `source/isekai-craft-gurashi/` (82 chapters):

```bash
# Build
go build -o manga2cbz ./cmd/manga2cbz

# Test on real data
./manga2cbz ./source/isekai-craft-gurashi --out ./archive -v

# Validate outputs
for f in ./archive/*.cbz; do unzip -t "$f"; done

# Check a sample CBZ
unzip -l "./archive/Chapter 1.cbz"
```
