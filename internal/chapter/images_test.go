package chapter

import (
	"os"
	"path/filepath"
	"testing"
)

// createTestFiles creates empty files in the given directory.
func createTestFiles(t *testing.T, dir string, names []string) {
	t.Helper()
	for _, name := range names {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", name, err)
		}
	}
}

func TestCollectImages_BasicJPG(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"01.jpg", "02.jpg", "03.jpg"})

	images, err := CollectImages(dir, []string{"jpg", "jpeg", "png"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(images) != 3 {
		t.Fatalf("expected 3 images, got %d", len(images))
	}

	// Verify paths are absolute
	for _, img := range images {
		if !filepath.IsAbs(img.Path) {
			t.Errorf("expected absolute path, got %s", img.Path)
		}
	}

	// Verify names are basenames
	expectedNames := []string{"01.jpg", "02.jpg", "03.jpg"}
	for i, img := range images {
		if img.Name != expectedNames[i] {
			t.Errorf("expected name %s, got %s", expectedNames[i], img.Name)
		}
	}
}

func TestCollectImages_MixedTypes(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"01.jpg", "02.png", "03.gif"})

	images, err := CollectImages(dir, []string{"jpg", "png", "gif"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(images) != 3 {
		t.Fatalf("expected 3 images, got %d", len(images))
	}
}

func TestCollectImages_FilterNonImages(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"01.jpg", "notes.txt", "readme.md"})

	images, err := CollectImages(dir, []string{"jpg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}

	if images[0].Name != "01.jpg" {
		t.Errorf("expected 01.jpg, got %s", images[0].Name)
	}
}

func TestCollectImages_CaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"01.JPG", "02.Png", "03.jpeg"})

	images, err := CollectImages(dir, []string{"jpg", "png", "jpeg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(images) != 3 {
		t.Fatalf("expected 3 images, got %d", len(images))
	}
}

func TestCollectImages_NaturalOrder(t *testing.T) {
	dir := t.TempDir()
	// Create files in non-natural order
	createTestFiles(t, dir, []string{"10.jpg", "2.jpg", "1.jpg"})

	images, err := CollectImages(dir, []string{"jpg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(images) != 3 {
		t.Fatalf("expected 3 images, got %d", len(images))
	}

	// Verify natural sort order: 1, 2, 10
	expectedOrder := []string{"1.jpg", "2.jpg", "10.jpg"}
	for i, img := range images {
		if img.Name != expectedOrder[i] {
			t.Errorf("position %d: expected %s, got %s", i, expectedOrder[i], img.Name)
		}
	}
}

func TestCollectImages_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	images, err := CollectImages(dir, []string{"jpg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(images) != 0 {
		t.Fatalf("expected empty slice, got %d images", len(images))
	}
}

func TestCollectImages_NotFound(t *testing.T) {
	_, err := CollectImages("/nonexistent/path/that/does/not/exist", []string{"jpg"})
	if err == nil {
		t.Fatal("expected error for nonexistent directory, got nil")
	}
}

// Additional edge case tests for robustness

func TestCollectImages_ExtensionsWithDots(t *testing.T) {
	// Extensions passed with leading dots should still work
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"image.jpg", "image.png"})

	images, err := CollectImages(dir, []string{".jpg", ".png"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(images))
	}
}

func TestCollectImages_SkipsSubdirectories(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"image.jpg"})

	// Create a subdirectory (should be ignored)
	subdir := filepath.Join(dir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	createTestFiles(t, subdir, []string{"nested.jpg"})

	images, err := CollectImages(dir, []string{"jpg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only the top-level image should be collected
	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}
	if images[0].Name != "image.jpg" {
		t.Errorf("expected image.jpg, got %s", images[0].Name)
	}
}

func TestCollectImages_NoMatchingExtensions(t *testing.T) {
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"image.jpg", "image.png"})

	// Search for extensions that don't exist
	images, err := CollectImages(dir, []string{"gif", "bmp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(images) != 0 {
		t.Fatalf("expected 0 images, got %d", len(images))
	}
}

func TestCollectImages_NaturalOrderWithPrefix(t *testing.T) {
	// Test natural sort with common manga naming patterns
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"page10.jpg", "page2.jpg", "page1.jpg", "page20.jpg"})

	images, err := CollectImages(dir, []string{"jpg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedOrder := []string{"page1.jpg", "page2.jpg", "page10.jpg", "page20.jpg"}
	for i, img := range images {
		if img.Name != expectedOrder[i] {
			t.Errorf("position %d: expected %s, got %s", i, expectedOrder[i], img.Name)
		}
	}
}

func TestCollectImages_PathFileConsistency(t *testing.T) {
	// Verify that Path ends with Name
	dir := t.TempDir()
	createTestFiles(t, dir, []string{"test.jpg"})

	images, err := CollectImages(dir, []string{"jpg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(images))
	}

	img := images[0]
	if filepath.Base(img.Path) != img.Name {
		t.Errorf("Path base (%s) doesn't match Name (%s)", filepath.Base(img.Path), img.Name)
	}
}
