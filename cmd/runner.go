package cmd

import (
	"agentos/pkg/db"
	"agentos/pkg/di"
	"agentos/pkg/localconfig"
	"agentos/pkg/logger"
	"agentos/pkg/server"
)

// NewServerInstance is a function that receives a configuration path and a service name and returns a server instance.
func NewServerInstance(configPath, serviceName string) (*server.Server, error) {
	config, err := localconfig.NewLocalConfig(configPath)
	if err != nil {
		return nil, err
	}

	loggerUseCase := logger.NewZeroLog(serviceName)

	dbUseCase, err := db.New(config.Database)
	if err != nil {
		return nil, err
	}

	serverUseCase := server.New(di.Container{
		Logger: loggerUseCase,
		Config: config,
		DB:     dbUseCase,
	}, serviceName)

	return serverUseCase, nil
}
