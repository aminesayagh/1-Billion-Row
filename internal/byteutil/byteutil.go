package byteutil

import (
	"encoding/binary"
	"math"
)

type DefaultByte [4]byte

// MaxBytes compares two DefaultByte values lexicographically and returns true
// if 'a' is less than 'b'. It acts like the `less` function in sorting.
func Max(a, b [4]byte) bool {
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return true
		} else if a[i] > b[i] {
			return false
		}
	}

	return false
}

func Add(a, b [4]byte) [4]byte {
	var result [4]byte
	var carry byte = 0 // carry is initially 0 for addition operation
	for i := 0; i < len(a); i++ {
		sum := a[i] + b[i] + carry
		result[i] = sum

		if a[i] > 0 && b[i] > 0 && sum < a[i] {
			carry = 1
		} else {
			carry = 0
		}
	}
	return result
}

func Divide(a, b [4]byte) [4]byte {
	var result [4]byte
	for i := 0; i < len(a); i++ {
		if b[i] == 0 {
			result[i] = 0 // or a[i], or 255, depending on your use case
		} else {
			result[i] = a[i] / b[i]
		}
	}
	return result
}

// ConvertByteToFloat64 converts a DefaultByte (8-byte array) into a float64.
func ConvertByteToFloat64(value [4]byte) float64 {
	bits := binary.LittleEndian.Uint64(value[:]) // convert the bytes to uint64
	return math.Float64frombits(bits)            // convert the bits to float64
}

func ConvertTo4Byte(value []byte) [4]byte {
	var result [4]byte
	copy(result[:], value)
	return result
}
