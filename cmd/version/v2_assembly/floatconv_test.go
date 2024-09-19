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
		expectedLen int  
		expectError int
	}{
		{
			name:        "Positive number with decimals",
			input:       []byte("123.45"),
			expected:    []byte{1, 2, 3, '.', 4, 5},
			expectedLen: 6,
			expectError: 0,
		},
		{
			name:        "Negative number with decimals",
			input:       []byte("-123.45"),
			expected:    []byte{'-', 1, 2, 3, '.', 4, 5},
			expectedLen: 7,
			expectError: 0,
		},
		{
			name:        "Positive integer with decimals",
			input:       []byte("+123.45"),
			expected:    []byte{1, 2, 3, '.', 4, 5},
			expectedLen: 6,
			expectError: 0,
		},
		{
			name:        "Positive integer",
			input:       []byte("6789."),
			expected:    []byte{6, 7, 8, 9, '.'},
			expectedLen: 5,
			expectError: 0,
		},
		{
			name:        "Negative integer",
			input:       []byte("-42"),
			expected:    []byte{'-', 4, 2},
			expectedLen: 3,
			expectError: 0,
		},
		{
			name:        "Positive number with leading spaces",
			input:       []byte("  123"),
			expected:    []byte{1, 2, 3},
			expectedLen: 3,
			expectError: 0,
		},
		{
			name:        "Positive number with trailing spaces",
			input:       []byte("123  "),
			expected:    []byte{1, 2, 3},
			expectedLen: 3,
			expectError: -1,
		},
		{
			name:        "Mixed characters with number",
			input:       []byte("abc123def"),
			expected:    []byte{1, 2, 3},
			expectedLen: 3,
			expectError: 0,
		},
		{
			name:        "Decimal only",
			input:       []byte("."),
			expected:    []byte{},
			expectedLen: -1, // Expect an error
			expectError: 0,
		},
		{
			name:        "Empty input",
			input:       []byte(""),
			expected:    []byte{},
			expectedLen: 0,
			expectError: 3,
		},
		{
			name:        "Invalid input",
			input:       []byte("abc"),
			expected:    []byte{},
			expectedLen: -1, // Expect an error
			expectError: 0,
		},
		{
			name:        "Multiple decimal points",
			input:       []byte("123.45.67"),
			expected:    []byte{1, 2, 3, '.', 4, 5},
			expectedLen: 6,
			expectError: -4,
		},
	}

	for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            output := make([]byte, len(tc.input)) // Allocate an output buffer with the same length as input
            errorCode := BytesToNumericBytes(tc.input, output)

            if errorCode != 0 && tc.expectError != 0 && errorCode != tc.expectError {
				t.Errorf("Expected error code %d, but got %d", tc.expectError, errorCode)
			} else if errorCode != 0 && tc.expectError == 0 {
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
