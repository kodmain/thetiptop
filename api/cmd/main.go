// Package main is the entry point for the application.
package main

//go:generate go run ../internal/docs/generator.go
//go:generate go fmt ../../api/internal/api/api.gen.go

import (
	"github.com/kodmain/thetiptop/api/internal/api"
	"github.com/kodmain/thetiptop/api/internal/application/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/application/observability/logger/levels"
	"github.com/kodmain/thetiptop/api/internal/architecture/kernel"
	"github.com/kodmain/thetiptop/api/internal/architecture/server"
	"github.com/spf13/cobra"
)

// Helper use Cobra package to create a CLI and give Args gesture
var Helper *cobra.Command = &cobra.Command{
	Use:                   "fizzbuzz",
	Short:                 "Fizzbuzz API Server",
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.SetLevel(levels.DEBUG)
		logger.Info("loading configuration")
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("starting application")
		srv := server.Create()
		srv.Register(api.Endpoints)
		srv.Start()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Info("waiting for application to shutdown")
		kernel.Wait()
	},
}

func main() {
	Helper.Execute()
}
