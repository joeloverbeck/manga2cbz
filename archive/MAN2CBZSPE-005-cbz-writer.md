# MAN2CBZSPE-005: CBZ Writer Implementation

## Description
Implement CBZ (ZIP) archive creation with proper streaming and error handling.

## Files to Touch
- `internal/cbz/writer.go` (create)
- `internal/cbz/writer_test.go` (create)

## Out of Scope
- Chapter discovery (separate ticket)
- Image collection (separate ticket)
- CLI integration
- Progress reporting/verbose output

## Dependencies
- MAN2CBZSPE-003 (ImageFile type from chapter package)

## Acceptance Criteria

### Tests That Must Pass
```go
// Test cases for writer_test.go:

// Basic CBZ creation
TestCreate_Basic:
  Setup: 3 temp image files
  Input: Create("output.cbz", images, opts)
  Verify:
    - File exists at output.cbz
    - unzip -t output.cbz passes
    - Archive contains exactly 3 entries

// Images at root (no nested folders)
TestCreate_FlatStructure:
  Setup: images with various source paths
  Input: Create("output.cbz", images, opts)
  Verify: All entries in archive are at root level (no "/" in entry names except filename)

// Store method (no compression)
TestCreate_StoreMethod:
  Setup: image files
  Input: Create("output.cbz", images, opts)
  Verify: ZIP entries use Store method (compression method 0)

// Cleanup on error
TestCreate_CleanupOnError:
  Setup: Make one image file unreadable after opening
  Input: Create("output.cbz", images, opts)
  Verify: Partial output.cbz file is removed on error

// Overwrite behavior with force=false
TestCreate_NoOverwrite:
  Setup: existing output.cbz file
  Input: Create("output.cbz", images, CreateOptions{Force: false})
  Verify: Returns error, existing file unchanged

// Overwrite behavior with force=true
TestCreate_ForceOverwrite:
  Setup: existing output.cbz file
  Input: Create("output.cbz", images, CreateOptions{Force: true})
  Verify: File is overwritten with new content

// Validate function
TestValidate_ValidArchive:
  Setup: Create valid CBZ
  Input: Validate("output.cbz")
  Verify: Returns nil

// Validate function with corrupt archive
TestValidate_CorruptArchive:
  Setup: Create file with invalid ZIP content
  Input: Validate("corrupt.cbz")
  Verify: Returns error
```

### Function Signatures
```go
type CreateOptions struct {
    Force bool // Overwrite existing files
}

func Create(outputPath string, images []chapter.ImageFile, opts CreateOptions) error
func Validate(cbzPath string) error
```

### Invariants That Must Remain True
- **Invariant #4**: `unzip -t *.cbz` passes for all outputs
- **Invariant #5**: Images at archive root, not nested
- Use `zip.Store` method (no compression)
- Stream files via `io.Copy` (no full memory loading)
- Clean up partial files on error

## Definition of Done
- [x] `Create()` function implemented with streaming
- [x] `Validate()` function implemented
- [x] ZIP uses Store method (no compression)
- [x] Partial files cleaned up on error
- [x] Force flag controls overwrite behavior
- [x] All test cases pass
- [x] `go test -race ./internal/cbz/...` passes

## Outcome

**Status**: ✅ Completed

**Implementation Details**:
- Created `internal/cbz/writer.go` with `Create()` and `Validate()` functions
- Created `internal/cbz/writer_test.go` with 11 test cases covering all acceptance criteria

**Test Results**:
- All 11 CBZ package tests pass
- 88.9% code coverage for cbz package
- All 40 project tests pass with race detection

**Key Design Decisions**:
1. Uses `zip.Store` method (no compression) since images are already compressed
2. Streams files via `io.Copy` to avoid loading entire images into memory
3. Uses deferred cleanup with success flag to remove partial files on error
4. Images stored at archive root using only `Name` field from ImageFile
5. Force flag checked before file creation to avoid unnecessary work

**Files Created**:
- `internal/cbz/writer.go` (~100 lines)
- `internal/cbz/writer_test.go` (~200 lines)
