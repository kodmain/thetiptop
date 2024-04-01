// Package application provides functionality for gracefully shutting down the application.
package application

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kodmain/thetiptop/api/internal/architecture/observability/logger"
)

type Exiter interface {
	Exit(code int)
}

type real struct{}

func (r *real) Exit(code int) {
	os.Exit(code)
}

var (
	PANIC chan error     = make(chan error, 1)
	SIGS  chan os.Signal = make(chan os.Signal, 1)
	PROG  Exiter         = &real{}
)

// Wait listens for signals and errors, and performs appropriate actions based on them.
// It waits for a signal to gracefully shut down the application or for an error to occur.
func Wait(exiters ...Exiter) {
	var exiter Exiter = &real{}
	if len(exiters) > 0 {
		exiter = exiters[0]
	}

	signal.Notify(SIGS, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case err := <-PANIC:
			if logger.Panic(err) {
				exiter.Exit(1)
			}
		case <-SIGS:
			logger.Info("Received signal, shutting down...")
			exiter.Exit(0)
			return
		}
	}
}
