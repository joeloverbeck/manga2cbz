# MAN2CBZSPE-004: Chapter Types Definition

**Status**: ✅ COMPLETED

## Description
Define the Chapter struct and related types used for chapter representation.

## Files to Touch
- `internal/chapter/chapter.go` (create with types only)

## Pre-Implementation Assumption Check
- ✅ `internal/chapter/` package exists (via `images.go`)
- ✅ `Chapter` struct does not yet exist
- ✅ Spec defines `Chapter` with `Name` and `Path` fields
- ✅ Pattern consistency: `ImageFile.Path` in existing code uses absolute paths

**Note**: The spec says "Full path to directory" without explicitly requiring absolute.
However, for consistency with `ImageFile.Path` (which is explicitly absolute in `images.go`),
and for robustness, this ticket interprets `Path` as absolute.

## Out of Scope
- Chapter discovery logic (separate ticket MAN2CBZSPE-006)
- Image collection (already done in MAN2CBZSPE-003)
- Any traversal or file system operations

## Acceptance Criteria

### Tests That Must Pass
No unit tests required for type definitions, but:
- Code must compile: `go build ./internal/chapter/...`

### Data Structure
```go
// Chapter represents a manga chapter directory
type Chapter struct {
    Name string // Directory name (becomes CBZ filename)
    Path string // Full absolute path to directory
}
```

### Invariants That Must Remain True
- Chapter.Name is the directory basename only (no path separators)
- Chapter.Path is always absolute
- Struct is exported for use by other packages

## Definition of Done
- [x] `Chapter` struct defined in `internal/chapter/chapter.go`
- [x] Package compiles without errors
- [x] Types match spec definition exactly

---

## Outcome

**Completed**: 2026-01-22

### What Was Originally Planned
- Create `internal/chapter/chapter.go` with `Chapter` struct type definition

### What Was Actually Changed
- Created `internal/chapter/chapter.go` with exactly the planned `Chapter` struct
- No deviations from plan

### Files Modified
- `internal/chapter/chapter.go` (created)

### Tests
- No new unit tests added (ticket scope: type definition only)
- Compilation verified: `go build ./internal/chapter/...` passes
- All existing tests continue to pass (29 tests, 95.5%+ coverage)

### Assumptions Validated
- The spec's "Full path" was interpreted as "absolute path" for consistency with `ImageFile.Path`
- This interpretation was documented in the ticket before implementation
