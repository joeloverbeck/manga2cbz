package chapter

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper to create a directory structure for tests
func createDir(t *testing.T, base string, path string) string {
	t.Helper()
	fullPath := filepath.Join(base, path)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", fullPath, err)
	}
	return fullPath
}

// Helper to create a file for tests
func createFile(t *testing.T, base string, path string) {
	t.Helper()
	fullPath := filepath.Join(base, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}
	if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create file %s: %v", fullPath, err)
	}
}

func TestDiscover_Flat(t *testing.T) {
	root := t.TempDir()

	// Create chapter directories with images
	createDir(t, root, "ChapterA")
	createFile(t, root, "ChapterA/page1.jpg")
	createDir(t, root, "ChapterB")
	createFile(t, root, "ChapterB/page1.jpg")

	chapters, err := Discover(root, false)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 2 {
		t.Fatalf("Expected 2 chapters, got %d", len(chapters))
	}

	if chapters[0].Name != "ChapterA" {
		t.Errorf("Expected first chapter name 'ChapterA', got %q", chapters[0].Name)
	}
	if chapters[1].Name != "ChapterB" {
		t.Errorf("Expected second chapter name 'ChapterB', got %q", chapters[1].Name)
	}

	// Verify paths are absolute
	if !filepath.IsAbs(chapters[0].Path) {
		t.Errorf("Expected absolute path, got %q", chapters[0].Path)
	}
}

func TestDiscover_Recursive(t *testing.T) {
	root := t.TempDir()

	// Create nested structure:
	// root/
	//   Volume1/
	//     Chapter1/
	//     Chapter2/
	//   Volume2/
	//     Chapter3/
	createDir(t, root, "Volume1/Chapter1")
	createFile(t, root, "Volume1/Chapter1/page1.jpg")
	createDir(t, root, "Volume1/Chapter2")
	createFile(t, root, "Volume1/Chapter2/page1.jpg")
	createDir(t, root, "Volume2/Chapter3")
	createFile(t, root, "Volume2/Chapter3/page1.jpg")

	chapters, err := Discover(root, true)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 3 {
		t.Fatalf("Expected 3 chapters, got %d", len(chapters))
	}

	// Verify chapter names are relative paths
	expectedNames := []string{
		filepath.Join("Volume1", "Chapter1"),
		filepath.Join("Volume1", "Chapter2"),
		filepath.Join("Volume2", "Chapter3"),
	}
	for i, expected := range expectedNames {
		if chapters[i].Name != expected {
			t.Errorf("Chapter %d: expected name %q, got %q", i, expected, chapters[i].Name)
		}
	}
}

func TestDiscover_SkipFiles(t *testing.T) {
	root := t.TempDir()

	// Create one directory and one file
	createDir(t, root, "Chapter1")
	createFile(t, root, "Chapter1/page1.jpg")
	createFile(t, root, "notes.txt")

	chapters, err := Discover(root, false)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 1 {
		t.Fatalf("Expected 1 chapter, got %d", len(chapters))
	}

	if chapters[0].Name != "Chapter1" {
		t.Errorf("Expected chapter name 'Chapter1', got %q", chapters[0].Name)
	}
}

func TestDiscover_EmptyDir(t *testing.T) {
	root := t.TempDir()

	chapters, err := Discover(root, false)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 0 {
		t.Errorf("Expected 0 chapters for empty directory, got %d", len(chapters))
	}
}

func TestDiscover_NotFound(t *testing.T) {
	_, err := Discover("/nonexistent/path/that/does/not/exist", false)
	if err == nil {
		t.Error("Expected error for nonexistent path, got nil")
	}
}

func TestDiscover_NaturalOrder(t *testing.T) {
	root := t.TempDir()

	// Create chapters in non-natural order
	createDir(t, root, "Chapter 10")
	createFile(t, root, "Chapter 10/page1.jpg")
	createDir(t, root, "Chapter 2")
	createFile(t, root, "Chapter 2/page1.jpg")
	createDir(t, root, "Chapter 1")
	createFile(t, root, "Chapter 1/page1.jpg")

	chapters, err := Discover(root, false)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 3 {
		t.Fatalf("Expected 3 chapters, got %d", len(chapters))
	}

	// Verify natural sort order: 1, 2, 10 (not 1, 10, 2)
	expectedOrder := []string{"Chapter 1", "Chapter 2", "Chapter 10"}
	for i, expected := range expectedOrder {
		if chapters[i].Name != expected {
			t.Errorf("Chapter %d: expected %q, got %q", i, expected, chapters[i].Name)
		}
	}
}

func TestDiscover_SkipHidden(t *testing.T) {
	root := t.TempDir()

	// Create one visible and one hidden directory
	createDir(t, root, "Chapter1")
	createFile(t, root, "Chapter1/page1.jpg")
	createDir(t, root, ".hidden")
	createFile(t, root, ".hidden/page1.jpg")

	chapters, err := Discover(root, false)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 1 {
		t.Fatalf("Expected 1 chapter (hidden skipped), got %d", len(chapters))
	}

	if chapters[0].Name != "Chapter1" {
		t.Errorf("Expected chapter name 'Chapter1', got %q", chapters[0].Name)
	}
}

func TestDiscover_SkipHiddenRecursive(t *testing.T) {
	root := t.TempDir()

	// Create visible and hidden nested directories
	createDir(t, root, "Volume1/Chapter1")
	createFile(t, root, "Volume1/Chapter1/page1.jpg")
	createDir(t, root, ".hiddenVol/Chapter2")
	createFile(t, root, ".hiddenVol/Chapter2/page1.jpg")

	chapters, err := Discover(root, true)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 1 {
		t.Fatalf("Expected 1 chapter (hidden tree skipped), got %d", len(chapters))
	}

	expected := filepath.Join("Volume1", "Chapter1")
	if chapters[0].Name != expected {
		t.Errorf("Expected chapter name %q, got %q", expected, chapters[0].Name)
	}
}

func TestDiscover_SpacesInNames(t *testing.T) {
	root := t.TempDir()

	// Test Invariant #6: Handle spaces in directory names
	createDir(t, root, "Chapter 5")
	createFile(t, root, "Chapter 5/page 1.jpg")
	createDir(t, root, "Chapter 15 - Special")
	createFile(t, root, "Chapter 15 - Special/page 1.jpg")

	chapters, err := Discover(root, false)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 2 {
		t.Fatalf("Expected 2 chapters, got %d", len(chapters))
	}

	// Verify natural sort: "Chapter 5" before "Chapter 15 - Special"
	if chapters[0].Name != "Chapter 5" {
		t.Errorf("Expected first chapter 'Chapter 5', got %q", chapters[0].Name)
	}
	if chapters[1].Name != "Chapter 15 - Special" {
		t.Errorf("Expected second chapter 'Chapter 15 - Special', got %q", chapters[1].Name)
	}
}

func TestDiscover_AbsolutePaths(t *testing.T) {
	root := t.TempDir()

	createDir(t, root, "Chapter1")
	createFile(t, root, "Chapter1/page1.jpg")

	chapters, err := Discover(root, false)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 1 {
		t.Fatalf("Expected 1 chapter, got %d", len(chapters))
	}

	// Verify path is absolute
	if !filepath.IsAbs(chapters[0].Path) {
		t.Errorf("Expected absolute path, got %q", chapters[0].Path)
	}

	// Verify path exists
	if _, err := os.Stat(chapters[0].Path); err != nil {
		t.Errorf("Chapter path does not exist: %v", err)
	}
}

func TestDiscover_RecursiveNaturalOrder(t *testing.T) {
	root := t.TempDir()

	// Create nested chapters in non-natural order
	createDir(t, root, "Vol 2/Chapter 10")
	createFile(t, root, "Vol 2/Chapter 10/page1.jpg")
	createDir(t, root, "Vol 1/Chapter 2")
	createFile(t, root, "Vol 1/Chapter 2/page1.jpg")
	createDir(t, root, "Vol 1/Chapter 1")
	createFile(t, root, "Vol 1/Chapter 1/page1.jpg")

	chapters, err := Discover(root, true)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(chapters) != 3 {
		t.Fatalf("Expected 3 chapters, got %d", len(chapters))
	}

	// Verify natural sort order
	expectedOrder := []string{
		filepath.Join("Vol 1", "Chapter 1"),
		filepath.Join("Vol 1", "Chapter 2"),
		filepath.Join("Vol 2", "Chapter 10"),
	}
	for i, expected := range expectedOrder {
		if chapters[i].Name != expected {
			t.Errorf("Chapter %d: expected %q, got %q", i, expected, chapters[i].Name)
		}
	}
}
