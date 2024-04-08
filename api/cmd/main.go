// Package main is the entry point for the application.
package main

//go:generate go run ../internal/docs/generator.go
//go:generate go fmt ../internal/interfaces/api.gen.go

import (
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger/levels"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
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
		logger.SetLevel(levels.DEBUG)
		cfg, err := config.Load("config.yml")

		if err != nil {
			return err
		}

		if err := database.New(cfg.Databases); err != nil {
			return err
		}

		if err := mail.New(cfg.Mail); err != nil {
			return err
		}

		/*
			err = mail.Send(&mail.Mail{
				To:      []string{"extazy937@gmail.com"},
				Cc:      []string{"alt.zo-8of03jkz@yopmail.com"},
				Subject: "TheTipTop API Server",
				Text:    []byte("TheTipTop API Server is running"),
				Html:    []byte("<h1>TheTipTop API Server is running</h1>"),
			})
		*/

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
