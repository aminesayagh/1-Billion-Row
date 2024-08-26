package hashutil


import (
	"hash/fnv"
)

type DefaultHash uint32

func HashBytes(data []byte) uint32 {
	h := fnv.New32() // 32-bit FNv-1a hash
	h.Write(data)
	return h.Sum32()
}

func ConvertSliceToArray(data []byte) [8]byte {
	var result [8]byte
	copy(result[:], data)
	return result
}