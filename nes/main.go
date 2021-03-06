package main

import (
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
	_ "github.com/skycoin/cx-aigym-nes/nes/cmd"
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
		Flags: []cli.Flag{
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
		},
		Action: func(c *cli.Context) error {
			if romPath != "" {
				return runUI(romPath, "rom")
			} else if jsonPath != "" {
				return runUI(jsonPath, "json")
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

func runUI(path, fileType string) error {
	if path == "" {
		log.Errorf("No %s files specified or found", fileType)
		os.Exit(1)
	}

	if true {
		cmd.Inject()
	}

	signalChan := make(chan os.Signal, 1)
	paths := []string{path}
	runtime.LockOSThread()
	ui.Run(paths, signalChan)

	defer close(signalChan)
	os.Exit(0)

	return nil
}
