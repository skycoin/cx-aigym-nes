### Summary

This scanner will scan a file and traverse through the file byte by byte and convert their locations to int8,int16,int32,int64 at each position, if the decoded value matches the input valie then it should output the type, value, and offset location in bytes

Note: All are in little endian.

### Usage

    go run scanner.go -i=[value to match] -f=[location of .rom file]
    go run scanner.go help 

### Example
   go run scanner.go -i=50 -f=../checkpoints/1614939776.ram  

1. If no arguments are specified, the program will look return an error

3. If the integer value and file value are specified, the program will run scanner.
