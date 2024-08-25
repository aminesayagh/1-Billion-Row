package byteutil

import (
	"testing"
)



func TestMaxByte(t *testing.T) {
	a := convertInt32ToByte(13)
	b := convertInt32ToByte(24)

	if !Max(a, b) {
		t.Errorf("Expected Max(a, b) to be true, got false")
	}

	if Max(b, a) {
		t.Errorf("Expected Max(b, a) to be false, got true")
	}
}

func TestMax2Byte(t *testing.T) {
	a := convertInt32To2Byte(13)
	b := convertInt32To2Byte(24)

	if !Max2Byte(a, b) {
		t.Errorf("Expected Max2Byte(a, b) to be true, got false")
	}

	if Max2Byte(b, a) {
		t.Errorf("Expected Max2Byte(b, a) to be false, got true")
	}
}



func TestAdd(t *testing.T) {
	a := convertInt32To4Byte(4)
	b := convertInt32To4Byte(12)
	expected := convertInt32To4Byte(16)

	result := Add(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected Add(a, b) to be %v, got %v", expected, result)
		}
	}
}

func TestDivide(t *testing.T) {
	a := convertInt32To4Byte(12)
	b := convertInt32To4Byte(3)
	expected := convertInt32To4Byte(4)

	result := Divide(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected Divide(a, b) to be %v, got %v", expected, result)
		}
	}
}
