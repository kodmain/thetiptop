// Package main is the entry point for the application.
package main

//go:generate go run ../internal/docs/generator.go
//go:generate go fmt ../internal/interfaces/api.gen.go

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/kodmain/thetiptop/api/internal/application/hook"
	"github.com/kodmain/thetiptop/api/internal/docs/generated"
	"github.com/kodmain/thetiptop/api/internal/domain/game/events"
	"github.com/kodmain/thetiptop/api/internal/domain/game/repositories"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger/levels"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
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
		hook.Call(hook.EventOnConfig)
		generated.SwaggerInfo.Version = env.BUILD_VERSION
		logger.SetLevel(levels.DEBUG)
		hook.Register(hook.EventOnDBInit, func() {
			events.HydrateDBWithTickets(
				repositories.NewGameRepository(database.Get(config.GetString("services.game.database", config.DEFAULT))),
				config.Get("project.tickets.required", 10000).(int),
				config.Get("project.tickets.types", map[string]int{}).(map[string]int),
			)
		})

		return config.Load(env.CONFIG_URI)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("starting application")
		hook.Call(hook.EventOnStart)
		srv := server.Create()
		srv.Register(interfaces.Endpoints)
		return srv.Start()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("waiting for application to shutdown")
		return application.Wait()
	},
}

// Version de l'application
var Version = generated.SwaggerInfo.Version

// versionCmd représente la commande de version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	Long:  "show version of the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version %s \n", env.BUILD_VERSION)
		fmt.Printf("Commit %s \n", env.BUILD_COMMIT)
	},
}

// @title		TheTipTop
// @version		dev
// @description	TheTipTop API
// @host		localhost
// @BasePath
func main() {
	env.CONFIG_URI = Helper.Flags().String("config", env.DEFAULT_CONFIG_URI, "URI de la configuration")
	env.AWS_PROFILE = Helper.Flags().String("profile", env.DEFAULT_AWS_PROFILE, "Profil AWS")
	env.PORT_HTTP = Helper.Flags().Int("http-port", env.DEFAULT_PORT_HTTP, "Port HTTP")
	env.PORT_HTTPS = Helper.Flags().Int("https-port", env.DEFAULT_PORT_HTTPS, "Port HTTPS")

	Helper.AddCommand(versionCmd)
	Helper.Execute()
}
