// Package application provides functionality for gracefully shutting down the application.
package application

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
)

var (
	PANIC chan error     = make(chan error)
	SIGS  chan os.Signal = make(chan os.Signal, 1)
)

// Wait listens for signals and errors, and performs appropriate actions based on them.
// It waits for a signal to gracefully shut down the application or for an error to occur.
func Wait() error {
	signal.Notify(SIGS, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case err := <-PANIC:
			fmt.Println(err)
			if logger.Panic(err) {
				return err
			}
		case <-SIGS:
			return nil
		}
	}
}
