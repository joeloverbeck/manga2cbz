// Package chapter provides functionality for manga chapter processing.
package chapter

import (
	"os"
	"path/filepath"
	"strings"

	"manga2cbz/internal/sort"
)

// ImageFile represents an image file to be included in a CBZ archive.
type ImageFile struct {
	Path string // Full absolute path to file
	Name string // Base filename (for archive entry)
}

// CollectImages finds all image files in a directory matching the given extensions.
// Returns images sorted in natural order (so "10.jpg" comes after "9.jpg").
// Extensions are matched case-insensitively without leading dots.
func CollectImages(dir string, extensions []string) ([]ImageFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Build extension lookup set (lowercase, without dots)
	extSet := make(map[string]bool, len(extensions))
	for _, ext := range extensions {
		extSet[strings.ToLower(strings.TrimPrefix(ext, "."))] = true
	}

	// Collect matching filenames
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := strings.TrimPrefix(filepath.Ext(name), ".")
		if extSet[strings.ToLower(ext)] {
			names = append(names, name)
		}
	}

	// Natural sort the filenames
	sort.Natural(names)

	// Build result with absolute paths
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	images := make([]ImageFile, len(names))
	for i, name := range names {
		images[i] = ImageFile{
			Path: filepath.Join(absDir, name),
			Name: name,
		}
	}

	return images, nil
}
