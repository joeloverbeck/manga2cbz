# MAN2CBZSPE-011: Integration Tests

## Description
Create integration tests that verify end-to-end functionality.

## Files to Touch
- `tests/integration_test.go` (create)

## Out of Scope
- Unit tests for individual components (already in component tickets)
- Performance benchmarks
- Fuzzing tests

## Dependencies
- All previous tickets (full implementation required)

## Acceptance Criteria

### Tests That Must Pass
```go
// Integration test cases:

// Test 1: Basic Multiple Chapters
TestIntegration_BasicMultipleChapters:
  Setup: Create TestManga/ with ChapterA/ (3 images), ChapterB/ (2 images)
  Run: manga2cbz TestManga/ --out output/
  Verify:
    - output/ChapterA.cbz exists with 3 images
    - output/ChapterB.cbz exists with 2 images
    - Both archives pass validation

// Test 2: Natural Sort Ordering
TestIntegration_NaturalSortOrdering:
  Setup: Create chapter with 1.png, 2.png, 10.png, 11.png, 3.png
  Run: manga2cbz input/ --out output/
  Verify: CBZ contains files ordered as 1, 2, 3, 10, 11

// Test 3: Mixed File Types
TestIntegration_MixedFileTypes:
  Setup: Chapter with 01.jpg, 02.png, 03.gif, notes.txt
  Run: manga2cbz input/ --out output/
  Verify: CBZ contains only image files, txt excluded

// Test 4: Empty Chapter
TestIntegration_EmptyChapter:
  Setup: Empty folder in input directory
  Run: manga2cbz input/ --out output/
  Verify: Warning logged, no CBZ created, exit code 0

// Test 5: Overwrite Behavior
TestIntegration_OverwriteBehavior:
  Setup: Run once, then run again
  Verify:
    - Without --force: skips existing, exit 0
    - With --force: overwrites, exit 0

// BONUS Test 6: Recursive Mode (not in original spec but useful)
TestIntegration_RecursiveMode:
  Setup:
    root/
      Volume1/Chapter1/
      Volume1/Chapter2/
      Volume2/Chapter3/
  Run: manga2cbz root/ -r --out output/
  Verify: 3 CBZ files created

// BONUS Test 7: Extension Filtering (not in original spec but useful)
TestIntegration_ExtensionFiltering:
  Setup: Chapter with mix of jpg, png, webp, bmp files
  Run: manga2cbz input/ -e jpg,png --out output/
  Verify: Only jpg and png files in CBZ
```

### Test Spec Compliance
| Spec Test | Integration Test |
|-----------|------------------|
| Basic Multiple Chapters | TestIntegration_BasicMultipleChapters |
| Natural Sort Ordering | TestIntegration_NaturalSortOrdering |
| Mixed File Types | TestIntegration_MixedFileTypes |
| Empty Chapter | TestIntegration_EmptyChapter |
| Overwrite Behavior | TestIntegration_OverwriteBehavior |

### Invariants Verified
- **#1 Complete Coverage**: All images included exactly once
- **#2 Correct Ordering**: Natural sort verified
- **#3 Proper Naming**: Folder name -> CBZ name
- **#4 Valid Archive**: unzip -t passes
- **#5 No Internal Folders**: Images at root
- **#6 Robustness**: Spaces, mixed case handled

## Definition of Done
- [x] Integration test file created
- [x] All 5 spec test cases implemented
- [x] 2 bonus tests implemented (recursive mode, extension filtering)
- [x] Tests create/cleanup temporary directories
- [x] `go test -v ./tests/...` passes
- [x] Tests verify all 6 invariants from spec

## Outcome

**Completed**: 2026-01-22

**Implementation**:
- Created `tests/integration_test.go` with 7 integration tests (5 spec + 2 bonus)
- All tests use temporary directories via `t.TempDir()` for automatic cleanup
- Helper functions: `runCLI()`, `createTestImage()`, `createChapterWithImages()`, `getCBZFileList()`, `validateCBZ()`

**Tests**:
1. `TestIntegration_BasicMultipleChapters` - Multiple chapters → CBZ archives (Invariants #1, #3, #4)
2. `TestIntegration_NaturalSortOrdering` - Natural sort: 1, 2, 3, 10, 11 (Invariant #2)
3. `TestIntegration_MixedFileTypes` - Image-only filtering (Invariants #1, #6)
4. `TestIntegration_EmptyChapter` - Warning, no crash, exit 0 (Invariant #6)
5. `TestIntegration_OverwriteBehavior` - Skip without --force, overwrite with --force (Invariant #6)
6. `TestIntegration_RecursiveMode` (BONUS) - Nested directories with -r flag (Invariants #1, #3, #5)
7. `TestIntegration_ExtensionFiltering` (BONUS) - -e flag filters extensions (Invariants #1, #6)

**Note**: Go's flag package requires flags BEFORE positional arguments. All CLI calls use `runCLI(t, "--flag", "value", inputDir)` pattern.
