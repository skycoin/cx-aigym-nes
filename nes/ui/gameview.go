package ui

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	log "github.com/sirupsen/logrus"
	"github.com/skycoin/cx-aigym-nes/nes/nes"
)

const padding = 0
const PATH_CHECKPOINTS = "../checkpoints"

var currentGameView *GameView

type GameView struct {
	director   *Director
	state      []byte
	StateHash  string `json:"state_hash"`
	console    *nes.Console
	RomPath    string `json:"rom_path"`
	RomName    string `json:"rom_name"`
	RomHash    string `json:"rom_hash"`
	texture    uint32
	record     bool
	FrameIndex int `json:"frame_index"`
	frames     []image.Image
	Timestamp  int64 `json:"timestamp"`
}

func NewGameView(director *Director, console *nes.Console, path string, hash string) View {
	var texture uint32
	if !director.glDisabled {
		texture = createTexture()
	}

	name := filepath.Base(path)
	currentGameView = &GameView{
		director:   director,
		console:    console,
		RomName:    name,
		RomPath:    path,
		RomHash:    hash,
		FrameIndex: 0,
		texture:    texture,
		record:     true,
		frames:     nil,
	}

	return currentGameView
}

func (view *GameView) Enter() {
	if !view.director.glDisabled {
		gl.ClearColor(0, 0, 0, 1)
		view.director.SetTitle(view.RomName)
		view.director.window.SetKeyCallback(view.onKey)
	}

	if !view.director.audioDisabled {
		view.console.SetAudioChannel(view.director.audio.channel)
		view.console.SetAudioSampleRate(view.director.audio.sampleRate)
	}

	// load state
	if err := view.console.LoadState(savePath(view.RomHash)); err == nil {
		return
	} else {
		view.console.Reset()
	}
	// load sram
	cartridge := view.console.Cartridge
	if cartridge.Battery != 0 {
		if sram, err := readSRAM(sramPath(view.RomHash)); err == nil {
			cartridge.SRAM = sram
		}
	}
}

func (view *GameView) Exit() {
	if !view.director.glDisabled {
		view.director.window.SetKeyCallback(nil)
		view.console.SetAudioChannel(nil)
		view.console.SetAudioSampleRate(0)
	}

	// save sram
	cartridge := view.console.Cartridge
	if cartridge.Battery != 0 {
		writeSRAM(sramPath(view.RomHash), cartridge.SRAM)
	}
	// save state
	view.console.SaveState(savePath(view.RomHash))

}

func (view *GameView) Update(t, dt float64) {
	if dt > 1 {
		dt = 0
	}
	window := view.director.window
	console := view.console
	if !view.director.glDisabled {
		if joystickReset(glfw.Joystick1) {
			view.director.ShowMenu()
		}
		if joystickReset(glfw.Joystick2) {
			view.director.ShowMenu()
		}
		if readKey(window, glfw.KeyEscape) {
			view.director.ShowMenu()
		}

	}

	updateControllers(view.director, console)
	console.StepSeconds(dt)

	if !view.director.glDisabled {
		gl.BindTexture(gl.TEXTURE_2D, view.texture)
		setTexture(console.Buffer())
		drawBuffer(view.director.window)
		gl.BindTexture(gl.TEXTURE_2D, 0)
	}

	if view.record {
		view.frames = append(view.frames, copyImage(console.Buffer()))
	}

	// count frame
	view.FrameIndex++
}

func (view *GameView) saveState() {
	log.Infof("Save state to buffer bytes: %v", view.console)
	view.state = view.console.SaveStateToBytes()
}

func (view *GameView) loadState() {
	log.Infof("Load state from checkpoint: %v", view.state)
	if view.state != nil {
		view.console.LoadStateFromBytes(view.state)
	}
}

func (view *GameView) saveStateToFiles(now int64) error {
	path := fmt.Sprintf("%s/%d.ram", PATH_CHECKPOINTS, now)
	return view.console.SaveState(path)
}

func (view *GameView) saveToJson(now int64) error {

	view.StateHash = fmt.Sprintf("%x", sha256.Sum256(view.state))
	file, err := json.MarshalIndent(view, "", " ")
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%d.json", PATH_CHECKPOINTS, now)
	return ioutil.WriteFile(path, file, 0644)
}

func (view *GameView) save() error {
	log.Infof("save state of game to %s: %v", PATH_CHECKPOINTS, view)
	now := time.Now().Unix()
	view.Timestamp = now
	view.saveStateToFiles(now)
	view.saveScreenshot(now)
	view.saveToJson(now)
	return nil
}

func (view *GameView) onKey(window *glfw.Window,
	key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeySpace:
			screenshot(view.console.Buffer())
		case glfw.KeyR:
			view.console.Reset()
		case glfw.KeyTab:
			if view.record {
				view.record = false
				animation(view.frames)
				view.frames = nil
			} else {
				view.record = true
			}
		case glfw.Key1:
			// save state to bytes
			view.saveState()

		case glfw.Key2:
			// load state from bytes
			view.loadState()

		case glfw.Key5:
			// save state to file
			view.save()
		}

	}
}

func (view *GameView) saveScreenshot(now int64) error {
	path := fmt.Sprintf("%s/%d.png", PATH_CHECKPOINTS, now)
	return savePNG(path, view.console.Buffer())
}

func (view *GameView) captureImageFrame() *image.RGBA {
	return view.console.PPU.Front
}

func drawBuffer(window *glfw.Window) {
	w, h := window.GetFramebufferSize()
	s1 := float32(w) / 256
	s2 := float32(h) / 240
	f := float32(1 - padding)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()
}

func updateControllers(director *Director, console *nes.Console) {
	turbo := console.PPU.Frame%6 < 3

	var j1, j2, k1 [8]bool

	if director.glDisabled || director.randomKeys {
		k1 = readRandomKeys()
	} else {
		k1 = readKeys(director.window, turbo)
	}

	if !director.glDisabled {
		j1 = readJoystick(glfw.Joystick1, turbo)
		j2 = readJoystick(glfw.Joystick2, turbo)
	}

	console.SetButtons1(combineButtons(k1, j1))
	console.SetButtons2(j2)
}

//Implement "GetRomFilename"" And "GetRomHash()" #22
func GetRomFilename() string {
	if currentGameView == nil || currentGameView.RomName == "" {
		log.Panic("no rom is loaded")
	}

	return currentGameView.RomName
}

func GetRomHash() string {
	if currentGameView == nil || currentGameView.RomName == "" {
		log.Panic("no rom is loaded")
	}

	return currentGameView.RomHash
}
