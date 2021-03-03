package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/skycoin/cx-aigym-nes/nes/cmd"
)

func init() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(runtime.NumCPU())

}
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
