package di

import (
	"agentos/internal/application/embedding"
	"agentos/internal/domain"
	"agentos/internal/domain/ports/app"
	embeddingHandler "agentos/internal/infrastructure/inbound/httprest/handlers/embedding"
	"agentos/internal/infrastructure/outbound/gemini"
	"agentos/internal/infrastructure/outbound/ollama"
	"agentos/internal/infrastructure/outbound/openai"
	embeddingRepo "agentos/internal/infrastructure/outbound/repository/postgres/oltp/embedding"
)

func InitEmbeddingUseCase(container Container) embedding.UseCase {
	repo := embeddingRepo.NewRepository(container.DB)

	return embedding.NewUseCase(
		container.Logger,
		repo,
		map[domain.EmbeddingProvider]app.EmbeddingProvider{
			domain.Ollama: ollama.NewClient(container.Config.Ollama.Url, container.Config.Ollama.Model),
			domain.OpenAI: openai.NewClient(container.Config.Openai.Url, container.Config.Openai.ApiKey, container.Config.Openai.Model),
			domain.Gemini: gemini.NewClient(container.Config.Gemini.Url, container.Config.Gemini.ApiKey, container.Config.Gemini.Model),
		},
	)
}

func InitEmbeddingHandler(container Container) embeddingHandler.Handler {
	useCase := InitEmbeddingUseCase(container)

	return embeddingHandler.New(useCase)
}
