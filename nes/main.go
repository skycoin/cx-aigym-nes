package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	log "github.com/sirupsen/logrus"
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
  cx-aigym-nes  loadrom  --file <romfile/s> [--range <range>...] [--verbose]
  cx-aigym-nes  loadjson  --file <jsonfile/s> [--range <range>...] [--verbose]
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
func main() {
	var romPath string
	var jsonPath string

	app := &cli.App{
		Name:    "cx-aigym-nes",
		Version: "1.0.0",
		Commands: []*cli.Command{
			{
				Name:    "loadrom",
				Aliases: []string{"lr"},
				Usage:   "load rom file/s",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "file",
						Value:       "",
						Aliases:     []string{"f"},
						Usage:       "load .rom file/s",
						Destination: &romPath,
						Required:    true,
					},
				},
				Action: func(c *cli.Context) error {
					return runUI(romPath, "rom")
				},
			},
			{
				Name:    "loadjson",
				Aliases: []string{"lj"},
				Usage:   "load json file/s",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "file",
						Value:       "",
						Aliases:     []string{"f"},
						Usage:       "load .json file/s",
						Destination: &jsonPath,
						Required:    true,
					},
				},
				Action: func(c *cli.Context) error {
					return runUI(jsonPath, "json")
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runUI(path, fileType string) error {
	signalChan := make(chan os.Signal, 1)
	done := make(chan int)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	if path == "" {
		log.Errorf("No %s files specified or found", fileType)
		os.Exit(1)

	}
	paths := []string{path}
	runtime.LockOSThread()
	ui.Run(paths, signalChan)
	code := <-done
	os.Exit(code)
	return nil
}
