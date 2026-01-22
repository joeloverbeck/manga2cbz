package convert

import (
	"encoding/base64"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"manga2cbz/internal/chapter"
)

// Minimal valid WebP image (1x1 pixel, red)
// Created from a real WebP encoder - this is the smallest valid WebP file
var minimalWebP = mustDecodeBase64("UklGRlYAAABXRUJQVlA4IEoAAADQAQCdASoBAAEAAUAmJYgCdAEO/hOMAAD++O9P/p3mX6w1v/xp/7H/lf+j/2b+Gf/Z/sV/on+x+AD+af5R/tv/K/1r2AP1O/7v/nfYAAAA")

func mustDecodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic("invalid base64 test data: " + err.Error())
	}
	return data
}

func TestIsWebP(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"lowercase webp", "image.webp", true},
		{"uppercase WEBP", "image.WEBP", true},
		{"mixed case WebP", "image.WebP", true},
		{"jpg extension", "image.jpg", false},
		{"jpeg extension", "image.jpeg", false},
		{"png extension", "image.png", false},
		{"no extension", "image", false},
		{"webp in name", "webp_image.jpg", false},
		{"double extension", "image.jpg.webp", true},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWebP(tt.filename)
			if result != tt.expected {
				t.Errorf("isWebP(%q) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestConvertWebPImages_NoWebPFiles(t *testing.T) {
	// Create temp directory with non-WebP files
	tempDir := t.TempDir()

	// Create a simple PNG file
	pngPath := filepath.Join(tempDir, "test.png")
	createTestPNG(t, pngPath)

	jpgPath := filepath.Join(tempDir, "test.jpg")
	createTestFile(t, jpgPath, []byte("fake jpg content"))

	images := []chapter.ImageFile{
		{Path: pngPath, Name: "test.png"},
		{Path: jpgPath, Name: "test.jpg"},
	}

	result, cleanup, err := ConvertWebPImages(images)
	if err != nil {
		t.Fatalf("ConvertWebPImages returned error: %v", err)
	}
	defer cleanup()

	// Should return same slice unchanged
	if len(result) != len(images) {
		t.Errorf("result length = %d, want %d", len(result), len(images))
	}

	for i, img := range result {
		if img.Path != images[i].Path {
			t.Errorf("result[%d].Path = %q, want %q", i, img.Path, images[i].Path)
		}
		if img.Name != images[i].Name {
			t.Errorf("result[%d].Name = %q, want %q", i, img.Name, images[i].Name)
		}
	}
}

func TestConvertWebPImages_WithWebPFiles(t *testing.T) {
	// Create temp directory with WebP files
	tempDir := t.TempDir()

	// Create a WebP file
	webpPath := filepath.Join(tempDir, "page1.webp")
	createTestFile(t, webpPath, minimalWebP)

	// Create a non-WebP file
	pngPath := filepath.Join(tempDir, "page2.png")
	createTestPNG(t, pngPath)

	images := []chapter.ImageFile{
		{Path: webpPath, Name: "page1.webp"},
		{Path: pngPath, Name: "page2.png"},
	}

	result, cleanup, err := ConvertWebPImages(images)
	if err != nil {
		t.Fatalf("ConvertWebPImages returned error: %v", err)
	}
	defer cleanup()

	if len(result) != 2 {
		t.Fatalf("result length = %d, want 2", len(result))
	}

	// First image should be converted to PNG
	if result[0].Name != "page1.png" {
		t.Errorf("result[0].Name = %q, want %q", result[0].Name, "page1.png")
	}
	if filepath.Ext(result[0].Path) != ".png" {
		t.Errorf("result[0].Path extension = %q, want .png", filepath.Ext(result[0].Path))
	}

	// Verify converted file exists and is valid PNG
	verifyPNG(t, result[0].Path)

	// Second image should be unchanged
	if result[1].Name != "page2.png" {
		t.Errorf("result[1].Name = %q, want %q", result[1].Name, "page2.png")
	}
	if result[1].Path != pngPath {
		t.Errorf("result[1].Path = %q, want %q", result[1].Path, pngPath)
	}
}

func TestConvertWebPImages_CleanupRemovesTempFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create a WebP file
	webpPath := filepath.Join(tempDir, "test.webp")
	createTestFile(t, webpPath, minimalWebP)

	images := []chapter.ImageFile{
		{Path: webpPath, Name: "test.webp"},
	}

	result, cleanup, err := ConvertWebPImages(images)
	if err != nil {
		t.Fatalf("ConvertWebPImages returned error: %v", err)
	}

	// Get the temp directory where converted file was created
	convertedPath := result[0].Path
	convertedDir := filepath.Dir(convertedPath)

	// Verify converted file exists
	if _, err := os.Stat(convertedPath); os.IsNotExist(err) {
		t.Fatal("converted file should exist before cleanup")
	}

	// Call cleanup
	cleanup()

	// Verify temp directory is removed
	if _, err := os.Stat(convertedDir); !os.IsNotExist(err) {
		t.Error("temp directory should be removed after cleanup")
	}
}

func TestConvertWebPImages_InvalidWebP(t *testing.T) {
	tempDir := t.TempDir()

	// Create an invalid WebP file (not valid WebP data)
	invalidWebpPath := filepath.Join(tempDir, "invalid.webp")
	createTestFile(t, invalidWebpPath, []byte("not valid webp data"))

	images := []chapter.ImageFile{
		{Path: invalidWebpPath, Name: "invalid.webp"},
	}

	_, _, err := ConvertWebPImages(images)
	if err == nil {
		t.Error("ConvertWebPImages should return error for invalid WebP")
	}
}

func TestConvertWebPImages_NonExistentFile(t *testing.T) {
	images := []chapter.ImageFile{
		{Path: "/nonexistent/path/test.webp", Name: "test.webp"},
	}

	_, _, err := ConvertWebPImages(images)
	if err == nil {
		t.Error("ConvertWebPImages should return error for nonexistent file")
	}
}

func TestConvertWebPImages_EmptySlice(t *testing.T) {
	images := []chapter.ImageFile{}

	result, cleanup, err := ConvertWebPImages(images)
	if err != nil {
		t.Fatalf("ConvertWebPImages returned error: %v", err)
	}
	defer cleanup()

	if len(result) != 0 {
		t.Errorf("result length = %d, want 0", len(result))
	}
}

func TestConvertWebPImages_PreservesOrder(t *testing.T) {
	tempDir := t.TempDir()

	// Create mixed files
	webp1 := filepath.Join(tempDir, "01.webp")
	png2 := filepath.Join(tempDir, "02.png")
	webp3 := filepath.Join(tempDir, "03.webp")
	jpg4 := filepath.Join(tempDir, "04.jpg")

	createTestFile(t, webp1, minimalWebP)
	createTestPNG(t, png2)
	createTestFile(t, webp3, minimalWebP)
	createTestFile(t, jpg4, []byte("fake jpg"))

	images := []chapter.ImageFile{
		{Path: webp1, Name: "01.webp"},
		{Path: png2, Name: "02.png"},
		{Path: webp3, Name: "03.webp"},
		{Path: jpg4, Name: "04.jpg"},
	}

	result, cleanup, err := ConvertWebPImages(images)
	if err != nil {
		t.Fatalf("ConvertWebPImages returned error: %v", err)
	}
	defer cleanup()

	expected := []string{"01.png", "02.png", "03.png", "04.jpg"}
	if len(result) != len(expected) {
		t.Fatalf("result length = %d, want %d", len(result), len(expected))
	}

	for i, img := range result {
		if img.Name != expected[i] {
			t.Errorf("result[%d].Name = %q, want %q", i, img.Name, expected[i])
		}
	}
}

// Helper functions

func createTestFile(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to create test file %s: %v", path, err)
	}
}

func createTestPNG(t *testing.T, path string) {
	t.Helper()
	// Create a 1x1 pixel image
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create PNG file %s: %v", path, err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("failed to encode PNG: %v", err)
	}
}

func verifyPNG(t *testing.T, path string) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("failed to open PNG file %s: %v", path, err)
	}
	defer f.Close()

	_, err = png.Decode(f)
	if err != nil {
		t.Errorf("file %s is not a valid PNG: %v", path, err)
	}
}
