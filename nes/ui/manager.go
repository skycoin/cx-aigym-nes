package ui

import (
	"fmt"
	"github.com/fogleman/nes/nes"
	"image"
	"io/ioutil"
	"time"
)

const PATH_STATE = "./tmp/state"

type Manager struct {
	console  *nes.Console

}


func NewManager(console  *nes.Console) *Manager{
	manager := Manager{}
	manager.console = console
	return &Manager{}
}

func (m *Manager) captureImageFrame() *image.RGBA {
	return m.console.PPU.Front
}

func (m *Manager) getRam() []byte {
	return m.console.RAM
}

func (m *Manager) setRam(state []byte) {
		m.console.RAM = state
}

func (m *Manager) saveStateToFile() error {
	path := fmt.Sprintf("%s/%d.ram", PATH_STATE, time.Now().Unix())
	err := ioutil.WriteFile(path, m.console.RAM, 0644)
	return err
}

//func (m *Manager) saveStateToJson() error {
//
//}
