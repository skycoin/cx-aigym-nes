package ui

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	log "github.com/sirupsen/logrus"
	extractor "github.com/skycoin/cx-aigym-nes/cmd/scoreextractor"
	"github.com/skycoin/cx-aigym-nes/nes/nes"
)

const padding = 0
const PATH_CHECKPOINTS = "../checkpoints"

type KeyReader func(window *glfw.Window, turbo bool) [8]bool

var ReadKeys KeyReader

// Speed - game speed
var Speed int

// CyclePerMS - number of cycles per milliseconds in current machine
var CyclePerMS float64
var FPS float64

func init() {
	ReadKeys = readKeys
}

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
		record:     false,
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
	// cartridge := view.console.Cartridge
	// if cartridge.Battery != 0 {
	// 	writeSRAM(sramPath(view.RomHash), cartridge.SRAM)
	// }
	// // save state
	// view.console.SaveState(savePath(view.RomHash))

}

// GetDt returns the step seconds according to the specified speed
func (view *GameView) GetDt(dt float64) float64 {
	cycles := int(nes.CPUFrequency * dt)
	if CyclePerMS == 0 {
		return dt
	}

	maxCycles := int((1000 / FPS) * CyclePerMS)
	newCycles := Speed * cycles
	if Speed == 0 {
		return float64(maxCycles) / float64(nes.CPUFrequency)
	}
	if newCycles < maxCycles {
		return float64(Speed) * dt
	}
	return float64(maxCycles) / float64(nes.CPUFrequency)
}

func (view *GameView) Update(t, dt float64) {
	if dt > 1 {
		dt = 0
	}

	dt = view.GetDt(dt)
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

	t1 := time.Now()
	cycles := console.StepSeconds(dt)
	tickTime := time.Now().Sub(t1)
	CyclePerMS = float64(cycles) / float64(tickTime.Milliseconds())

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

func (view *GameView) saveStateToFiles(saveDirectory string, now int64) error {
	filepath := fmt.Sprintf("%s/%d.ram",
		strings.TrimSuffix(saveDirectory, "/"), now)
	return view.console.SaveState(filepath)
}

func (view *GameView) saveToJson(saveDirectory string, now int64) error {

	view.StateHash = fmt.Sprintf("%x", sha256.Sum256(view.state))
	data, err := json.MarshalIndent(view, "", " ")
	if err != nil {
		return err
	}

	filepath := fmt.Sprintf("%s/%d.json",
		strings.TrimSuffix(saveDirectory, "/"), now)
	return ioutil.WriteFile(filepath, data, 0644)
}

func (view *GameView) save() error {
	saveDirectory := view.director.savedirectory
	if saveDirectory == "" {
		saveDirectory = PATH_CHECKPOINTS
	}

	log.Infof("save state of game to %s", saveDirectory)
	view.Timestamp = time.Now().Unix()
	view.saveStateToFiles(saveDirectory, view.Timestamp)
	view.saveScreenshot(saveDirectory, view.Timestamp)
	view.saveToJson(saveDirectory, view.Timestamp)
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

		case glfw.Key6:
			// print out extractor
			extractor.ExtractGameDetails(view.RomHash, view.console.RAM)
		}

	}
}

func (view *GameView) saveScreenshot(saveDirectory string, now int64) error {
	filepath := fmt.Sprintf("%s/%d.png",
		strings.TrimSuffix(saveDirectory, "/"), now)
	return savePNG(filepath, view.console.Buffer())
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

	k1 = ReadKeys(director.window, turbo)

	// if director.glDisabled || director.randomKeys {
	// 	k1 = readRandomKeys()
	// } else {
	// 	k1 = readKeys(director.window, turbo)
	// }

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
