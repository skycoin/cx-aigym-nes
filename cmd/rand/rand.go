package rand

import (
	"log"
	"math/rand"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/skycoin/cx-aigym-nes/nes/nes"
	"github.com/skycoin/cx-aigym-nes/nes/ui"
)

func randomKeys() int {
	return rand.Intn(64)
}

func convertIntToBool(v int) bool {
	if v != 0 {
		return true
	} else {
		return false
	}
}

// ReadRandomKeys generates random keys
func ReadRandomKeys(window *glfw.Window, turbo bool) [8]bool {
	var result [8]bool
	keys := randomKeys()
	log.Printf("%b", keys)

	result[nes.ButtonA] = convertIntToBool(keys & 1)
	result[nes.ButtonB] = convertIntToBool(keys & 2)
	//result[nes.ButtonSelect] = convertIntToBool(keys & 4)
	//result[nes.ButtonStart] = convertIntToBool(keys & 8)
	result[nes.ButtonUp] = convertIntToBool(keys & 4)
	result[nes.ButtonDown] = convertIntToBool(keys & 8)
	result[nes.ButtonLeft] = convertIntToBool(keys & 16)
	result[nes.ButtonRight] = convertIntToBool(keys & 32)
	return result
}

// Inject injects ReadRandomKeys to gameview
func Inject() {
	ui.ReadKeys = ReadRandomKeys
}
