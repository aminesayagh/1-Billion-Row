#include "textflag.h"

TEXT ·BytesToNumericBytes(SB), NOSPLIT, $0
    MOVQ    len+8(FP), CX      // Load the length of the byte slice into CX
    MOVQ    b+0(FP), AX        // Load the pointer to the byte slice into AX
    MOVQ    AX, DX             // Copy the pointer to DX for iteration
    JMP     loop               // Jump to the loop

loop:
    CMPQ    CX, $0             // Check if CX (length) is 0
    JE      done               // If so, we're done

    MOVB    (DX), AX           // Load the next byte from the slice into AX, zero-extended
    CMPB    AL, $'-'           // Compare the byte with '-'
    JE      handle_sign        // If it's '-', handle it as a sign
    CMPB    AL, $'+'           // Compare the byte with '+'
    JE      handle_sign        // If it's '+', handle it as a sign
    CMPB    AL, $'0'           // Compare the byte with '0'
    JL      next_char          // If less than '0', skip to the next character
    CMPB    AL, $'9'           // Compare the byte with '9'
    JG      next_char          // If greater than '9', skip to the next character

    SUBB    $'0', AL           // Convert the byte to a numeric value
    MOVB    AL, (DX)           // Store the converted value back in the slice

next_char:
    INCQ    DX                 // Move to the next byte in the slice
    DECQ    CX                 // Decrement the loop counter
    JMP     loop               // Repeat the loop

handle_sign:
    MOVB    AL, (DX)           // Directly copy the sign character
    JMP     next_char          // Move to the next character

done:
    MOVQ    DX, AX             // Return the pointer to the byte slice
    RET                        // Return from the function
