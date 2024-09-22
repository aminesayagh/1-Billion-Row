package v2_assembly

import (
	"strconv"
)

func CalculateAverage(sum DecimalNumber, count int32) DecimalNumber {
    if count == 0 {
        return DecimalNumber{} // Return zero if count is zero
    }

    // Create a DecimalNumber representation of the count
    countDecimal := DecimalNumber{}
    countStr := strconv.Itoa(int(count))
    countDecimal.Normalize([]byte(countStr))

    // Perform division
    return Divide(sum, countDecimal)
}

// Divide divides one DecimalNumber by another
func Divide(a, b DecimalNumber) DecimalNumber {
    result := DecimalNumber{}
    
    // Handle division by zero
    if b.Length == 0 || (b.Length == 1 && b.Digits[0] == 0) {
        return result // Return zero for division by zero
    }

    // Set the sign of the result
    result.Sign = a.Sign ^ b.Sign

    // Prepare for long division
    numerator := make([]byte, len(a.Digits)+6) // Extra precision
    copy(numerator, a.Digits[:])
    
    divisor := make([]byte, len(b.Digits))
    copy(divisor, b.Digits[:])

    // Perform long division
    for i := 0; i < 6; i++ { // Up to 6 digits of precision
        digit := byte(0)
        for compareBytes(numerator, divisor) >= 0 {
            subtractBytes(numerator, divisor)
            digit++
        }
        result.Digits[i] = digit
        result.Length++
        
        // Shift numerator
        copy(numerator, numerator[1:])
        numerator[len(numerator)-1] = 0
    }

    // Adjust scale
    result.Scale = a.Scale - b.Scale + 6 // Adjust as needed

    // Normalize result
    for result.Length > 0 && result.Digits[0] == 0 {
        result.Digits = result.Digits[1:]
        result.Length--
        result.Scale--
    }

    return result
}

// Helper function to compare byte slices
func compareBytes(a, b []byte) int {
    for i := 0; i < len(a) && i < len(b); i++ {
        if a[i] != b[i] {
            if a[i] > b[i] {
                return 1
            }
            return -1
        }
    }
    if len(a) > len(b) {
        return 1
    }
    if len(a) < len(b) {
        return -1
    }
    return 0
}

// Helper function to subtract byte slices
func subtractBytes(a, b []byte) {
    borrow := byte(0)
    for i := len(a) - 1; i >= 0; i-- {
        diff := a[i] - borrow
        if i < len(b) {
            if diff < b[i] {
                diff += 10
                borrow = 1
            } else {
                borrow = 0
            }
            diff -= b[i]
        } else if borrow > 0 {
            if diff < borrow {
                diff += 10
                borrow = 1
            } else {
                borrow = 0
            }
            diff -= borrow
        }
        a[i] = diff
    }
}