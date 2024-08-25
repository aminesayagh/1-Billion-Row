package byteutil

import (
	"encoding/binary"
)

func convertInt32To4Byte(value int32) [4]byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(value))
	return buf
}

func convertInt32ToByte(value int32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(value))
	return buf[:]
}

func convertInt32To2Byte(value int32) [2]byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(value))
	return buf
}

func convert4ByteToInt32(value [4]byte) int32 {
	return int32(binary.LittleEndian.Uint32(value[:]))
}

func convertByteToInt32(value []byte) int32 {
	return int32(binary.LittleEndian.Uint32(value))
}
func ConvertByteToInt16(value []byte) int16 {
	return int16(binary.LittleEndian.Uint16(value))
}

func convert2ByteToInt32(value [2]byte) int32 {
	return int32(binary.LittleEndian.Uint32(value[:]))
}

func convertFloat64To4Byte(value float64) [4]byte {
	return convertInt32To4Byte(int32(value))
}

func convert4ByteToFloat64(value [4]byte) float64 {
	return float64(convert4ByteToInt32(value))
}

func convert2ByteToFloat64(value [2]byte) float64 {
	return float64(convert2ByteToInt32(value))
}

func convertByteToFloat64(value []byte) float64 {
	return float64(convertByteToInt32(value))
}

func convertByteTo16Byte(value []byte) [16]byte {
	var result [16]byte
	copy(result[:], value)
	return result
}

func ConvertByteTo8Byte(value []byte) [8]byte {
	var result [8]byte
	copy(result[:], value)
	return result
}

func ConvertInt16ToByte(value int16) []byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(value))
	return buf[:]
}

func ConvertInt32ToByte(value int32) []byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(value))
	return buf[:]
}

func convertByteTo4Byte(value []byte) [4]byte {
	var result [4]byte
	copy(result[:], value)
	return result
}

func convertByteTo2Byte(value []byte) [2]byte {
	var result [2]byte
	copy(result[:], value)
	return result
}
