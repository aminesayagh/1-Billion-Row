package v2_assembly

import (
	"unsafe"
)

type DecimalNumber struct {
	Digits [6]byte // Fixed-size array for digits
	Sign   byte    // 0 for positive, 255 for negative
	Scale  byte    // Number of digits after the decimal point
	Length byte    // Number of significant digits
}

func (bn *DecimalNumber) Normalize(num []byte) {
	bn.Sign = 0
	bn.Scale = 0
	bn.Length = 0
	decimalPointSeen := false

	for i := 0; i < len(num) && bn.Length < 6; i++ {
		c := num[i]
		switch c {
		case '-':
			bn.Sign = 255
		case '.':
			decimalPointSeen = true
		default:
			bn.Digits[bn.Length] = c  // Store the byte value directly
			bn.Length++
			if decimalPointSeen {
				bn.Scale++
			}
		}
	}
}


func (bn *DecimalNumber) Add(b DecimalNumber) DecimalNumber {
	result := DecimalNumber{
		Sign:  bn.Sign,
		Scale: max(bn.Scale, b.Scale),
	}

	carry := byte(0)
	for i := 5; i >= 0; i-- {
		sum := carry
		if i >= 6-bn.Length {
			sum += bn.Digits[i-(6-bn.Length)]
		}
		if i >= 6-b.Length {
			sum += b.Digits[i-(6-b.Length)]
		}
		result.Digits[i] = sum % 10
		carry = sum / 10
	}

	// Adjust length
	result.Length = 6
	for result.Length > 0 && result.Digits[6-result.Length] == 0 {
		result.Length--
	}
	if result.Length == 0 {
		result.Length = 1
		result.Sign = 0
	}

	return result
}

// Compare compares two DecimalNumbers. Returns:
//   -1 if bn < b
//    0 if bn == b
//    1 if bn > b
func (bn *DecimalNumber) Compare(b DecimalNumber) int {
	if bn.Sign != b.Sign {
		if bn.Sign == 0 {
			return 1
		}
		return -1
	}

	for i := 0; i < 6; i++ {
		if bn.Digits[i] != b.Digits[i] {
			if bn.Digits[i] > b.Digits[i] {
				return 1
			}
			return -1
		}
	}

	return 0
}

func (bn *DecimalNumber) String() string {
	if bn.Length == 0 {
		return "0"
	}

	result := make([]byte, 0, 8) // max 6 digits + sign + decimal point

	if bn.Sign == 255 {
		result = append(result, '-')
	}

	integerPart := int(bn.Length) - int(bn.Scale)
	for i := 0; i < int(bn.Length); i++ {
		if i == integerPart && bn.Scale > 0 {
			result = append(result, '.')
		}
		result = append(result, bn.Digits[6-bn.Length+i]+'0')
	}

	return *(*string)(unsafe.Pointer(&result))
}