// Package sort provides sorting utilities for manga file ordering.
package sort

import (
	"sort"
	"unicode"
)

// chunk represents a segment of a string, either numeric or non-numeric.
type chunk struct {
	value     string
	isNumeric bool
}

// Natural sorts strings in natural/alphanumeric order in-place.
// This ensures "10" comes after "9", not after "1".
// The sort is stable: equal elements maintain their relative order.
func Natural(items []string) {
	sort.SliceStable(items, func(i, j int) bool {
		return naturalLess(items[i], items[j])
	})
}

// naturalLess returns true if a should come before b in natural order.
func naturalLess(a, b string) bool {
	chunksA := splitChunks(a)
	chunksB := splitChunks(b)

	minLen := len(chunksA)
	if len(chunksB) < minLen {
		minLen = len(chunksB)
	}

	for i := 0; i < minLen; i++ {
		cmp := compareChunks(chunksA[i], chunksB[i])
		if cmp != 0 {
			return cmp < 0
		}
	}

	// All compared chunks are equal; shorter string comes first
	return len(chunksA) < len(chunksB)
}

// splitChunks divides a string into alternating numeric and non-numeric segments.
func splitChunks(s string) []chunk {
	if s == "" {
		return nil
	}

	var chunks []chunk
	runes := []rune(s)
	start := 0

	for start < len(runes) {
		isDigit := unicode.IsDigit(runes[start])
		end := start + 1

		for end < len(runes) && unicode.IsDigit(runes[end]) == isDigit {
			end++
		}

		chunks = append(chunks, chunk{
			value:     string(runes[start:end]),
			isNumeric: isDigit,
		})
		start = end
	}

	return chunks
}

// compareChunks compares two chunks and returns:
// -1 if a < b, 0 if a == b, 1 if a > b.
func compareChunks(a, b chunk) int {
	// If both are numeric, compare by numeric value
	if a.isNumeric && b.isNumeric {
		return compareNumeric(a.value, b.value)
	}

	// If types differ, numeric chunks come before non-numeric
	if a.isNumeric != b.isNumeric {
		if a.isNumeric {
			return -1
		}
		return 1
	}

	// Both are non-numeric: lexicographic comparison
	return compareStrings(a.value, b.value)
}

// compareNumeric compares two numeric strings by their integer value.
// For equal values with different representations (001 vs 1),
// shorter strings come first.
func compareNumeric(a, b string) int {
	// Strip leading zeros for comparison
	aStripped := stripLeadingZeros(a)
	bStripped := stripLeadingZeros(b)

	// Compare by length first (longer = bigger number)
	if len(aStripped) != len(bStripped) {
		if len(aStripped) < len(bStripped) {
			return -1
		}
		return 1
	}

	// Same length, compare lexicographically (works for same-length numbers)
	cmp := compareStrings(aStripped, bStripped)
	if cmp != 0 {
		return cmp
	}

	// Equal numeric value; shorter original string comes first (fewer leading zeros)
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}
	return 0
}

// stripLeadingZeros removes leading zeros from a numeric string.
// Returns "0" for all-zero strings.
func stripLeadingZeros(s string) string {
	i := 0
	for i < len(s)-1 && s[i] == '0' {
		i++
	}
	return s[i:]
}

// compareStrings returns -1, 0, or 1 for string comparison.
func compareStrings(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
