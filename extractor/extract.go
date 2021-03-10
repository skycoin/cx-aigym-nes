package main

import (
	"os"

	extractor "github.com/kenje4090/cx-aigym-nes/extractor/scoreextractor"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	var gameName string
	var filename string

	extractorApp := &cli.App{
		Name:    "Score Extractor",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "Game",
				Value:       "",
				Aliases:     []string{"g"},
				Usage:       "Name of the game",
				Destination: &gameName,
			},
			&cli.StringFlag{
				Name:        "filename",
				Value:       "",
				Aliases:     []string{"f", "file"},
				Usage:       ".rom file to scan",
				Destination: &filename,
			},
		},
		Action: func(c *cli.Context) error {
			if gameName != "" && filename != "" {
				extractor.GameExtractorFromFile(filename, gameName)
			} else {
				log.Error("No files specified or found")
				os.Exit(1)
			}
			return nil
		},
	}

	err := extractorApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
