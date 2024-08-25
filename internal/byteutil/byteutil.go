package byteutil

import (
	"encoding/binary"
	"math"
)


func Max(a, b []byte) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return false
}

func Max4Bytes(a, b [4]byte) bool {
	if a[0] != b[0] {
		return a[0] < b[0]
	}
	if a[1] != b[1] {
		return a[1] < b[1]
	}
	if a[2] != b[2] {
		return a[2] < b[2]
	}
	return a[3] < b[3]
}

func Max2Byte(a, b [2]byte) bool {
	if a[0] != b[0] {
		return a[0] < b[0]
	}
	return a[1] < b[1]
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

func ConvertTo2Byte(value []byte) [2]byte {
	var result [2]byte
	copy(result[:], value)
	return result
}

func ConvertTo16Byte(value []byte) [16]byte {
	var result [16]byte
	copy(result[:], value)
	return result
}
