package app

import (
	"agentos/internal/domain"
	"agentos/internal/infrastructure/outbound/ollama"
	"context"
	"mime/multipart"
)

type EmbeddingUseCase interface {
	Search(ctx context.Context, embeddingProvider domain.EmbeddingProvider, query string, limit int) (map[string]any, error)
	Upload(ctx context.Context, embeddingProvider domain.EmbeddingProvider, file *multipart.FileHeader) error
}

type EmbeddingProvider interface {
	GenerateEmbeddings(chunks []string) (domain.Embeddings, error)
	Chat(messages ollama.Messages, model domain.OllamaModel, options *ollama.Options) (string, error)
}
