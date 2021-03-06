package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	maxInt8  = 127
	minInt8  = -128
	maxInt16 = 32767
	minInt16 = -32768
	maxInt32 = 2147483647
	minInt32 = -2147483648
	maxInt64 = 9223372036854775807
	minInt64 = -9223372036854775808
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var inputString string
	var inputInteger int64
	var filename string

	scannerApp := &cli.App{
		Name:    "Scanner",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "integer",
				Value:       "",
				Aliases:     []string{"i"},
				Usage:       "Value to seek from .rom file",
				Destination: &inputString,
			},
			&cli.StringFlag{
				Name:        "filename",
				Value:       "",
				Aliases:     []string{"f", "file"},
				Usage:       ".rom file to scan",
				Destination: &filename,
			},
		},
		Action: func(c *cli.Context) error {
			if inputString != "" && filename != "" {
				inputInteger, _ = strconv.ParseInt(inputString, 10, 64)
				scanner(filename, inputInteger)
			} else {
				log.Error("No files specified or found")
				os.Exit(1)
			}
			return nil
		},
	}

	err := scannerApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func scanner(filename string, inputInteger int64) {
	var value8 int8
	var value16 int16
	var value32 int32
	var value64 int64
	var found bool

	data, err := ioutil.ReadFile(filename)
	check(err)

	fmt.Printf("Filename: %v\n", filename)
	for i := range data {
		// Finding match for int64 values
		if inputInteger <= maxInt64 && inputInteger >= minInt64 {
			if i >= 7 {
				value64 |= int64(data[i-7])
				value64 |= int64(data[i-6]) << 8
				value64 |= int64(data[i-5]) << 16
				value64 |= int64(data[i-4]) << 24
				value64 |= int64(data[i-3]) << 32
				value64 |= int64(data[i-2]) << 40
				value64 |= int64(data[i-1]) << 48
				value64 |= int64(data[i]) << 54

				if value64 == int64(inputInteger) {
					fmt.Printf("Int64,%v,byte offset=%v\n", value64, i-7)
					found = true
				}
			}
		}

		// Finding match for int32 values
		if inputInteger <= maxInt32 && inputInteger >= minInt32 {
			if i >= 3 {
				value32 |= int32(data[i-3])
				value32 |= int32(data[i-2]) << 8
				value32 |= int32(data[i-1]) << 16
				value32 |= int32(data[i]) << 24

				if value32 == int32(inputInteger) {
					fmt.Printf("Int32,%v,byte offset=%v\n", value32, i-3)
					found = true
				}
			}
		}

		// Finding match for int16 values
		if inputInteger <= maxInt16 && inputInteger >= minInt16 {
			if i >= 1 {
				value16 |= int16(data[i-1])
				value16 |= int16(data[i]) << 8

				if value16 == int16(inputInteger) {
					fmt.Printf("Int16,%v,byte offset=%v\n", value16, i-1)
					found = true
				}
			}
		}

		// Finding match for int8 values
		if inputInteger <= maxInt8 && inputInteger >= minInt8 {
			value8 = int8(data[i])
			if value8 == int8(inputInteger) {
				fmt.Printf("Int8,%v,byte offset=%v\n", value8, i)
				found = true
			}
		}

	}

	if !found {
		fmt.Printf("found no matches")
	}
}
