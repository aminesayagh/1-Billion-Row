package v2_assembly

import (
	"reflect"
	"testing"
)

func TestBytesToNumericStruct(t *testing.T) {
	testCases := []struct {
		name          string
		input         []byte
		expectedDigits [6]byte
		expectedSign  byte
		expectedScale byte
		expectedLen   byte
		expectError   int
	}{
		{
			name:          "Positive number with decimals",
			input:         []byte("123.45"),
			expectedDigits: [6]byte{1, 2, 3, 4, 5, 0},
			expectedSign:  0,
			expectedScale: 2,
			expectedLen:   5,
			expectError:   0,
		},
		{
			name:          "Negative number with decimals",
			input:         []byte("-123.45"),
			expectedDigits: [6]byte{1, 2, 3, 4, 5, 0},
			expectedSign:  255,
			expectedScale: 2,
			expectedLen:   5,
			expectError:   0,
		},
		{
			name:          "Positive integer",
			input:         []byte("6789"),
			expectedDigits: [6]byte{6, 7, 8, 9, 0, 0},
			expectedSign:  0,
			expectedScale: 0,
			expectedLen:   4,
			expectError:   0,
		},
		{
			name:          "Negative integer",
			input:         []byte("-42"),
			expectedDigits: [6]byte{4, 2, 0, 0, 0, 0},
			expectedSign:  255,
			expectedScale: 0,
			expectedLen:   2,
			expectError:   0,
		},
		{
			name:          "Decimal only",
			input:         []byte("."),
			expectedDigits: [6]byte{},
			expectedSign:  0,
			expectedScale: 0,
			expectedLen:   0,
			expectError:   -1,
		},
		{
			name:          "Empty input",
			input:         []byte(""),
			expectedDigits: [6]byte{},
			expectedSign:  0,
			expectedScale: 0,
			expectedLen:   0,
			expectError:   0,
		},
		{
			name:          "Invalid input",
			input:         []byte("abc"),
			expectedDigits: [6]byte{},
			expectedSign:  0,
			expectedScale: 0,
			expectedLen:   0,
			expectError:   -1,
		},
		{
			name:          "Multiple decimal points",
			input:         []byte("123.45.67"),
			expectedDigits: [6]byte{1, 2, 3, 4, 5, 0},
			expectedSign:  0,
			expectedScale: 2,
			expectedLen:   5,
			expectError:   -4,
		},
		{
			name:          "Multiple negative signs",
			input:         []byte("--123.45"),
			expectedDigits: [6]byte{1, 2, 3, 4, 5, 0},
			expectedSign:  255,
			expectedScale: 2,
			expectedLen:   5,
			expectError:   -1,
		},
		{
			name:          "Multiple positive signs",
			input:         []byte("++123.45"),
			expectedDigits: [6]byte{1, 2, 3, 4, 5, 0},
			expectedSign:  0,
			expectedScale: 2,
			expectedLen:   5,
			expectError:   -1,
		},
		{
			name: 		"Close with 0",
			input: []byte("32.0"),
			expectedDigits: [6]byte{3, 2, 0, 0, 0, 0},
			expectedSign:  0,
			expectedScale: 0,
			expectedLen:   2,
			expectError:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var digits [6]byte
			var sign, scale, length byte
			
			errorCode := BytesToNumericBytes(tc.input, &digits, &sign, &scale, &length)

			if errorCode != tc.expectError {
				t.Errorf("Expected error code %d, but got %d", tc.expectError, errorCode)
			}

			if !reflect.DeepEqual(digits, tc.expectedDigits) {
				t.Errorf("Expected digits %v, but got %v", tc.expectedDigits, digits)
			}

			if sign != tc.expectedSign {
				t.Errorf("Expected sign %d, but got %d", tc.expectedSign, sign)
			}

			if scale != tc.expectedScale {
				t.Errorf("Expected scale %d, but got %d", tc.expectedScale, scale)
			}

			if length != tc.expectedLen {
				t.Errorf("Expected length %d, but got %d", tc.expectedLen, length)
			}

			t.Logf("Input: %s, Digits: %v, Sign: %d, Scale: %d, Length: %d, Error: %d",
				string(tc.input), digits, sign, scale, length, errorCode)
		})
	}
}
