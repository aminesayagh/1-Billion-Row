package byteutil

import (
	"math/rand"
	"onBillion/internal/tracker"
	"testing"
	"strconv"
	"strings"
	"fmt"
)

func maxInt(a, b int32) bool {
	return a < b
}

func maxInt16(a, b int16) bool {
	return a < b
}

func maxFloat64(a, b float64) bool {
	return a < b
}

func ConvertByteToInt32(b []byte) (int16, error) {
	str := strings.TrimRight(string(b), ".")
	str = strings.SplitN(str, ".", 2)[0]

	// Convert string to integer using strconv.Atoi (as we're handling integers)
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("error converting %s to int: %v", str, err)
	}

	if i < -32768 || i > 32767 {
		return 0, fmt.Errorf("integer out of int16 range: %d", i)
	}

	return int16(i), nil
}

func convertByteTo8ByteNumber(b []byte) [8]byte {
	var result [8]byte
	copy(result[:], b)
	return result
}



func TestMaxPerformance(t *testing.T) {
	randomStringNumbers := []string{
		"15.2", "57.7", "23.4", "12.9", "45.1", "78.6", "34.7", "-67.2", "89.1", "-23.6", "45.5", "67.4", "89.5", "23.2", "45.4", "67.6",
	}
	
	randomByteNumbers := make([][]byte, len(randomStringNumbers))
	for i, v := range randomStringNumbers {
		randomByteNumbers[i] = []byte(v)
	}
	randomBytes := func() []byte {
		return randomByteNumbers[rand.Intn(len(randomByteNumbers))]
	}
	fByte := func() {
		Max(randomBytes(), randomBytes())
	}
	f2byte := func() {
		a2Byte := convertByteTo2Byte(randomBytes())
		b2Byte := convertByteTo2Byte(randomBytes())
		Max2Byte(a2Byte, b2Byte)
	}
	f4byte := func() {
		a4Byte := convertByteTo4Byte(randomBytes())
		b4Byte := convertByteTo4Byte(randomBytes())
		Max4Bytes(a4Byte, b4Byte)
	}
	fInt16 := func() {
		aInt := ConvertByteToInt16(randomBytes())
		bInt := ConvertByteToInt16(randomBytes())
		_ = maxInt16(aInt, bInt)
	}
	fInt32 := func() {
		aInt := convertByteToInt32(randomBytes())
		bInt := convertByteToInt32(randomBytes())
		_ = maxInt(aInt, bInt)
	}
	fFloat64 := func() {
		aFloat := convertByteToFloat64(randomBytes())
		bFloat := convertByteToFloat64(randomBytes())
		_ = maxFloat64(aFloat, bFloat)
	}

	// Number of iterations
	iterations := 100000

	// Measure performance of MaxBytes vs maxInt
	funcs := []func(){
		fByte,    // Function 1
		f2byte,   // Function 2
		f4byte,   // Function 3
		fInt16,   // Function 4
		fInt32,   // Function 5
		fFloat64, // Function 6
	}

	tracker.MeasureExecutionTime(funcs, iterations)
}

func TestAddPerformance(t *testing.T) {

}
