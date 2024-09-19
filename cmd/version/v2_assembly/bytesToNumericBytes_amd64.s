#include "textflag.h"

TEXT ·BytesToNumericBytes(SB), NOSPLIT, $0
    // Register Usage:
    // CX: Input length
    // DX: Input pointer
    // SI: Output pointer
    // BX: Automaton state
    // R8: Output length
    // R9: End of output buffer
    // R10: Error code

    // Function Parameters:
    // input+0(FP): Input pointer
    // input_len+8(FP): Input length
    // output+24(FP): Output pointer
    // output_len+32(FP): Output length
    // ret+48(FP): Return value

    MOVQ    input_len+8(FP), CX            // Load the length of the byte slice into CX
    MOVQ    input+0(FP), DX                // Load the pointer to the byte slice into DX
    MOVQ    output+24(FP), SI              // Load output slice data pointer
    XORQ    AX, AX                         // Clear AX to indicate no errors

    MOVQ    output_len+32(FP), R8          // Load output slice length
    MOVQ    output+24(FP), R9              // Load output slice start address
    ADDQ    R8, R9                         // R9 now points to end of output slice

    MOVQ    $0, R10                        // R10 will store our error code

    MOVQ    $0, ret+48(FP)                 // Initialize return value to 0 (to error)

    CMPQ    CX, $0                         // Check if the length is 0
    JE      done                           // If it is, we are done

    MOVQ    $0, BX                         // Set initial state to q0

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
    JA      skip_char                // If greater than '9', jump to error, JA is the same as JG but for unsigned

    JMP     set_state_q3              // Process the character

    q0_check_zero:
    CMPB    AL, $'0'                 // Check if AL is different than '0'
    JNE     q0_char_sign             // If less than '0', jump to q0_char_sign

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    JMP     set_state_q2             // Jump to q2

    q0_char_sign:
    CMPB    AL, $'-'                 // Check if AL is '-'
    JE      set_state_q1             // If equal, jump to q1

    CMPB    AL, $'+'                 // Check if AL is '+'
    JE      skip_char             // If equal, jump to q1

    JMP     skip_char       // If we get here, we are in an invalid state

    set_state_q1:
    MOVQ    $1, BX                   // Set the next state to q1
    JMP     process_next_byte        // Process the next byte

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
    JE      set_state_q4             // If equal, jump to q4

    JMP     error_invalid_char       // If we get here, we are in an invalid state

    set_state_q4:
    MOVQ    $4, BX                   // Set the next state to q4
    JMP     process_next_byte        // Process the next byte

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

    JMP     process_next_byte        // Process the next byte

    q4_check_zero: 
    CMPB    AL, $'0'                 // Check if AL is '0'
    JNE     error_multiple_decimals  // If equal, jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $5, BX                   // Set the next state to q5
    
    JMP     process_next_byte        // Process the next byte

// State q5: Decimal state
q5:
    CMPB    AL, $'1'                 // Check if AL is less than '1'
    JB      q4_check_zero           // If less than '1', jump to q5_not_decimal

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JA      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $6, BX                   // Set the next state to q6
    
    JMP     process_next_byte        // Process the next byte

// State q6: Decimal state
q6:
    CMPB    AL, $'1'                 // Check if AL is equal to '0'
    JB      q4_check_zero           // If equal, jump to q6_not_decimal

    CMPB    AL, $'9'                 // Check if AL is greater than '9'
    JA      skip_char                // If greater than '9', jump to error

    SUBB    $'0', AL                 // Subtract '0' from AL to get the numeric value
    MOVQ    $6, BX                   // Set the next state to q6

    JMP     process_next_byte        // Process the next byte

process_next_byte:
    CMPQ    SI, R9                   // Compare current position with end of output
    JGE     error_buffer_overflow    // If we've reached or passed the end, error out
    
    MOVB    AL, (SI)                 // Store the numeric value
    INCQ    SI                       // Move the destination pointer
    INCQ    DX                       // Move the source pointer
    DECQ    CX                       // Decrement the length
    
    CMPQ    CX, $0                   // Check if the length is 0
    JE      accept                   // If it is not, continue the loop

    JMP     main_loop                // If we reach the end of the input, accept the result

skip_char:
    INCQ    DX                       // Move the source pointer
    DECQ    CX                       // Decrement the length
    JNZ     main_loop                // If CX != 0, continue the loop
    JMP     accept                   // If CX == 0, jump to accept

accept:
    CMPQ    CX, $0                   // Check if the length is 0
    JNE     error_invalid_char       // If it is not, we are done
    JMP     done

done:
    RET

error_invalid_char: // code -1 is 255 in uint8
    MOVQ    $-1, R10          // Set error code for invalid character
    JMP     save_error

error_invalid_state: // code -2 is 254 in uint8
    MOVQ    $-2, R10          // Set error code for invalid state
    JMP     save_error

error_unexpected_decimal: // code -3 is 253 in uint8
    MOVQ    $-3, R10          // Set error code for unexpected decimal point
    JMP     save_error

error_multiple_decimals: // code -4 is 252 in uint8
    MOVQ    $-4, R10          // Set error code for multiple decimal points
    JMP     save_error
    
error_buffer_overflow:
    MOVQ    $-5, R10          // Set error code for buffer overflow
    JMP     save_error

save_error:
    MOVQ    SI, AX
    SUBQ    output+24(FP), AX        // Calculate bytes written
    CMPQ    R10, $0                  // Check if there's an error
    JE      no_error
    MOVQ    R10, ret+48(FP)          // If there's an error, return the error code
    RET

no_error:
    MOVQ    AX, ret+48(FP)           // If no error, return the number of bytes written
    RET
