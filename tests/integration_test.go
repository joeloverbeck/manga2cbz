// Package tests contains integration tests for manga2cbz.
package tests

import (
	"archive/zip"
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runCLI executes the manga2cbz CLI with the given arguments.
// Returns stdout, stderr, and exit code.
func runCLI(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	// Get the module root (one level up from tests/)
	moduleRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("failed to get module root: %v", err)
	}

	cmd := exec.Command("/usr/local/go/bin/go", append([]string{"run", "./cmd/manga2cbz"}, args...)...)
	cmd.Dir = moduleRoot

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	exitCode = 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("failed to run CLI: %v", err)
	}

	return stdoutBuf.String(), stderrBuf.String(), exitCode
}

// createTestImage creates a minimal PNG image file at the given path.
func createTestImage(t *testing.T, path string) {
	t.Helper()

	// Create a 1x1 pixel PNG
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test image %s: %v", path, err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		t.Fatalf("failed to encode PNG %s: %v", path, err)
	}
}

// createChapterWithImages creates a chapter directory with numbered image files.
func createChapterWithImages(t *testing.T, baseDir, chapterName string, imageNames []string) string {
	t.Helper()

	chapterPath := filepath.Join(baseDir, chapterName)
	if err := os.MkdirAll(chapterPath, 0755); err != nil {
		t.Fatalf("failed to create chapter dir %s: %v", chapterPath, err)
	}

	for _, name := range imageNames {
		createTestImage(t, filepath.Join(chapterPath, name))
	}

	return chapterPath
}

// getCBZFileList opens a CBZ and returns the list of files in archive order.
func getCBZFileList(t *testing.T, cbzPath string) []string {
	t.Helper()

	reader, err := zip.OpenReader(cbzPath)
	if err != nil {
		t.Fatalf("failed to open CBZ %s: %v", cbzPath, err)
	}
	defer reader.Close()

	var files []string
	for _, f := range reader.File {
		files = append(files, f.Name)
	}
	return files
}

// validateCBZ verifies that a CBZ is a valid ZIP archive.
func validateCBZ(t *testing.T, cbzPath string) {
	t.Helper()

	reader, err := zip.OpenReader(cbzPath)
	if err != nil {
		t.Fatalf("CBZ validation failed for %s: %v", cbzPath, err)
	}
	defer reader.Close()

	// Verify each entry can be read
	for _, f := range reader.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("CBZ entry %s is corrupt in %s: %v", f.Name, cbzPath, err)
		}
		rc.Close()
	}
}

// ============================================================================
// SPEC TEST CASES (5)
// ============================================================================

// TestIntegration_BasicMultipleChapters verifies that multiple chapters
// are correctly converted to CBZ archives.
// Invariants verified: #1 (complete), #3 (naming), #4 (valid)
func TestIntegration_BasicMultipleChapters(t *testing.T) {
	// Setup: Create temp directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "manga")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Create ChapterA with 3 images
	createChapterWithImages(t, inputDir, "ChapterA", []string{
		"01.png", "02.png", "03.png",
	})

	// Create ChapterB with 2 images
	createChapterWithImages(t, inputDir, "ChapterB", []string{
		"01.png", "02.png",
	})

	// Run CLI (flags must come before positional args in Go's flag package)
	stdout, stderr, exitCode := runCLI(t, "--out", outputDir, inputDir)

	// Verify exit code 0
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify ChapterA.cbz exists with 3 images
	cbzA := filepath.Join(outputDir, "ChapterA.cbz")
	if _, err := os.Stat(cbzA); os.IsNotExist(err) {
		t.Fatalf("ChapterA.cbz not created")
	}
	filesA := getCBZFileList(t, cbzA)
	if len(filesA) != 3 {
		t.Errorf("ChapterA.cbz: expected 3 images, got %d: %v", len(filesA), filesA)
	}
	validateCBZ(t, cbzA) // Invariant #4

	// Verify ChapterB.cbz exists with 2 images
	cbzB := filepath.Join(outputDir, "ChapterB.cbz")
	if _, err := os.Stat(cbzB); os.IsNotExist(err) {
		t.Fatalf("ChapterB.cbz not created")
	}
	filesB := getCBZFileList(t, cbzB)
	if len(filesB) != 2 {
		t.Errorf("ChapterB.cbz: expected 2 images, got %d: %v", len(filesB), filesB)
	}
	validateCBZ(t, cbzB) // Invariant #4
}

// TestIntegration_NaturalSortOrdering verifies that images are sorted
// in natural order within CBZ archives (2 before 10).
// Invariants verified: #2 (ordering)
func TestIntegration_NaturalSortOrdering(t *testing.T) {
	// Setup: Create temp directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "manga")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Create chapter with out-of-order numbered images
	// Lexicographic order would be: 1, 10, 11, 2, 3
	// Natural order should be: 1, 2, 3, 10, 11
	createChapterWithImages(t, inputDir, "Chapter1", []string{
		"1.png", "2.png", "10.png", "11.png", "3.png",
	})

	// Run CLI (flags must come before positional args in Go's flag package)
	stdout, stderr, exitCode := runCLI(t, "--out", outputDir, inputDir)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify file order in CBZ
	cbzPath := filepath.Join(outputDir, "Chapter1.cbz")
	files := getCBZFileList(t, cbzPath)

	expected := []string{"1.png", "2.png", "3.png", "10.png", "11.png"}
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d: %v", len(expected), len(files), files)
	}

	for i, exp := range expected {
		if files[i] != exp {
			t.Errorf("position %d: expected %s, got %s\nfull order: %v", i, exp, files[i], files)
		}
	}
}

// TestIntegration_MixedFileTypes verifies that only image files are included
// in CBZ archives and non-image files are excluded.
// Invariants verified: #1 (complete), #6 (robustness)
func TestIntegration_MixedFileTypes(t *testing.T) {
	// Setup: Create temp directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "manga")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Create chapter dir
	chapterPath := filepath.Join(inputDir, "Chapter1")
	if err := os.MkdirAll(chapterPath, 0755); err != nil {
		t.Fatalf("failed to create chapter dir: %v", err)
	}

	// Create image files
	createTestImage(t, filepath.Join(chapterPath, "01.jpg"))
	createTestImage(t, filepath.Join(chapterPath, "02.png"))
	createTestImage(t, filepath.Join(chapterPath, "03.gif"))

	// Create non-image file
	if err := os.WriteFile(filepath.Join(chapterPath, "notes.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create notes.txt: %v", err)
	}

	// Run CLI (flags must come before positional args in Go's flag package)
	stdout, stderr, exitCode := runCLI(t, "--out", outputDir, inputDir)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify CBZ contains only image files
	cbzPath := filepath.Join(outputDir, "Chapter1.cbz")
	files := getCBZFileList(t, cbzPath)

	if len(files) != 3 {
		t.Errorf("expected 3 image files, got %d: %v", len(files), files)
	}

	for _, f := range files {
		if f == "notes.txt" {
			t.Errorf("non-image file 'notes.txt' should not be in CBZ")
		}
	}
}

// TestIntegration_EmptyChapter verifies that empty chapters (no images)
// are handled gracefully with a warning and no crash.
// Invariants verified: #6 (robustness)
func TestIntegration_EmptyChapter(t *testing.T) {
	// Setup: Create temp directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "manga")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Create empty chapter directory (no images)
	emptyChapter := filepath.Join(inputDir, "EmptyChapter")
	if err := os.MkdirAll(emptyChapter, 0755); err != nil {
		t.Fatalf("failed to create empty chapter dir: %v", err)
	}

	// Run CLI (flags must come before positional args in Go's flag package)
	stdout, stderr, exitCode := runCLI(t, "--out", outputDir, inputDir)

	// Should complete without crashing, exit 0 (empty chapters are skipped, not failed)
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify warning was logged
	if !strings.Contains(stderr, "Warning") && !strings.Contains(stderr, "no images") {
		t.Errorf("expected warning about empty chapter in stderr: %s", stderr)
	}

	// Verify no CBZ was created for empty chapter
	cbzPath := filepath.Join(outputDir, "EmptyChapter.cbz")
	if _, err := os.Stat(cbzPath); !os.IsNotExist(err) {
		t.Errorf("EmptyChapter.cbz should not be created for empty chapter")
	}
}

// TestIntegration_OverwriteBehavior verifies the --force flag behavior.
// Without --force, existing files should be skipped.
// With --force, existing files should be overwritten.
// Invariants verified: #6 (robustness)
func TestIntegration_OverwriteBehavior(t *testing.T) {
	// Setup: Create temp directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "manga")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Create chapter with images
	createChapterWithImages(t, inputDir, "Chapter1", []string{
		"01.png", "02.png",
	})

	// First run - should create CBZ (flags must come before positional args)
	stdout, stderr, exitCode := runCLI(t, "--out", outputDir, inputDir)
	if exitCode != 0 {
		t.Fatalf("first run: expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	cbzPath := filepath.Join(outputDir, "Chapter1.cbz")

	// Get original file info
	origInfo, err := os.Stat(cbzPath)
	if err != nil {
		t.Fatalf("CBZ not created on first run: %v", err)
	}

	// Second run without --force - should fail with error
	stdout, stderr, exitCode = runCLI(t, "--out", outputDir, inputDir)

	// Should have non-zero exit due to file existing
	// With 1 chapter failing out of 1, failures == len(chapters), so exit 2
	if exitCode == 0 {
		t.Errorf("second run without --force: expected non-zero exit, got 0")
	}

	// Third run with --force - should overwrite successfully
	stdout, stderr, exitCode = runCLI(t, "--force", "--out", outputDir, inputDir)
	if exitCode != 0 {
		t.Errorf("run with --force: expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify file was overwritten (modtime should be different or equal)
	newInfo, err := os.Stat(cbzPath)
	if err != nil {
		t.Fatalf("CBZ should exist after --force: %v", err)
	}

	// File should still be valid
	validateCBZ(t, cbzPath)

	// Size should be the same (same content)
	if newInfo.Size() != origInfo.Size() {
		t.Logf("Note: file size changed from %d to %d (expected same)", origInfo.Size(), newInfo.Size())
	}
}

// ============================================================================
// BONUS TEST CASES (2)
// ============================================================================

// TestIntegration_RecursiveMode verifies the -r flag processes nested directories.
// BONUS: Not in original spec test cases but useful for coverage.
// Invariants verified: #1 (complete), #3 (naming), #5 (no internal folders)
func TestIntegration_RecursiveMode(t *testing.T) {
	// Setup: Create nested directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "manga")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Create nested structure:
	// manga/Volume1/Chapter1/
	// manga/Volume1/Chapter2/
	// manga/Volume2/Chapter3/
	createChapterWithImages(t, inputDir, filepath.Join("Volume1", "Chapter1"), []string{"01.png"})
	createChapterWithImages(t, inputDir, filepath.Join("Volume1", "Chapter2"), []string{"01.png"})
	createChapterWithImages(t, inputDir, filepath.Join("Volume2", "Chapter3"), []string{"01.png"})

	// Run CLI with recursive flag (flags must come before positional args)
	stdout, stderr, exitCode := runCLI(t, "-r", "--out", outputDir, inputDir)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify 3 CBZ files created
	// Note: In recursive mode, names include path separators replaced with underscores
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("failed to read output dir: %v", err)
	}

	cbzCount := 0
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".cbz") {
			cbzCount++

			// Verify each CBZ is valid and has images at root (Invariant #5)
			cbzPath := filepath.Join(outputDir, e.Name())
			files := getCBZFileList(t, cbzPath)
			for _, f := range files {
				if strings.Contains(f, "/") || strings.Contains(f, "\\") {
					t.Errorf("CBZ %s has nested path: %s (violates invariant #5)", e.Name(), f)
				}
			}
			validateCBZ(t, cbzPath)
		}
	}

	if cbzCount != 3 {
		t.Errorf("expected 3 CBZ files, got %d", cbzCount)
	}
}

// TestIntegration_ExtensionFiltering verifies the -e flag filters extensions.
// BONUS: Not in original spec test cases but useful for coverage.
// Invariants verified: #1 (complete per filter), #6 (robustness)
func TestIntegration_ExtensionFiltering(t *testing.T) {
	// Setup: Create temp directory structure
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "manga")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Create chapter dir
	chapterPath := filepath.Join(inputDir, "Chapter1")
	if err := os.MkdirAll(chapterPath, 0755); err != nil {
		t.Fatalf("failed to create chapter dir: %v", err)
	}

	// Create image files with various extensions
	createTestImage(t, filepath.Join(chapterPath, "01.jpg"))
	createTestImage(t, filepath.Join(chapterPath, "02.png"))
	createTestImage(t, filepath.Join(chapterPath, "03.webp"))
	createTestImage(t, filepath.Join(chapterPath, "04.bmp"))

	// Run CLI with extension filter for jpg,png only (flags must come before positional args)
	stdout, stderr, exitCode := runCLI(t, "-e", "jpg,png", "--out", outputDir, inputDir)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	// Verify CBZ contains only jpg and png files
	cbzPath := filepath.Join(outputDir, "Chapter1.cbz")
	files := getCBZFileList(t, cbzPath)

	if len(files) != 2 {
		t.Errorf("expected 2 files (jpg, png only), got %d: %v", len(files), files)
	}

	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f))
		if ext != ".jpg" && ext != ".png" {
			t.Errorf("unexpected file extension in CBZ: %s", f)
		}
	}
}
