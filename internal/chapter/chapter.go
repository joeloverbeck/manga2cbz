// Package chapter provides functionality for manga chapter processing.
package chapter

import (
	"os"
	"path/filepath"
	"strings"

	"manga2cbz/internal/sort"
)

// Chapter represents a manga chapter directory.
type Chapter struct {
	Name string // Directory name (becomes CBZ filename)
	Path string // Full absolute path to directory
}

// Discover finds chapter directories in the input directory.
// If recursive is false, only immediate subdirectories are considered.
// If recursive is true, all nested directories containing images are found.
// Results are sorted in natural order (Chapter 2 before Chapter 10).
// Hidden directories (starting with .) are skipped.
func Discover(inputDir string, recursive bool) ([]Chapter, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(inputDir)
	if err != nil {
		return nil, err
	}

	// Verify directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &os.PathError{Op: "discover", Path: absPath, Err: os.ErrInvalid}
	}

	if recursive {
		return discoverRecursive(absPath)
	}
	return discoverFlat(absPath)
}

// discoverFlat finds chapter directories one level deep.
func discoverFlat(inputDir string) ([]Chapter, error) {
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}

	var chapters []Chapter
	for _, entry := range entries {
		// Skip non-directories
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		chapters = append(chapters, Chapter{
			Name: entry.Name(),
			Path: filepath.Join(inputDir, entry.Name()),
		})
	}

	// Sort chapters naturally
	sortChapters(chapters)

	return chapters, nil
}

// discoverRecursive finds chapter directories at all depths.
// A directory is considered a chapter if it contains no subdirectories
// (leaf node in the directory tree).
func discoverRecursive(inputDir string) ([]Chapter, error) {
	var chapters []Chapter

	err := filepath.WalkDir(inputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip non-directories
		if !d.IsDir() {
			return nil
		}

		// Skip hidden directories
		if strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		// Skip root directory
		if path == inputDir {
			return nil
		}

		// Check if this directory is a leaf (no subdirectories)
		isLeaf, err := isLeafDirectory(path)
		if err != nil {
			return err
		}

		if isLeaf {
			// Build name relative to input directory
			relPath, err := filepath.Rel(inputDir, path)
			if err != nil {
				return err
			}

			chapters = append(chapters, Chapter{
				Name: relPath,
				Path: path,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort chapters naturally
	sortChapters(chapters)

	return chapters, nil
}

// isLeafDirectory returns true if the directory contains no subdirectories.
func isLeafDirectory(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			return false, nil
		}
	}

	return true, nil
}

// sortChapters sorts chapters by name in natural order.
func sortChapters(chapters []Chapter) {
	// Extract names for sorting
	names := make([]string, len(chapters))
	nameToChapter := make(map[string]Chapter)
	for i, ch := range chapters {
		names[i] = ch.Name
		nameToChapter[ch.Name] = ch
	}

	// Sort names naturally
	sort.Natural(names)

	// Rebuild chapters slice in sorted order
	for i, name := range names {
		chapters[i] = nameToChapter[name]
	}
}
