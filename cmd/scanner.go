package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var value8 int8
	var value16 int16
	var value32 int32
	var value64 int64
	var inputInteger int64
	var found bool

	input := os.Getenv("integer")
	inputInteger, _ = strconv.ParseInt(input, 10, 64)

	filename := os.Getenv("file")
	data, err := ioutil.ReadFile(filename)
	check(err)

	fmt.Printf("Filename: %v\n", filename)
	for i, _ := range data {
		if (i+1)%8 == 0 {
			value64 |= int64(data[i-7])
			value64 |= int64(data[i-6]) << 8
			value64 |= int64(data[i-5]) << 16
			value64 |= int64(data[i-4]) << 24
			value64 |= int64(data[i-3]) << 32
			value64 |= int64(data[i-2]) << 40
			value64 |= int64(data[i-1]) << 48
			value64 |= int64(data[i]) << 54

			if int64(inputInteger) > 0 && value64 == int64(inputInteger) {
				fmt.Printf("Int64,%v,byte offset=%v\n", value64, i-7)
				found = true
			}
		}

		if (i+1)%4 == 0 {
			value32 |= int32(data[i-3])
			value32 |= int32(data[i-2]) << 8
			value32 |= int32(data[i-1]) << 16
			value32 |= int32(data[i]) << 24

			if int32(inputInteger) > 0 && value32 == int32(inputInteger) {
				fmt.Printf("Int32,%v,byte offset=%v\n", value32, i-3)
				found = true
			}
		}

		if (i+1)%2 == 0 {
			value16 |= int16(data[i-1])
			value16 |= int16(data[i]) << 8

			if int16(inputInteger) > 0 && value16 == int16(inputInteger) {
				fmt.Printf("Int16,%v,byte offset=%v\n", value16, i-1)
				found = true
			}
		}

		value8 = int8(data[i])
		if int8(inputInteger) > 0 && value8 == int8(inputInteger) {
			fmt.Printf("Type=Int8,%v,byte offset=%v\n", value8, i)
			found = true
		}
	}

	if !found {
		fmt.Printf("found no matches")
	}
}
