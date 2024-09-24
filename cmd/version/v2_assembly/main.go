package v2_assembly

import (
	"bufio"
	"bytes"
	"fmt"
	"oneBillion/config"
	"io"
	"os"
	"sync"
)

type Measurement struct {
	Min     DecimalNumber
	Max     DecimalNumber
	Average DecimalNumber
	Count   int32
}

// Assembly function declaration
func BytesToNumericBytes(input []byte, digits *[6]byte, sign *byte, scale *byte, length *byte) int

var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 64)
    },
}

var digits = [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func intToBytes(i int32, buf []byte) int {
    // Handle 0 specially
    if i == 0 {
        buf[0] = '0'
        return 1
    }

    // Handle negative numbers
    var negative bool
    if i < 0 {
        negative = true
        i = -i
    }

    // Convert the absolute value to bytes
    pos := len(buf)
    for i > 0 {
        pos--
        buf[pos] = digits[i%10]
        i /= 10
    }

    // Add the negative sign if necessary
    if negative {
        pos--
        buf[pos] = '-'
    }

    // Move the bytes to the beginning of the buffer
    copy(buf, buf[pos:])

    return len(buf) - pos
}

func DecimalNumberToBytes(i DecimalNumber, buf []byte) int {
    if i.Length == 0 {
        buf[0] = '0'
        buf[1] = '.'
        buf[2] = '0'
        return 3
    }

    pos := 0

    // Add negative sign if necessary
    if i.Sign == 255 {
        buf[pos] = '-'
        pos++
    }

    // Add digits before decimal point
    integerPartLength := int(i.Length) - int(i.Scale)
    for j := 0; j < integerPartLength; j++ {
        buf[pos] = digits[i.Digits[j]]
        pos++
    }

    // Add decimal point if there's a fractional part
    if i.Scale > 0 {
        buf[pos] = '.'
        pos++

        // Add digits after decimal point
        for j := integerPartLength; j < int(i.Length); j++ {
            buf[pos] = digits[i.Digits[j]]
            pos++
        }
    } else {
        // If no fractional part, add ".0"
        buf[pos] = '.'
        buf[pos+1] = '0'
        pos += 2
    }

    return pos
}


func Parsing(config *config.Config) {
	// Open the input file
	dataFile, err := os.Open(config.InputFilePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer dataFile.Close()

	measurements := make(map[string]*Measurement, 500000)
	temperatureCache := make(map[string]DecimalNumber, 10000)

	// Initialize the file scanner
    reader := bufio.NewReaderSize(dataFile, 64*1024) // 64 KB buffer

	var (
		line     		[]byte
		sepIndex        int
		station  		string
		tempBytes       []byte
		currentTemp     DecimalNumber
		measurement     *Measurement
		exists          bool
	)

	for {
        line = bufferPool.Get().([]byte)[:0]
        line, err = reader.ReadSlice('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Println("Error reading line:", err)
            continue
        }

		sepIndex = bytes.IndexByte(line, ';') // Split by semicolon without memory allocation
		if sepIndex == -1 {
			continue
		}

		station = string(line[:sepIndex]) // Extract the station as a string (key for the map)
		tempBytes = line[sepIndex+1:]
		
        // Use the temperature bytes as a key for caching
        tempKey := string(tempBytes)

		if currentTemp, exists = temperatureCache[tempKey]; !exists {
            errorCode := BytesToNumericBytes(tempBytes, &currentTemp.Digits, &currentTemp.Sign, &currentTemp.Scale, &currentTemp.Length)

			if errorCode  < 0 {
				fmt.Println("Error parsing temperature: ", errorCode, " for station: ", station)
				continue
			}
            temperatureCache[tempKey] = currentTemp
		}
		// Update measurements
		if measurement, exists = measurements[station]; exists {
			if currentTemp.Compare(measurement.Max) > 0 {
				measurement.Max = currentTemp
			}
			if currentTemp.Compare(measurement.Min) < 0 {
				measurement.Min = currentTemp
			}
			measurement.Count++
			// measurement.Average = 
		} else {
			measurements[station] = &Measurement{
				Min:     currentTemp,
				Max:     currentTemp,
				Average: currentTemp,
				Count:   1,
			}
		}
        bufferPool.Put(line)
		
	}

	// Write output to file
	outputFile, err := os.Create(config.OutputFilePath)
	if err != nil {
		fmt.Println("Error creating output file: ", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriterSize(outputFile, 4*1024*1024) // 4 MB buffer
    var buffer [256]byte // 256 bytes buffer
	
    for station, m := range measurements {
        n := 0
        n += copy(buffer[n:], station) // Copy the station name to the buffer
        buffer[n] = ';' // Add a semicolon
        n++
		
        n += DecimalNumberToBytes(m.Min, buffer[n:]) // Copy the minimum temperature to the buffer
        buffer[n] = ';'
        n++
        n += DecimalNumberToBytes(m.Max, buffer[n:]) // Copy the maximum temperature to the buffer
        buffer[n] = ';'
        n++
        n += DecimalNumberToBytes(m.Average, buffer[n:]) // Copy the average temperature to the buffer
        buffer[n] = ';'
        n++
		n += intToBytes(m.Count, buffer[n:]) // Convert the count to bytes and copy it to the buffer
        buffer[n] = '\n'
        n++

        if _, err := writer.Write(buffer[:n]); err != nil {
            fmt.Println("Error writing to output file:", err)
            return
        }
    }

	if err := writer.Flush(); err != nil {
		fmt.Println("Error flushing output buffer:", err)
	}
}
