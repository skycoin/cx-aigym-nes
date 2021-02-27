package cmd

import (
	"github.com/fogleman/nes/ui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Load the .json file",
		Run: func(cmd *cobra.Command, args []string) {

			if file == "" {
				log.Error("no rom files specified or found")
				os.Exit(1)
			}

			paths := []string{file}

			signalChan := make(chan os.Signal, 1)
			done := make(chan int)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			go func() {
				// we need to keep OpenGL calls on a single thread
				runtime.LockOSThread()
				ui.Run(paths, signalChan)
				done <- 0
			}()

			code := <-done
			os.Exit(code)
		},
	}

	file string
)

func init() {
	loadCmd.Flags().StringVarP(&file, "file", "f", "", "The path of .json file")
	loadCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(loadCmd)
}
