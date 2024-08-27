section .data
    num1 dq 5 ; 64-bit floating point number 5.0
    num2 dq 3 ; 64-bit floating point number 3.0
    result dq 0 ; 64-bit floating point number 0.0

section .text
    global _start ; Entry point for the program

_start:
    mov rax, [num1] ; Load num1 into rax register, the [num1] is the memory address of num1 variable
    mov rdx, [num2] ; Load num2 into rdx register
    mov [result], rax ; Store the value of rax register into result variable

    add rax, rdx ; Add the value of rdx register to rax register
    mov [result], rax ; Store the result of addition into result variable

    ; Exit the program
    mov rax, 60 ; syscall number for exit
    xor rdi, rdi ; Exit code 0
    syscall