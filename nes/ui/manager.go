package ui

import (
	"github.com/fogleman/nes/nes"
	"image"
)

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
