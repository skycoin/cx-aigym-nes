package ui

import (
	"github.com/fogleman/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"os"
	"time"
)

type View interface {
	Enter()
	Exit()
	Update(t, dt float64)
}

type Director struct {
	window        *glfw.Window
	audio         *Audio
	view          View
	menuView      View
	timestamp     float64
	glDisabled    bool
	audioDisabled bool
	randomKeys    bool
	doneChan      chan int
	signalChan    chan os.Signal
}

func NewDirector(window *glfw.Window, audio *Audio,
	signalChan chan os.Signal, doneChan chan int, glDisabled bool,
	audioDisabled bool, randomKeys bool) *Director {
	director := Director{}
	director.window = window
	director.audio = audio
	director.glDisabled = glDisabled
	director.audioDisabled = audioDisabled
	director.randomKeys = randomKeys
	director.doneChan = doneChan
	director.signalChan = signalChan
	return &director
}

func (d *Director) SetGlDisabled(glDisabled bool) {
	d.glDisabled = glDisabled
}

func (d *Director) SetAudioDisabled(audioDisabled bool) {
	d.audioDisabled = audioDisabled
}

func (d *Director) setRandomKeys(randomKeys bool) {
	d.randomKeys = randomKeys
}

func (d *Director) SetTitle(title string) {
	d.window.SetTitle(title)
}

func (d *Director) SetView(view View) {
	if d.view != nil {
		d.view.Exit()
	}
	d.view = view
	if d.view != nil {
		d.view.Enter()
	}

	if d.glDisabled {
		d.timestamp = float64(time.Now().UnixNano())
	} else {
		d.timestamp = glfw.GetTime()
	}
}

func (d *Director) Step() {
	var timestamp float64
	if !d.glDisabled {
		gl.Clear(gl.COLOR_BUFFER_BIT)
		timestamp = glfw.GetTime()

	} else {
		timestamp = float64(time.Now().UnixNano())
	}

	dt := timestamp - d.timestamp
	d.timestamp = timestamp
	if d.view != nil {
		d.view.Update(timestamp, dt)
	}
}

func (d *Director) Start(paths []string) {
	d.menuView = NewMenuView(d, paths)
	if len(paths) == 1 {
		d.PlayGame(paths[0])
	} else {
		d.ShowMenu()
	}
	d.Run()
}

func (d *Director) Run() {
	if d.glDisabled {
		for {
			select {
			case <-d.signalChan:
				d.SetView(nil)
				d.doneChan <- 0
				return
			default:
				d.Step()
				time.Sleep(10 * time.Millisecond)
			}
		}
	} else {
		for !d.window.ShouldClose() {
			d.Step()
			d.window.SwapBuffers()
			glfw.PollEvents()
		}

		d.SetView(nil)
	}

}

func (d *Director) PlayGame(path string) {
	hash, err := hashFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	console, err := nes.NewConsole(path)
	if err != nil {
		log.Fatalln(err)
	}

	d.SetView(NewGameView(d, console, path, hash))
}

func (d *Director) ShowMenu() {
	d.SetView(d.menuView)
}
