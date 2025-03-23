package main

import (
	"agentos/pkg/di"
	"agentos/pkg/server"
)

func SetAPIRoutes(server *server.Server) {
	server.Echo.GET("/health", server.HealthCheckController)

	embeddingRoutes(server)
}

func embeddingRoutes(server *server.Server) {
	campaignHandler := di.InitEmbeddingHandler(server.Container)

	publicApi := server.PublicAPI.Group("/embedding")
	publicApi.POST("/upload", campaignHandler.Upload)
	publicApi.GET("/search", campaignHandler.Search)
}
