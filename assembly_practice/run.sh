#!/bin/bash

# Name of the assembly source file (without extension)
SOURCE_FILE=$1

# Assemble the source file
nasm -f elf64 ${SOURCE_FILE}.asm -o ${SOURCE_FILE}.o
if [ $? -ne 0 ]; then
    echo "Assembly failed."
    exit 1
fi

# Link the object file to create an executable
ld ${SOURCE_FILE}.o -o ${SOURCE_FILE}
if [ $? -ne 0 ]; then
    echo "Linking failed."
    exit 1
fi

# Run the executable
./${SOURCE_FILE}
if [ $? -ne 0 ]; then
    echo "Execution failed."
    exit 1
fi

echo "Program executed successfully."
