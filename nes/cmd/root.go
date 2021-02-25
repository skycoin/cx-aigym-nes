package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

/*
cx-aigym-nes

Play game from rom files.

Usage:
  cx-aigym-nes  load  --file <file> [--range <range>...] [--verbose]
  cx-aigym-nes -h | --help
  cx-aigym-nes --version

Required options:
  -f --file <file>	    The path of .json file.

Options:
  -h --help             Shows this screen.
  --version             Shows version.

*/

var (
	rootCmd = &cobra.Command{
		Use:     "cx-aigym-nes",
		Version: "1.0.0",
		Short:   "Play game from rom file",
	}

	levelLogging string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	configLogging()
	rootCmd.PersistentFlags().StringVarP(&levelLogging, "level", "l",
		"info", "level logging")

}

func configLogging() {

	// Set the logging level
	switch levelLogging {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// Set the TextFormatter
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
	})

	log.Infoln("cx-aigym-nes is starting")
}
