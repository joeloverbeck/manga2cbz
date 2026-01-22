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
