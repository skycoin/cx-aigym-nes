package main

import (
	log2 "log"
	"net/http"
	"os"
	"runtime"

	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
	"github.com/skycoin/cx-aigym-nes/cmd/rand"
	"github.com/skycoin/cx-aigym-nes/nes/ui"
	"github.com/urfave/cli/v2"
)

func init() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(runtime.NumCPU())

}

/*
cx-aigym-nes

Play game from rom files.

Usage:
  cx-aigym-nes  loadrom  --file <romfile/s> [--range <range>...] [--verbose] --random <rand> [--speed <speed>]
  cx-aigym-nes  loadjson  --file <jsonfile/s> [--range <range>...] [--verbose] --random <rand> [--speed <speed>]
  cx-aigym-nes  loadrom  --help
  cx-aigym-nes  loadjson --help
  cx-aigym-nes -h | --help
  cx-aigym-nes --version

Required options:
  -f --file <file>	    The path of .json file.

Options:
  -h --help             Shows this screen.
  --version             Shows version.

*/
var PROFILING_ENABLE bool = false

func main() {

	var (
		romPath       string
		jsonPath      string
		savedirectory string
		disableAudio  bool
		disableVideo  bool
		random        bool
		dt            float64
		speed         int
		fps           float64
	)

	app := &cli.App{
		Name:    "cx-aigym-nes",
		Version: "1.0.0",
		Flags: []cli.Flag{

			&cli.BoolFlag{
				Name:        "disable-audio",
				Usage:       "disable audio",
				Destination: &disableAudio,
			},
			&cli.BoolFlag{
				Name:        "disable-video",
				Usage:       "disable video",
				Destination: &disableVideo,
			},
			&cli.StringFlag{
				Name:        "savedirectory",
				Usage:       "Path to store the state of games",
				Destination: &savedirectory,
			},
			&cli.StringFlag{
				Name:        "loadrom",
				Value:       "",
				Aliases:     []string{"lr"},
				Usage:       "load .rom file/s",
				Destination: &romPath,
			},
			&cli.StringFlag{
				Name:        "loadjson",
				Value:       "",
				Aliases:     []string{"lj"},
				Usage:       "load .json file/s",
				Destination: &jsonPath,
			},
			&cli.BoolFlag{
				Name:        "random",
				Value:       false,
				Aliases:     []string{"r"},
				Usage:       "play random",
				Destination: &random,
			},
			&cli.Float64Flag{
				Name:        "step-seconds",
				Value:       0.016,
				Aliases:     []string{"ss"},
				Usage:       "step seconds (dt)",
				Destination: &dt,
			},
			&cli.IntFlag{
				Name:        "speed",
				Value:       1,
				Aliases:     []string{"s"},
				Usage:       "speed",
				Destination: &speed,
			},
			&cli.BoolFlag{
				Name:        "pprofile",
				Value:       false,
				Aliases:     []string{"pp"},
				Usage:       "pprofileing",
				Destination: &PROFILING_ENABLE,
			},
			&cli.Float64Flag{
				Name:        "frames-per-second",
				Value:       60,
				Aliases:     []string{"fps"},
				Usage:       "frames per second",
				Destination: &fps,
			},
		},
		Action: func(c *cli.Context) error {
			//enables profilling with pprof
			//ex. 'go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30'
			if PROFILING_ENABLE == true {
				go func() {
					log2.Println(http.ListenAndServe("localhost:6060", nil))
				}()
			}

			ui.Speed = speed
			ui.FPS = fps
			if random {
				rand.Inject()
			}
			if romPath != "" {
				return runUI(romPath, "rom", savedirectory, disableAudio, disableVideo, dt)
			} else if jsonPath != "" {
				return runUI(jsonPath, "json", savedirectory, disableAudio, disableVideo, dt)
			} else {
				log.Error("No files specified or found")
				os.Exit(1)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runUI(path, fileType string, savedirectory string,
	disableAudio bool, disableVideo bool, dt float64) error {
	if path == "" {
		log.Errorf("No %s files specified or found", fileType)
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	paths := []string{path}
	runtime.LockOSThread()
	ui.Run(paths, signalChan, savedirectory, disableAudio, disableVideo, dt)

	defer close(signalChan)
	os.Exit(0)

	return nil
}
