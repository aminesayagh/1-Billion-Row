package mapStruct

import (
	"onBillion/internal/byteutil"
	"onBillion/internal/hashutil"
)

type Measurement struct {
	Name 	[]byte
	Min     [4]byte
	Max     [4]byte
	Average [4]byte
	Count   [4]byte
}


type Map struct {
	data map[uint32]*Measurement
}

func New() *Map {
	return &Map{
		data: make(map[uint32]*Measurement, 100000),
	}
}

func (m *Map) Add(key []byte, temperature [4]byte) {
	measurement := m.Get(key)
	// convert the temperature to DefaultByte
	if measurement == nil {
		measurement = &Measurement{
			Name:    key,
			Min:     temperature,
			Max:     temperature,
			Average: temperature,
			Count:   [4]byte{1},
		}
	} else {
		if byteutil.Max(measurement.Max, temperature) {
			measurement.Max = temperature
		}
		if byteutil.Max(temperature, measurement.Min) {
			measurement.Min = temperature
		}
		measurement.Average = byteutil.Divide(byteutil.Add(measurement.Average, temperature), measurement.Count)
		measurement.Count = byteutil.Add(measurement.Count, [4]byte{1})
	}
	m.Set(key, measurement)
}

func (m *Map) Set(key []byte, measurement *Measurement) {
	hash := hashutil.HashBytes(key)
	m.data[hash] = measurement
}

func (m *Map) Get(key []byte) (*Measurement) {
	hash := hashutil.HashBytes(key)
	return m.data[hash];
}

func (m *Map) List() map[uint32]*Measurement {
	return m.data
}