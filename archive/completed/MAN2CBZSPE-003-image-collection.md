# MAN2CBZSPE-003: Image File Collection

## Description
Implement image file filtering and collection from a chapter directory.

## Files to Touch
- `internal/chapter/images.go` (create)
- `internal/chapter/images_test.go` (create)

## Out of Scope
- Chapter discovery/directory traversal (separate ticket)
- CBZ creation (separate ticket)
- Natural sort implementation (dependency, already done)
- CLI argument parsing

## Dependencies
- MAN2CBZSPE-002 (natural sort) - uses `sort.Natural()` for ordering

## Assumptions Verified (2026-01-22)
- ✅ `internal/sort/natural.go` exists with `Natural(items []string)` function
- ✅ Go module name is `manga2cbz` (import as `manga2cbz/internal/sort`)
- ✅ `internal/chapter/` directory exists (contains only `.gitkeep`)
- **Design Note**: Spec defines separate `SortImages(images []ImageFile)` function, but this
  ticket integrates sorting directly into `CollectImages` for simplicity. This is acceptable
  as the spec's `SortImages` can be added later as a public wrapper if needed externally.

## Acceptance Criteria

### Tests That Must Pass
```go
// Test cases for images_test.go:

// Basic image collection with default extensions
TestCollectImages_BasicJPG:
  Setup: temp dir with 01.jpg, 02.jpg, 03.jpg
  Input: CollectImages(dir, []string{"jpg", "jpeg", "png"})
  Verify: Returns 3 ImageFile structs, paths are absolute, names are basenames

// Mixed image types
TestCollectImages_MixedTypes:
  Setup: temp dir with 01.jpg, 02.png, 03.gif
  Input: CollectImages(dir, []string{"jpg", "png", "gif"})
  Verify: All 3 files collected

// Filter non-images
TestCollectImages_FilterNonImages:
  Setup: temp dir with 01.jpg, notes.txt, readme.md
  Input: CollectImages(dir, []string{"jpg"})
  Verify: Only 01.jpg returned, txt and md excluded

// Case insensitive extensions
TestCollectImages_CaseInsensitive:
  Setup: temp dir with 01.JPG, 02.Png, 03.jpeg
  Input: CollectImages(dir, []string{"jpg", "png", "jpeg"})
  Verify: All 3 files collected regardless of case

// Natural sort ordering
TestCollectImages_NaturalOrder:
  Setup: temp dir with 10.jpg, 2.jpg, 1.jpg
  Input: CollectImages(dir, []string{"jpg"})
  Verify: Returned order is 1.jpg, 2.jpg, 10.jpg

// Empty directory
TestCollectImages_EmptyDir:
  Setup: empty temp dir
  Input: CollectImages(dir, []string{"jpg"})
  Verify: Returns empty slice, no error

// Directory not found
TestCollectImages_NotFound:
  Input: CollectImages("/nonexistent", []string{"jpg"})
  Verify: Returns error
```

### Data Structure
```go
type ImageFile struct {
    Path string // Full absolute path to file
    Name string // Base filename (for archive entry)
}
```

### Invariants That Must Remain True
- **Invariant #1**: Every image file matching extensions is included exactly once
- **Invariant #2**: Images are returned in natural sort order
- **Invariant #6**: Handle mixed case extensions robustly
- ImageFile.Path is always absolute
- ImageFile.Name is always just the basename (no directory)

## Definition of Done
- [x] `ImageFile` struct defined
- [x] `CollectImages(dir string, extensions []string) ([]ImageFile, error)` implemented
- [x] Images are filtered by extension (case-insensitive)
- [x] Images are sorted using natural sort
- [x] All test cases pass
- [x] `go test -race ./internal/chapter/...` passes

---

## Status: COMPLETED (2026-01-22)

## Outcome

### What Was Actually Changed vs Originally Planned

**Originally Planned:**
- Create `internal/chapter/images.go` with `ImageFile` struct and `CollectImages` function
- Create `internal/chapter/images_test.go` with 7 test cases

**Actually Implemented:**
- ✅ Created `internal/chapter/images.go` (57 lines) - as planned
- ✅ Created `internal/chapter/images_test.go` (179 lines) - extended beyond plan

### New/Modified Tests

| Test | Rationale |
|------|-----------|
| `TestCollectImages_BasicJPG` | Core ticket requirement - validates basic functionality |
| `TestCollectImages_MixedTypes` | Core ticket requirement - validates multi-extension support |
| `TestCollectImages_FilterNonImages` | Core ticket requirement - validates filtering behavior |
| `TestCollectImages_CaseInsensitive` | Core ticket requirement - validates Invariant #6 |
| `TestCollectImages_NaturalOrder` | Core ticket requirement - validates Invariant #2 |
| `TestCollectImages_EmptyDir` | Core ticket requirement - validates empty directory handling |
| `TestCollectImages_NotFound` | Core ticket requirement - validates error handling |
| `TestCollectImages_ExtensionsWithDots` | **Extra** - Edge case: extensions passed with leading dots (e.g., ".jpg") |
| `TestCollectImages_SkipsSubdirectories` | **Extra** - Validates flat (non-recursive) collection behavior |
| `TestCollectImages_NoMatchingExtensions` | **Extra** - Edge case: no files match any extension |
| `TestCollectImages_NaturalOrderWithPrefix` | **Extra** - Real manga naming pattern (page1.jpg, page2.jpg, page10.jpg) |
| `TestCollectImages_PathFileConsistency` | **Extra** - Validates Path/Name invariant consistency |

### Design Decisions

1. **Sorting integrated into CollectImages**: The spec mentioned a separate `SortImages` function, but the ticket specified that `CollectImages` should return sorted results. This approach is simpler and ensures callers always get correctly ordered images.

2. **Extension normalization**: Extensions are normalized (lowercase, dot-stripped) to handle user input variations robustly.

3. **Absolute paths only**: `filepath.Abs()` is called on the directory to guarantee `ImageFile.Path` is always absolute.

### Test Coverage

- **Coverage**: 95.5% of statements in `internal/chapter`
- **Race detection**: Passed with `-race` flag
- **All 12 tests pass**
