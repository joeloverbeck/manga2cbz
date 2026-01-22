// Package convert provides image format conversion functionality.
package convert

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"

	"manga2cbz/internal/chapter"
)

// ConvertWebPImages converts any WebP images in the slice to PNG format.
// Converted files are written to a temporary directory.
// Returns an updated ImageFile slice with converted paths and a cleanup function.
// The cleanup function should be called (typically via defer) to remove temp files.
// Non-WebP files are passed through unchanged.
func ConvertWebPImages(images []chapter.ImageFile) ([]chapter.ImageFile, func(), error) {
	// Quick check: any WebP files?
	hasWebP := false
	for _, img := range images {
		if isWebP(img.Name) {
			hasWebP = true
			break
		}
	}

	// No WebP files, return as-is with no-op cleanup
	if !hasWebP {
		return images, func() {}, nil
	}

	// Create temp directory for converted files
	tempDir, err := os.MkdirTemp("", "manga2cbz-convert-*")
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	// Process images, converting WebP to PNG
	result := make([]chapter.ImageFile, len(images))
	for i, img := range images {
		if !isWebP(img.Name) {
			// Pass through non-WebP files unchanged
			result[i] = img
			continue
		}

		// Convert WebP to PNG
		converted, err := convertWebPToPNG(img, tempDir)
		if err != nil {
			cleanup()
			return nil, nil, err
		}
		result[i] = converted
	}

	return result, cleanup, nil
}

// isWebP checks if a filename has a .webp extension (case-insensitive).
func isWebP(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".webp"
}

// convertWebPToPNG converts a single WebP image to PNG format.
// The converted file is written to the temp directory.
// Returns an ImageFile with updated path and name.
func convertWebPToPNG(img chapter.ImageFile, tempDir string) (chapter.ImageFile, error) {
	// Open source WebP file
	srcFile, err := os.Open(img.Path)
	if err != nil {
		return chapter.ImageFile{}, err
	}
	defer srcFile.Close()

	// Decode WebP image
	decodedImg, err := webp.Decode(srcFile)
	if err != nil {
		return chapter.ImageFile{}, err
	}

	// Generate new filename with .png extension
	baseName := strings.TrimSuffix(img.Name, filepath.Ext(img.Name))
	newName := baseName + ".png"
	newPath := filepath.Join(tempDir, newName)

	// Create destination file
	dstFile, err := os.Create(newPath)
	if err != nil {
		return chapter.ImageFile{}, err
	}
	defer dstFile.Close()

	// Encode as PNG
	if err := png.Encode(dstFile, decodedImg); err != nil {
		return chapter.ImageFile{}, err
	}

	return chapter.ImageFile{
		Path: newPath,
		Name: newName,
	}, nil
}
