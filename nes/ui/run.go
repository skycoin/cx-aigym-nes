package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/gordonklaus/portaudio"
	"log"
	"os"
	"runtime"
)

const (
	width  = 256
	height = 240
	scale  = 3
	title  = "NES"
)

func init() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	// we need to keep OpenGL calls on a single thread
	runtime.LockOSThread()
}

func Run(paths []string, signalChan chan os.Signal, doneChan chan int) {
	var (
		glDisabled    = true
		audioDisabled = true
		randomKeys    = false
		window        *glfw.Window
		audio         *Audio
	)

	if !audioDisabled {
		// initialize audio
		portaudio.Initialize()
		defer portaudio.Terminate()

		audio = NewAudio()
		if err := audio.Start(); err != nil {
			log.Fatalln(err)
		}
		defer audio.Stop()
	}

	if !glDisabled {
		// initialize glfw
		if err := glfw.Init(); err != nil {
			log.Fatalln(err)
		}
		defer glfw.Terminate()

		// create window
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
		var err error
		window, err = glfw.CreateWindow(width*scale, height*scale, title, nil, nil)
		if err != nil {
			log.Fatalln(err)
		}
		window.MakeContextCurrent()

		// initialize gl
		if err := gl.Init(); err != nil {
			log.Fatalln(err)
		}
		gl.Enable(gl.TEXTURE_2D)
	}

	// run director
	director := NewDirector(window, audio, signalChan, doneChan, glDisabled, audioDisabled, randomKeys)
	director.Start(paths)

}
