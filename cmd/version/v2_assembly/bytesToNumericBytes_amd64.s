#include "textflag.h"

TEXT ·BytesToNumericBytes(SB), NOSPLIT, $0
    // Register Usage:
    // CX: Input length
    // DX: Input pointer
    // SI: Output pointer
    // BX: Automaton state
    // R8: Sign pointer
    // R9: Scale pointer
    // R10: Length pointer
    // R11: Current digit count
    // R12: Current scale
    // R13: Error code

    // Function Parameters:
    // input+0(FP): Input pointer
    // input_len+8(FP): Input length
    // output+24(FP): Output pointer
    // output_len+32(FP): Output length
    // sign+40(FP): Sign pointer
    // scale+48(FP): Scale pointer
    // length+56(FP): Length pointer
    

    MOVQ    input_len+8(FP), CX            // Load the length of the byte slice into CX
    MOVQ    input+0(FP), DX                // Load the pointer to the byte slice into DX
    MOVQ    output+24(FP), SI              // Load output slice data pointer
    MOVQ    sign+32(FP), R8                // Load sign pointer
    MOVQ    scale+40(FP), R9               // Load scale pointer
    MOVQ    length+48(FP), R10             // Load length pointer

    XORQ    AX, AX                         // Clear AX to indicate no errors
    MOVQ    $0, R11                        // R11 will store current digit count
    MOVQ    $0, R12                        // R12 will store current scale
    MOVQ    $0, R13                        // R13 will store our error code

    MOVQ    $0, ret+56(FP)          // Initialize return value to 0 (no error)

    MOVB    $0, (R8)                // Initialize sign to positive (0)
    MOVB    $0, (R9)                // Initialize scale to 0
    MOVB    $0, (R10)               // Initialize length to 0

    CMPQ    CX, $0                  // Check if the input length is 0
    JE      done                    // If it is, we are done

    MOVQ    $0, BX                  // Set initial state to q0

main_loop:
    MOVB    (DX), AL                 // Load the current byte into AL

    // The most common
    CMPQ    BX, $3                   // If BX == 3, jump to q3
    JE      q3
    CMPQ    BX, $0                   // If BX == 0, jump to q0
    JE      q0
    CMPQ    BX, $1                   // If BX == 1, jump to q1
    JE      q1
    CMPQ    BX, $4                   // If BX == 4, jump to q4
    JE      q4
    CMPQ    BX, $6                   // If BX == 6, jump to q6
    JE      q6
    CMPQ    BX, $5                   // If BX == 5, jump to q5
    JE      q5
    CMPQ    BX, $2                   // If BX == 2, jump to q2
    JE      q2

    JMP     error_invalid_state      // If we get here, we are in an invalid state

// State q0: Initial state
q0:
    CMPB    AL, $'1'                 // Check if AL is '1'
    JB      q0_check_zero            // If less than '1', jump to q0_check_zero, JB is the same as JL but for unsigned

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JA      error_invalid_char       // If greater than '9', jump to error, JA is the same as JG but for unsigned

    JMP     set_state_q3              // Process the character

    q0_check_zero:
    CMPB    AL, $'0'                 // Check if AL is different than '0'
    JNE     q0_char_sign             // If less than '0', jump to q0_char_sign

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    JMP     set_state_q2             // Jump to q2

    q0_char_sign:
    CMPB    AL, $'-'                 // Check if AL is '-'
    JNE     error_invalid_char       // If equal, jump to q1

    MOVQ    $1, BX                   // Set the next state to q1
    MOVB    $255, (R8)               // Set the sign to negative
    JMP     skip_char       // If we get here, we are in an invalid state

    set_state_q2:
    MOVQ    $2, BX                   // Set the next state to q2
    JMP     process_next_byte        // Process the next byte

    set_state_q3:
    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value

    MOVQ    $3, BX                   // Set the next state to q3
    JMP     process_next_byte        // Process the next byte


// State q1: Sign state
q1:
    CMPB    AL, $'1'                 // Check if AL is less than '1'
    JB      q1_check_zero            // If less than '1', jump to q1_check_zero, (Jump if Bellow)

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JA      skip_char                // If greater than '9', jump to skip_char, (Jump if Above)

    JMP     set_state_q3             // Process the character

    q1_check_zero:
    CMPB    AL, $'0'                 // Check if AL is '0'
    JNE     error_invalid_char       // If not equal, jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $2, BX                   // Set the next state to q2

    JMP     process_next_byte        // Process the next byte

// State q2: Sign state
q2:
    CMPB    AL, $'.'                // Check if AL is '.'
    JNE     error_invalid_char      // an invalid state

    set_state_q4:
    MOVQ    $4, BX                   // Set the next state to q4

    JMP     skip_char                // Process the next byte

// State q3: Integer No-Zero state
q3:
    CMPB    AL, $'0'                 // Check if AL is less than '0'
    JB      q3_check_point           // If less than '0', skip the character

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JA      skip_char                // If greater than '9', skip the character

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $3, BX                   // Set the next state to q3

    JMP     process_next_byte        // Process the next byte

    q3_check_point:
    CMPB    AL, $'.'                 // Check if AL is '.'
    JE      set_state_q4             // If equal, jump to q5

    JMP     error_invalid_char       // If we get here, we are in an invalid state

// State q4: Decimal point state
q4:
    CMPB    AL, $'1'                 // Check if AL is less than '1'
    JB      q4_check_zero           // If less than '1', jump to q4_not_decimal

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JA      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $6, BX                   // Set the next state to q5

    INCQ    R12
    MOVB    R12, (R9)              // Update scale

    JMP     process_next_byte        // Process the next byte

    q4_check_zero: 
    CMPB    AL, $'0'                 // Check if AL is '0'
    JNE     q4_close                // If equal, jump to error

    // check if this is the last digit in the buffer
    CMPQ    CX, $1                   // Check if we've reached the end of the buffer
    JE      skip_char                // If we have, skip the character

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $5, BX                   // Set the next state to q5

    INCQ    R12
    MOVB    R12, (R9)              // Update scale
    
    JMP     process_next_byte        // Process the next byte

    q4_close:
    MOVQ    $0, CX                   // Set the length to 0
    JMP     accept                   // Accept the number

// State q5: Decimal state
q5:
    INCQ    R12
    MOVB    R12, (R9)              // Update scale

    CMPB    AL, $'1'                 // Check if AL is less than '1'
    JB      q4_check_zero           // If less than '1', jump to q5_not_decimal

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JA      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $6, BX                   // Set the next state to q6

    
    JMP     process_next_byte        // Process the next byte

// State q6: Decimal state
q6:

    CMPB    AL, $'1'                 // Check if AL is less than '1'
    JB      q4_check_zero            // If less than '1', jump to q6_not_decimal

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JA      skip_char                // If greater than '9', jump to error
    
    INCQ    R12
    MOVB    R12, (R9)               // Update scale

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $6, BX                   // Set the next state to q6

    JMP     process_next_byte        // Process the next byte

process_next_byte:
    CMPQ    R11, $6                  // Check if we've reached max digits
    JGE     error_buffer_overflow    // If we've reached or passed the end, error out
    
    MOVB    AL, (SI)(R11*1)          // Store digit in array
    INCQ    R11                      // Increment digit count
    MOVB    R11B, (R10)              // Update length

    JMP     skip_char                // Skip the character

error_invalid_char: // code -1 is 255 in uint8
    MOVQ    $-1, R13          // Set error code for invalid character
    JMP     save_error_and_continue

error_invalid_state: // code -2 is 254 in uint8
    MOVQ    $-2, R13          // Set error code for invalid state
    JMP     save_error_and_continue

error_unexpected_decimal: // code -3 is 253 in uint8
    MOVQ    $-3, R13          // Set error code for unexpected decimal point
    JMP     save_error_and_continue

save_error_and_continue:
    MOVQ    R13, ret+56(FP)   // Return the error code

skip_char:
    INCQ    DX                       // Move the source pointer
    DECQ    CX                       // Decrement the length
    CMPQ    CX, $0                   // Check if the length is 0
    JE      accept                   // If it is, we are done
    JMP     main_loop                // Otherwise, continue the loop

accept:
    CMPQ    CX, $0                   // Check if the length is 0
    JNE     error_invalid_char       // If it is not, we are done
    JMP     done

done:
    RET

error_multiple_decimals: // code -4 is 252 in uint8
    MOVQ    $-4, R13          // Set error code for multiple decimal points
    JMP     save_error
    
error_buffer_overflow:
    MOVQ    $-5, R13          // Set error code for buffer overflow
    JMP     save_error

save_error:
    MOVQ    R13, ret+56(FP)   // Return the error code
    RET

no_error:
    MOVQ    R11, ret+56(FP)   // If no error, return the number of digits processed
    RET
