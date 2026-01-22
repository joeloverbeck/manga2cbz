// Package cbz provides functionality for creating CBZ (Comic Book ZIP) archives.
package cbz

import (
	"archive/zip"
	"errors"
	"io"
	"os"

	"manga2cbz/internal/chapter"
)

// CreateOptions configures CBZ archive creation behavior.
type CreateOptions struct {
	Force bool // Overwrite existing files if true
}

// Create creates a CBZ archive at outputPath containing the given images.
// Images are stored at the archive root level using their Name field.
// Uses Store method (no compression) since images are already compressed.
// Streams files via io.Copy to avoid loading entire images into memory.
// Cleans up partial files on error.
func Create(outputPath string, images []chapter.ImageFile, opts CreateOptions) (err error) {
	// Check if file exists when Force is false
	if !opts.Force {
		if _, statErr := os.Stat(outputPath); statErr == nil {
			return errors.New("file already exists: " + outputPath)
		}
	}

	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	// Track whether we completed successfully for cleanup
	success := false
	defer func() {
		outFile.Close()
		if !success {
			os.Remove(outputPath)
		}
	}()

	// Create ZIP writer
	zipWriter := zip.NewWriter(outFile)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil && err == nil {
			err = closeErr
			success = false
		}
	}()

	// Add each image to the archive
	for _, img := range images {
		if addErr := addImageToArchive(zipWriter, img); addErr != nil {
			return addErr
		}
	}

	success = true
	return nil
}

// addImageToArchive adds a single image file to the ZIP archive.
// Uses Store method (no compression) and streams the file content.
func addImageToArchive(zw *zip.Writer, img chapter.ImageFile) error {
	// Open source file
	srcFile, err := os.Open(img.Path)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get file info for size
	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create ZIP entry header with Store method
	header := &zip.FileHeader{
		Name:   img.Name, // Store at root level
		Method: zip.Store,
	}
	header.SetModTime(info.ModTime())

	// Create the entry in the archive
	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}

	// Stream the file content
	_, err = io.Copy(writer, srcFile)
	return err
}

// Validate checks if a CBZ file is a valid ZIP archive.
// Returns nil if valid, or an error describing the problem.
func Validate(cbzPath string) error {
	reader, err := zip.OpenReader(cbzPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Attempt to iterate through entries to verify archive integrity
	for _, f := range reader.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		rc.Close()
	}

	return nil
}
