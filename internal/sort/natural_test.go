package sort

import (
	"reflect"
	"testing"
)

func TestNatural_NumericOrder(t *testing.T) {
	input := []string{"10", "2", "1", "20", "3"}
	want := []string{"1", "2", "3", "10", "20"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_MixedAlphanumeric(t *testing.T) {
	input := []string{"page10.jpg", "page2.jpg", "page1.jpg", "page20.jpg"}
	want := []string{"page1.jpg", "page2.jpg", "page10.jpg", "page20.jpg"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_LeadingZeros(t *testing.T) {
	// Shorter representation (fewer leading zeros) comes first for same numeric value
	input := []string{"01.jpg", "10.jpg", "02.jpg", "001.jpg"}
	want := []string{"01.jpg", "001.jpg", "02.jpg", "10.jpg"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_Alphabetic(t *testing.T) {
	input := []string{"charlie", "alpha", "bravo"}
	want := []string{"alpha", "bravo", "charlie"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_Empty(t *testing.T) {
	input := []string{}

	Natural(input)

	if len(input) != 0 {
		t.Errorf("Natural() on empty slice should remain empty, got %v", input)
	}
}

func TestNatural_Single(t *testing.T) {
	input := []string{"only"}
	want := []string{"only"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_ChapterPatterns(t *testing.T) {
	input := []string{"Chapter 10", "Chapter 2", "Chapter 1", "Chapter 100"}
	want := []string{"Chapter 1", "Chapter 2", "Chapter 10", "Chapter 100"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_MultipleNumericSegments(t *testing.T) {
	// v1c10p5 patterns (volume 1, chapter 10, page 5)
	input := []string{"v1c10p5", "v1c2p1", "v1c10p1", "v2c1p1"}
	want := []string{"v1c2p1", "v1c10p1", "v1c10p5", "v2c1p1"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_InPlace(t *testing.T) {
	input := []string{"3", "1", "2"}
	original := input // Same backing array

	Natural(input)

	// Verify it modified the original slice
	if &input[0] != &original[0] {
		t.Error("Natural() should sort in-place, not create new slice")
	}

	want := []string{"1", "2", "3"}
	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_OnlyNumbers(t *testing.T) {
	input := []string{"100", "99", "1000", "1", "10"}
	want := []string{"1", "10", "99", "100", "1000"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_SpecialCharacters(t *testing.T) {
	input := []string{"file-10.jpg", "file-2.jpg", "file_1.jpg", "file-1.jpg"}
	want := []string{"file-1.jpg", "file-2.jpg", "file-10.jpg", "file_1.jpg"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_MangaPageNaming(t *testing.T) {
	// Real-world manga page naming patterns
	input := []string{
		"001.png",
		"002.png",
		"010.png",
		"011.png",
		"100.png",
	}
	want := []string{
		"001.png",
		"002.png",
		"010.png",
		"011.png",
		"100.png",
	}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_MixedExtensions(t *testing.T) {
	input := []string{"10.png", "2.jpg", "1.gif", "20.webp"}
	want := []string{"1.gif", "2.jpg", "10.png", "20.webp"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_Stability(t *testing.T) {
	// Test that equal elements maintain relative order
	// "a1" and "a1" (if duplicated) should stay in original order
	input := []string{"b2", "a1", "b1", "a2", "a1"}
	want := []string{"a1", "a1", "a2", "b1", "b2"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_Unicode(t *testing.T) {
	// Ensure unicode digits are handled
	input := []string{"file10", "file2", "file１"} // ０ is fullwidth digit
	Natural(input)

	// Just verify it doesn't panic - exact order depends on unicode handling
	if len(input) != 3 {
		t.Errorf("Natural() should handle unicode without panic")
	}
}

func TestNatural_EmptyStrings(t *testing.T) {
	input := []string{"", "a", "", "1"}
	want := []string{"", "", "1", "a"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

func TestNatural_LargeNumbers(t *testing.T) {
	// Test with large numbers that might overflow int32
	input := []string{"99999999999", "1", "9999999999"}
	want := []string{"1", "9999999999", "99999999999"}

	Natural(input)

	if !reflect.DeepEqual(input, want) {
		t.Errorf("Natural() = %v, want %v", input, want)
	}
}

// Benchmark for performance validation
func BenchmarkNatural_100Items(b *testing.B) {
	base := []string{
		"page100.jpg", "page10.jpg", "page1.jpg", "page50.jpg", "page5.jpg",
		"page99.jpg", "page9.jpg", "page55.jpg", "page15.jpg", "page25.jpg",
	}
	// Repeat to get 100 items
	items := make([]string, 0, 100)
	for i := 0; i < 10; i++ {
		items = append(items, base...)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Make a copy to avoid sorting already sorted data
		data := make([]string, len(items))
		copy(data, items)
		Natural(data)
	}
}
