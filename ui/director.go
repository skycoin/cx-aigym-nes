package ui

import (
	"github.com/fogleman/nes/nes"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
	"time"
)

type View interface {
	Enter()
	Exit()
	Update(t, dt float64)
}

type Director struct {
	window    *glfw.Window
	audio     *Audio
	view      View
	menuView  View
	timestamp float64
	GlDisabled bool
	AudioDisabled bool
}

func NewDirector(window *glfw.Window, audio *Audio, glDisabled bool, audioDisabled bool) *Director {
	director := Director{}
	director.window = window
	director.audio = audio
	director.GlDisabled = glDisabled
	director.AudioDisabled = audioDisabled
	return &director
}

func (d *Director) SetGlDisabled(glDisabled bool) {
	d.GlDisabled = glDisabled
}

func (d *Director) SetAudioDisabled(audioDisabled bool) {
	d.AudioDisabled = audioDisabled
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

	if d.GlDisabled {
		d.timestamp = float64(time.Now().UnixNano())
	} else {
		d.timestamp = glfw.GetTime()
	}
}

func (d *Director) Step() {
	var timestamp float64
	if !d.GlDisabled {
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
	if d.GlDisabled {
		for  {
			d.Step()
		}
	} else {
		for !d.window.ShouldClose() {
			d.Step()
			d.window.SwapBuffers()
			glfw.PollEvents()
		}
	}

	d.SetView(nil)
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
