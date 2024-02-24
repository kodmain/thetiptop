// Package main is the entry point for the application.
package main

import (
	"github.com/kodmain/thetiptop/project/config"
	v1 "github.com/kodmain/thetiptop/project/internal/api/v1"
	"github.com/kodmain/thetiptop/project/internal/lib"
	"github.com/kodmain/thetiptop/project/pkg/server"
	"github.com/spf13/cobra"
)

func main() {
	config.Helper.Run = func(cmd *cobra.Command, args []string) {
		srv := server.Create()
		srv.API(v1.Status)
		srv.API(v1.Fizzbuzz)
		srv.Start()
	}

	if err := config.Helper.Execute(); err != nil {
		lib.WithCriticalError(err)
	}
}
