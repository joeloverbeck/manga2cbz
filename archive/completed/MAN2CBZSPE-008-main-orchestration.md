# MAN2CBZSPE-008: Main Orchestration and Error Handling

## Description
Implement the main processing loop that ties all components together.

## Files to Touch
- `cmd/manga2cbz/main.go` (extend with orchestration logic)

## Out of Scope
- Verbose/quiet output formatting (separate ticket)
- Makefile (separate ticket)
- Integration tests (separate ticket)

## Dependencies
- MAN2CBZSPE-006 (chapter discovery)
- MAN2CBZSPE-003 (image collection)
- MAN2CBZSPE-005 (CBZ writer)
- MAN2CBZSPE-007 (CLI flags)

## Acceptance Criteria

### Tests That Must Pass
Manual verification with test data:
```bash
# Basic processing
./manga2cbz ./source/isekai-craft-gurashi --out ./archive
-> Creates CBZ files in ./archive/

# Verify output
ls ./archive/*.cbz | wc -l
-> Should show 82 (or number of chapters)

# Validate archives
for f in ./archive/*.cbz; do unzip -t "$f" > /dev/null || echo "FAIL: $f"; done
-> No failures
```

### Processing Flow
```
1. Parse CLI arguments
2. Validate input directory exists
3. Discover chapters (flat or recursive based on -r)
4. For each chapter:
   a. Collect images with extension filter
   b. If no images, log warning and skip
   c. Create CBZ file
   d. Log success/failure
5. Report summary
6. Exit with appropriate code
```

### Exit Codes
- `0`: Success (all chapters processed)
- `1`: Partial failure (some chapters failed)
- `2`: Total failure (invalid args, no chapters found)

### Invariants That Must Remain True
- **Invariant #1**: Every image file included exactly once in CBZ
- **Invariant #3**: Proper naming (folder name -> CBZ name)
- Empty chapters logged as warnings, not errors
- Processing continues after individual chapter failures

## Definition of Done
- [x] Main processing loop implemented
- [x] Chapters discovered and processed
- [x] CBZ files created in output directory
- [x] Exit codes match spec
- [x] Partial failures don't stop processing
- [x] Empty chapters handled gracefully

## Outcome

**Status**: ✅ Completed

**Implementation**: `cmd/manga2cbz/main.go`

**Key Implementation Details**:
- Main processing loop in `process()` function (lines 196-266)
- Chapter discovery via `chapter.Discover()` (line 198)
- Image collection via `chapter.CollectImages()` per chapter (line 225)
- CBZ creation via `cbz.Create()` (line 246)
- Exit codes defined as constants (lines 22-26):
  - `exitSuccess = 0`: All chapters processed
  - `exitPartial = 1`: Some chapters failed
  - `exitError = 2`: Total failure (invalid args, no chapters)
- Partial failure handling: `continue` statements after errors (lines 229, 249)
- Empty chapter handling: Warning logged, processing continues (lines 232-237)
- Summary reporting at end of processing (lines 253-256)

**Test Coverage**:
- 48 tests pass with race detection
- `cmd/manga2cbz`: 82.5% coverage (15 tests)
- Orchestration-specific tests: `TestRun_EndToEnd`, `TestRun_NoChaptersFound`, `TestRun_QuietMode`

**Note**: Implementation was completed alongside MAN2CBZSPE-007 (CLI flags). All orchestration acceptance criteria verified through existing test suite.
