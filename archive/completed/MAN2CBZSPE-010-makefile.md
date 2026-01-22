# MAN2CBZSPE-010: Makefile and Build Configuration

## Description
Create Makefile with build, test, and release targets.

## Files to Touch
- `Makefile` (create)

## Out of Scope
- CI/CD configuration
- Docker configuration
- Installation scripts
- Package manager integration

## Acceptance Criteria

### Tests That Must Pass
```bash
# Build target
make build
-> Creates ./manga2cbz binary

# Test target
make test
-> Runs all tests with race detection and coverage

# Release target
make release
-> Creates optimized binary with stripped symbols

# Windows cross-compile
make windows
-> Creates ./manga2cbz.exe

# Clean target
make clean
-> Removes built binaries
```

### Makefile Contents
```makefile
.PHONY: build release test clean windows

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

# Clean build artifacts
clean:
	rm -f manga2cbz manga2cbz.exe
```

### Invariants That Must Remain True
- `make build` produces working binary
- `make test` runs with race detection
- Makefile targets match spec commands

## Definition of Done
- [x] Makefile created with all targets
- [x] `make build` works
- [x] `make test` works
- [x] `make release` produces smaller binary
- [x] `make windows` cross-compiles
- [x] `make clean` removes artifacts

## Outcome

**Status**: COMPLETED

**Implementation**:
- Created `/home/joeloverbeck/projects/manga2cbz/Makefile` with all required targets
- Makefile uses `.PHONY` declarations for all targets
- All targets use standard Go build commands

**Verification Results**:
| Target | Result |
|--------|--------|
| `make build` | ✅ Creates working binary (2.7MB) |
| `make test` | ✅ All 63 tests pass with race detection, 82-96% coverage |
| `make release` | ✅ Creates optimized binary (1.77MB, 34% smaller) |
| `make windows` | ✅ Creates valid PE32+ executable for Windows x86-64 |
| `make clean` | ✅ Removes both binaries |

**Notes**:
- The Makefile uses bare `go` commands (not `/usr/local/go/bin/go`) for portability across systems
- Release build achieves ~34% size reduction through `-trimpath` and `-ldflags "-s -w"`
- Cross-compilation to Windows verified via `file` command showing PE32+ format
