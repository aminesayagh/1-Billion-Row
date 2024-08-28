#include "textflag.h"

TEXT ·BytesToNumericBytes(SB), NOSPLIT, $0 
    MOVQ    len+8(FP), CX       // Load the length of the byte slice into CX
    MOVQ    B+0(FP), DX         // Load the pointer to the byte slice into DX
    MOVQ    DX, SI              // SI will be our destination pointer for valid bytes
    
loop:
    CMPQ    CX, $0              // Check if CX (length) is 0
    JLE     done                // If it is, we are done, JLE is jump if less than or equal to

    MOVB    (DX), AL           // Load the next byte from the slice into AL

    CMPB    AL, $'-'            // Compare AL to '-'
    JE      copy_char           // If they are equal, jump to copy_char

    CMPB    AL, $'+'            // Compare AL to '+'
    JE      copy_char           // If they are equal, jump to copy_char

    CMPB    AL, $'.'            // Compare AL to '.'
    JE      copy_char           // If they are equal, jump to copy_char

    CMPB    AL, $'0'            // Compare AL to '0'
    JL      skip_char           // If AL is less than '0', jump to skip

    CMPB    AL, $'9'            // Compare AL to '9'
    JG      skip_char           // If AL is greater than '9', jump to skip

    SUBB    $'0', AL            // Subtract '0' from AL to get the numeric value

copy_char:
    MOVB    AL, (SI)            // Copy the valid character to the current SI position
    INCQ    SI                  // Increment the destination pointer
    JMP     next_char           // Jump to next_char

skip_char:
    MOVB    $0, (SI)            // Copy a null byte to the current SI position
    JMP     next_char           // Jump to next_char

next_char:
    INCQ    DX                  // Increment the source pointer
    DECQ    CX                  // Decrement the loop counter
    JMP     loop                // Jump to loop

copy_null:
    MOVB    $0, (SI)            // Copy a null byte to the current SI position
    INCQ    SI                  // Increment the destination pointer
    JMP     copy_null           // Jump to copy_null

done:
fill_nulls:
    CMPQ   SI, DX              // Compare the destination pointer to the end of the slice
    JGE     finish              // If the destination pointer is greater than or equal to the end of the slice, jump to finish

    MOVB    $0, (SI)            // Copy a null byte to the current SI position
    INCQ    SI                  // Increment the destination pointer
    JMP     fill_nulls          // Jump to fill_nulls

finish:
    RET
