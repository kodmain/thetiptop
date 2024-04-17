// Package main is the entry point for the application.
package main

//go:generate go run ../internal/docs/generator.go
//go:generate go fmt ../internal/interfaces/api.gen.go

import (
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/kodmain/thetiptop/api/internal/docs/generated"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger/levels"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/server"
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
		logger.Info("loading configuration")
		generated.SwaggerInfo.Version = config.BUILD_VERSION
		logger.SetLevel(levels.DEBUG)
		// cfg, err := config.Load("config.yml")
		return config.Load(config.DEFAULT_CONFIG)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("starting application")
		srv := server.Create()
		srv.Register(interfaces.Endpoints)
		return srv.Start()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("waiting for application to shutdown")
		return application.Wait()
	},
}

// @title		TheTipTop
// @version		1.0 // BUILD_VERSION
// @description	TheTipTop API
// @host		localhost
// @BasePath
func main() {
	Helper.Execute()
}

// test bug sonar.
func tydoNothing() {}
