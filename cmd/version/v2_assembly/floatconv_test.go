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
			name:    "Positive integer with decimals",
			input:   []byte("+123.45"),
			expected: []byte{1, 2, 3, '.', 4, 5},
		},
		{
			name:     "Positive integer",
			input:    []byte("6789."),
			expected: []byte{6, 7, 8, 9, '.'},
		},
		{
			name:     "Negative integer",
			input:    []byte("-42"),
			expected: []byte{'-', 4, 2},
		},
		{
			name:     "Positive number with leading spaces before the number",
			input:    []byte("  123"),
			expected: []byte{1, 2, 3},
		},
		{
			name:     "Positive number with trailing spaces after the number",
			input:    []byte("123  "),
			expected: []byte{1, 2, 3},
		},
		{
			name:     "Positive number with leading and trailing spaces",
			input:    []byte("  123  "),
			expected: []byte{1, 2, 3},
		},
		{
			name:     "Mixed characters with number",
			input:    []byte("abc123def"),
			expected: []byte{1, 2, 3},
		},
		{
			name:     "Decimal only",
			input:    []byte("."),
			expected: []byte{46},
		},
		{
			name:     "Negative number with leading and trailing spaces",
			input:    []byte("  -123  "),
			expected: []byte{'-', 1, 2, 3},
		},
		{
			name:     "Negative number with leading and trailing spaces and decimals",
			input:    []byte("  -123.45  "),
			expected: []byte{'-', 1, 2, 3, '.', 4, 5},
		},
		{
			name:    "Last characters are a '.'",
			input:   []byte("123."),
			expected: []byte{1, 2, 3, 46},
		},
		{
			name:    "Last characters are a '.' with leading spaces",
			input:   []byte("  123."),
			expected: []byte{1, 2, 3, 46},
		},
		{
			name:    "Last characters are a '.' with trailing spaces",
			input:   []byte("123.  "),
			expected: []byte{1, 2, 3, 46},
		},
	}

	for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            output := make([]byte, len(tc.input)) // Allocate an output buffer with the same length as input
            errorCode := BytesToNumericBytes(tc.input, output)

            if errorCode != 0 {
                t.Errorf("Expected error code 0, but got %d", errorCode)
            }

            // Determine the actual length of valid bytes written
            actualLength := 0
            for i := 0; i < len(output); i++ {
                if output[i] != 0 {
                    actualLength++
                } else {
                    break
                }
            }

            actual := output[:actualLength]

            if !reflect.DeepEqual(actual, tc.expected) {
                t.Errorf("Expected %v, but got %v", tc.expected, actual)
            } else {
                t.Logf("Expected %v, and got %v", tc.expected, actual)
            }
        })
    }
}
