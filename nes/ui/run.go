package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/gordonklaus/portaudio"
	"log"
	"os"
)

const (
	width  = 256
	height = 240
	scale  = 3
	title  = "NES"
)

func Run(paths []string, signalChan chan os.Signal, disableAudio bool, disableVideo bool) {
	var (
		randomKeys = false
		window     *glfw.Window
		audio      *Audio
	)

	if !disableAudio {
		// initialize audio
		portaudio.Initialize()
		defer portaudio.Terminate()

		audio = NewAudio()
		if err := audio.Start(); err != nil {
			log.Fatalln(err)
		}
		defer audio.Stop()
	}

	if !disableVideo {

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
	director := NewDirector(window, audio, signalChan, disableVideo, disableAudio, randomKeys)
	director.Start(paths)

}
