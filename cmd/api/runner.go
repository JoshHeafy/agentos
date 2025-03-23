package main

import (
	"context"
	"os"

	"agentos/cmd"
	"agentos/pkg/logger"
)

func Run() {
	logger := logger.NewZeroLog(ServiceName)
	server, err := cmd.NewServerInstance(os.Getenv("CONFIGURATION_FILEPATH"), ServiceName)
	if err != nil {
		logger.Error(context.Background(), err.Error())
		return
	}

	SetAPIRoutes(server)

	if err := server.Start(); err != nil {
		logger.Error(context.Background(), err.Error())
	}
}
