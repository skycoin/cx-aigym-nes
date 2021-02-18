package ui

import (
	"fmt"
	"github.com/fogleman/nes/nes"
	"image"
	"io/ioutil"
	"time"
)

const PATH_CHECKPOINTS = "../../checkpoints"
type Manager struct {
	console  *nes.Console
	hash string
}


func NewManager(console  *nes.Console) *Manager{
	manager := Manager{}
	manager.console = console
	return &manager
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
	path := fmt.Sprintf("%s/%d.ram", PATH_CHECKPOINTS, time.Now().Unix())
	err := ioutil.WriteFile(path, m.console.RAM, 0644)
	return err
}

func (m *Manager) saveScreenshot() error {
	path := fmt.Sprintf("%s/%d.png", PATH_CHECKPOINTS, time.Now().Unix())
	return savePNG(path, m.console.Buffer())
}
//func (m *Manager) saveStateToJson() error {
//
//}

func (m *Manager) hashingRomFile(filename string)  {

}


func (m *Manager) loadGame() {

}