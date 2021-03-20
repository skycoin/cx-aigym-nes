### Summary

This score extractor will extract game infos on the 5 games located in roms.

### Usage

    go run extract.go -g=[name of the game] -f=[location of .rom file]
    go run extract.go help 

### Example
   go run extract.go -i=pacman -f=scoreextractor/checkpoints_test_data/1615018620.ram

1. If no arguments are specified, the program will return an error

3. If the integer value and file value are specified, the program will run extractor.
