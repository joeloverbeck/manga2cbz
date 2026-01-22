# MAN2CBZSPE-006: Chapter Discovery

## Description
Implement chapter directory discovery with flat and recursive modes.

## Files to Touch
- `internal/chapter/chapter.go` (add discovery functions)
- `internal/chapter/chapter_test.go` (create)

## Out of Scope
- Image collection (already done)
- CBZ creation (separate ticket)
- CLI integration
- Verbose/quiet output modes

## Dependencies
- MAN2CBZSPE-004 (Chapter type definition)

## Acceptance Criteria

### Tests That Must Pass
```go
// Test cases for chapter_test.go:

// Flat discovery (one level)
TestDiscover_Flat:
  Setup:
    root/
      ChapterA/ (with images)
      ChapterB/ (with images)
  Input: Discover(root, recursive=false)
  Verify: Returns 2 chapters, names are "ChapterA" and "ChapterB"

// Recursive discovery
TestDiscover_Recursive:
  Setup:
    root/
      Volume1/
        Chapter1/
        Chapter2/
      Volume2/
        Chapter3/
  Input: Discover(root, recursive=true)
  Verify: Returns 3 chapters from nested structure

// Skip files (only directories)
TestDiscover_SkipFiles:
  Setup:
    root/
      Chapter1/ (dir)
      notes.txt (file)
  Input: Discover(root, recursive=false)
  Verify: Only Chapter1 returned, notes.txt ignored

// Empty input directory
TestDiscover_EmptyDir:
  Setup: empty root/
  Input: Discover(root, recursive=false)
  Verify: Returns empty slice, no error

// Directory not found
TestDiscover_NotFound:
  Input: Discover("/nonexistent", recursive=false)
  Verify: Returns error

// Natural sort order of chapters
TestDiscover_NaturalOrder:
  Setup:
    root/
      Chapter 10/
      Chapter 2/
      Chapter 1/
  Input: Discover(root, recursive=false)
  Verify: Order is Chapter 1, Chapter 2, Chapter 10

// Hidden directories skipped
TestDiscover_SkipHidden:
  Setup:
    root/
      Chapter1/
      .hidden/
  Input: Discover(root, recursive=false)
  Verify: Only Chapter1 returned
```

### Function Signatures
```go
func Discover(inputDir string, recursive bool) ([]Chapter, error)
```

### Invariants That Must Remain True
- **Invariant #3**: "Chapter 5" folder -> "Chapter 5.cbz" naming preserved
- **Invariant #6**: Handle spaces in directory names
- Chapters returned in natural sort order
- Only directories are considered (files skipped)
- Hidden directories (starting with .) are skipped
- Chapter.Path is always absolute

## Definition of Done
- [x] `Discover()` function implemented
- [x] Flat mode works (one level deep)
- [x] Recursive mode works (all depths)
- [x] Chapters sorted naturally
- [x] Hidden directories skipped
- [x] All test cases pass
- [x] `go test -race ./internal/chapter/...` passes

## Outcome

**Status**: COMPLETED

**Implementation Summary**:
- Added `Discover(inputDir string, recursive bool) ([]Chapter, error)` - public API
- Added `discoverFlat()` - single-level directory enumeration
- Added `discoverRecursive()` - full-depth traversal using `filepath.WalkDir`
- Added `isLeafDirectory()` - determines if directory is a chapter (no subdirectories)
- Added `sortChapters()` - natural sort helper using `sort.Natural()`

**Test Coverage**: 89.2% (11 chapter discovery tests + 13 image collection tests)

**Tests Added**:
- `TestDiscover_Flat` - basic flat discovery
- `TestDiscover_Recursive` - nested directory traversal
- `TestDiscover_SkipFiles` - ignores non-directories
- `TestDiscover_EmptyDir` - handles empty input
- `TestDiscover_NotFound` - error on missing directory
- `TestDiscover_NaturalOrder` - "Chapter 1, 2, 10" ordering
- `TestDiscover_SkipHidden` - ignores `.hidden` directories
- `TestDiscover_SkipHiddenRecursive` - skips hidden in recursive mode
- `TestDiscover_SpacesInNames` - handles "Chapter 5 - Special"
- `TestDiscover_AbsolutePaths` - verifies absolute paths
- `TestDiscover_RecursiveNaturalOrder` - natural sort in recursive mode

**Key Design Decisions**:
1. Recursive mode treats "leaf" directories (no subdirectories) as chapters
2. Chapter.Name in recursive mode is the relative path from input directory
3. Hidden directories are completely skipped (including all children in recursive mode)
