package cbz

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"manga2cbz/internal/chapter"
)

// createTestImages creates temporary image files for testing.
// Returns the temp directory and slice of ImageFile structs.
func createTestImages(t *testing.T, count int) (string, []chapter.ImageFile) {
	t.Helper()
	tmpDir := t.TempDir()

	images := make([]chapter.ImageFile, count)
	for i := 0; i < count; i++ {
		name := filepath.Join(tmpDir, "image"+string(rune('0'+i))+".jpg")
		content := []byte("fake image content " + string(rune('0'+i)))
		if err := os.WriteFile(name, content, 0644); err != nil {
			t.Fatalf("failed to create test image: %v", err)
		}
		images[i] = chapter.ImageFile{
			Path: name,
			Name: filepath.Base(name),
		}
	}
	return tmpDir, images
}

func TestCreate_Basic(t *testing.T) {
	tmpDir, images := createTestImages(t, 3)
	outputPath := filepath.Join(tmpDir, "output.cbz")

	err := Create(outputPath, images, CreateOptions{})
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("output file does not exist")
	}

	// Verify archive has exactly 3 entries
	reader, err := zip.OpenReader(outputPath)
	if err != nil {
		t.Fatalf("failed to open archive: %v", err)
	}
	defer reader.Close()

	if len(reader.File) != 3 {
		t.Errorf("expected 3 entries, got %d", len(reader.File))
	}
}

func TestCreate_FlatStructure(t *testing.T) {
	tmpDir := t.TempDir()

	// Create images with nested source paths
	subDir := filepath.Join(tmpDir, "nested", "deep")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	images := []chapter.ImageFile{
		{Path: filepath.Join(subDir, "image1.jpg"), Name: "image1.jpg"},
		{Path: filepath.Join(subDir, "image2.jpg"), Name: "image2.jpg"},
	}

	// Create the source files
	for _, img := range images {
		if err := os.WriteFile(img.Path, []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	outputPath := filepath.Join(tmpDir, "output.cbz")
	err := Create(outputPath, images, CreateOptions{})
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify all entries are at root level (no directory separators)
	reader, err := zip.OpenReader(outputPath)
	if err != nil {
		t.Fatalf("failed to open archive: %v", err)
	}
	defer reader.Close()

	for _, f := range reader.File {
		if filepath.Dir(f.Name) != "." {
			t.Errorf("entry %q is not at root level", f.Name)
		}
	}
}

func TestCreate_StoreMethod(t *testing.T) {
	tmpDir, images := createTestImages(t, 2)
	outputPath := filepath.Join(tmpDir, "output.cbz")

	err := Create(outputPath, images, CreateOptions{})
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify ZIP entries use Store method (compression method 0)
	reader, err := zip.OpenReader(outputPath)
	if err != nil {
		t.Fatalf("failed to open archive: %v", err)
	}
	defer reader.Close()

	for _, f := range reader.File {
		if f.Method != zip.Store {
			t.Errorf("entry %q uses compression method %d, expected Store (0)", f.Name, f.Method)
		}
	}
}

func TestCreate_CleanupOnError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create one valid image and one that will fail
	validPath := filepath.Join(tmpDir, "valid.jpg")
	if err := os.WriteFile(validPath, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	images := []chapter.ImageFile{
		{Path: validPath, Name: "valid.jpg"},
		{Path: filepath.Join(tmpDir, "nonexistent.jpg"), Name: "nonexistent.jpg"},
	}

	outputPath := filepath.Join(tmpDir, "output.cbz")
	err := Create(outputPath, images, CreateOptions{})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}

	// Verify partial file was cleaned up
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Error("partial output file was not cleaned up")
	}
}

func TestCreate_NoOverwrite(t *testing.T) {
	tmpDir, images := createTestImages(t, 1)
	outputPath := filepath.Join(tmpDir, "existing.cbz")

	// Create existing file
	originalContent := []byte("original content")
	if err := os.WriteFile(outputPath, originalContent, 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Attempt to create without Force
	err := Create(outputPath, images, CreateOptions{Force: false})
	if err == nil {
		t.Fatal("expected error when file exists and Force=false")
	}

	// Verify original file is unchanged
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != string(originalContent) {
		t.Error("existing file was modified")
	}
}

func TestCreate_ForceOverwrite(t *testing.T) {
	tmpDir, images := createTestImages(t, 1)
	outputPath := filepath.Join(tmpDir, "existing.cbz")

	// Create existing file
	if err := os.WriteFile(outputPath, []byte("original"), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	// Create with Force=true
	err := Create(outputPath, images, CreateOptions{Force: true})
	if err != nil {
		t.Fatalf("Create() with Force=true failed: %v", err)
	}

	// Verify file was overwritten (should be valid ZIP now)
	if err := Validate(outputPath); err != nil {
		t.Errorf("overwritten file is not a valid archive: %v", err)
	}
}

func TestValidate_ValidArchive(t *testing.T) {
	tmpDir, images := createTestImages(t, 2)
	outputPath := filepath.Join(tmpDir, "valid.cbz")

	err := Create(outputPath, images, CreateOptions{})
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	err = Validate(outputPath)
	if err != nil {
		t.Errorf("Validate() returned error for valid archive: %v", err)
	}
}

func TestValidate_CorruptArchive(t *testing.T) {
	tmpDir := t.TempDir()
	corruptPath := filepath.Join(tmpDir, "corrupt.cbz")

	// Create file with invalid ZIP content
	if err := os.WriteFile(corruptPath, []byte("not a zip file"), 0644); err != nil {
		t.Fatalf("failed to create corrupt file: %v", err)
	}

	err := Validate(corruptPath)
	if err == nil {
		t.Error("Validate() should return error for corrupt archive")
	}
}

func TestCreate_EmptyImageList(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "empty.cbz")

	err := Create(outputPath, []chapter.ImageFile{}, CreateOptions{})
	if err != nil {
		t.Fatalf("Create() failed with empty list: %v", err)
	}

	// Verify empty archive is valid
	reader, err := zip.OpenReader(outputPath)
	if err != nil {
		t.Fatalf("failed to open archive: %v", err)
	}
	defer reader.Close()

	if len(reader.File) != 0 {
		t.Errorf("expected 0 entries, got %d", len(reader.File))
	}
}

func TestCreate_InvalidOutputPath(t *testing.T) {
	_, images := createTestImages(t, 1)

	// Try to write to a path that doesn't exist
	outputPath := "/nonexistent/directory/output.cbz"
	err := Create(outputPath, images, CreateOptions{})
	if err == nil {
		t.Error("expected error for invalid output path")
	}
}

func TestValidate_NonexistentFile(t *testing.T) {
	err := Validate("/nonexistent/file.cbz")
	if err == nil {
		t.Error("Validate() should return error for nonexistent file")
	}
}
