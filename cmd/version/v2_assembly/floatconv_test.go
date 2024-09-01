package v2_assembly

import (
	"reflect" // Added reflect package used for deep equal comparison
	"testing"
)

func TestBytesToNumericBytes(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Positive number with decimals",
			input:    []byte("123.45"),
			expected: []byte{1, 2, 3, '.', 4, 5},
		},
		{
			name:     "Negative number with decimals",
			input:    []byte("-123.45"),
			expected: []byte{'-', 1, 2, 3, '.', 4, 5},
		},
		{
			name:     "Positive integer",
			input:    []byte("6789"),
			expected: []byte{6, 7, 8, 9},
		},
		{
			name:     "Negative integer",
			input:    []byte("-42"),
			expected: []byte{'-', 4, 2},
		},
		{
			name:     "Positive number with leading and trailing spaces",
			input:    []byte("  123  "),
			expected: []byte{1, 2, 3, 0, 0, 0, 0},
		},
		// {
		// 	name:     "Mixed characters with number",
		// 	input:    []byte("abc123def"),
		// 	expected: []byte{1, 2, 3, 0, 0, 0, 0},
		// },
		// {
		// 	name:     "Decimal only",
		// 	input:    []byte("."),
		// 	expected: []byte{'.'},
		// },
		// {
		// 	name:     "Empty input",
		// 	input:    []byte(""),
		// 	expected: []byte{},
		// },
		// {
		// 	name:     "Negative number with leading and trailing spaces",
		// 	input:    []byte("  -123  "),
		// 	expected: []byte{'-', 1, 2, 3, 0, 0, 0, 0},
		// },
		// {
		// 	name:     "Negative number with leading and trailing spaces and decimals",
		// 	input:    []byte("  -123.45  "),
		// 	expected: []byte{'-', 1, 2, 3, '.', 4, 5},
		// },
		// {
		// 	name:    "Last characters are a '.'",
		// 	input:   []byte("123."),
		// 	expected: []byte{1, 2, 3},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Input: %v", tc.input)
			BytesToNumericBytes(tc.input) 
			
			if !reflect.DeepEqual(tc.input, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, tc.input)
			} else {
				t.Logf("Expected %v, and got %v", tc.expected, tc.input)
			}
		})
	}
}
