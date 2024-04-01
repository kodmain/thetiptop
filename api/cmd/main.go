// Package main is the entry point for the application.
package main

//go:generate go run ../internal/docs/generator.go
//go:generate go fmt ../../api/internal/interfaces/api.gen.go

import (
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/kodmain/thetiptop/api/internal/architecture/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/architecture/observability/logger/levels"
	"github.com/kodmain/thetiptop/api/internal/architecture/persistence"
	"github.com/kodmain/thetiptop/api/internal/architecture/server"
	"github.com/kodmain/thetiptop/api/internal/interfaces"
	"github.com/spf13/cobra"
)

// Helper use Cobra package to create a CLI and give Args gesture
var Helper *cobra.Command = &cobra.Command{
	Use:                   "thetiptop",
	Short:                 "TheTipTop API Server",
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		logger.SetLevel(levels.DEBUG)
		if err := persistence.New(
			persistence.Config{
				Protocol: persistence.SQLite,
				Name:     "file",
				DBName:   config.DEFAULT_DB_PATH,
			},
		); err != nil {
			return err
		}

		logger.Info("loading configuration")
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("starting application")
		srv := server.Create()
		srv.Register(interfaces.Endpoints)
		srv.Start()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Info("waiting for application to shutdown")
		application.Wait()
	},
}

func main() {
	Helper.Execute()
}
