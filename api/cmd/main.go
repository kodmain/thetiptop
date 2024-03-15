// Package main is the entry point for the application.
package main

//go:generate go run ../components/api/gen.go
//go:generate go fmt ../../api/internal/api/api.gen.go

import (
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/api"
	"github.com/kodmain/thetiptop/api/internal/lib"
	"github.com/kodmain/thetiptop/api/pkg/server"
	"github.com/spf13/cobra"
)

func main() {
	config.Helper.Run = func(cmd *cobra.Command, args []string) {
		srv := server.Create()
		srv.Register(api.Endpoints)
		srv.Start()
	}

	if err := config.Helper.Execute(); err != nil {
		lib.WithCriticalError(err)
	}
}
