package hashutil

import (
	"testing"	
	"oneBillion/internal/tracker"
)

// TestHashBytes tests the HashBytes function
func TestHashBytes(t *testing.T) {
	data := []byte("Hello, World!")
	hash := HashBytes(data)
	if hash != 0x4291a886 {
		t.Errorf("Hash mismatch: expected 0x4291a886, got %#x", hash)
	}
}

func TestHashBytesEmpty(t *testing.T) {
	data := []byte("")
	hash := HashBytes(data)
	if hash != 0 {
		t.Errorf("Hash mismatch: expected 0, got %#x", hash)
	}
}

func TestHashBytesPerformance(t *testing.T) {
	data := []byte("Hello, World!")
	fBytes := func() {
		var result [16]byte
		copy(result[:], data)
	}
	fHash := func() {
		var _ = HashBytes(data)
	}
	fString := func () {
		var _ = string(data)
	}

	funcs := make(map[string]func())
	funcs["ConvertSliceToArray"] = fBytes
	funcs["HashBytes"] = fHash
	funcs["String"] = fString
	tracker.MeasureExecutionTime(funcs, 1000000)
}



