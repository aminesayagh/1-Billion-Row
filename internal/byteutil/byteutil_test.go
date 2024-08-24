package byteutil

import (
	"testing"
)

func TestMax(t *testing.T) {
	a := [4]byte{1, 2, 3, 4}
	b := [4]byte{1, 2, 3, 5}

	if !Max(a, b) {
		t.Errorf("Expected Max(a, b) to be true, got false")
	}

	if Max(b, a) {
		t.Errorf("Expected Max(b, a) to be false, got true")
	}
}

func TestAdd(t *testing.T) {
	a := [4]byte{1, 2, 3, 4}
	b := [4]byte{4, 3, 2, 1}
	expected := [4]byte{5, 5, 5, 5}

	result := Add(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected Add(a, b) to be %v, got %v", expected, result)
		}
	}
}

func TestDivide(t *testing.T) {
	a := [4]byte{8, 4, 6, 2}
	b := [4]byte{2, 2, 2, 2}
	expected := [4]byte{4, 2, 3, 1}

	result := Divide(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected Divide(a, b) to be %v, got %v", expected, result)
		}
	}
}

