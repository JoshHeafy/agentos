package out

import (
	"context"

	"github.com/pgvector/pgvector-go"

	"agentos/internal/domain"
)

type EmbeddingRepository interface {
	Search(ctx context.Context, provider domain.EmbeddingProvider, queryEmbedding pgvector.Vector, limit int) (domain.Embeddings, error)
	CreateBulk(ctx context.Context, provider domain.EmbeddingProvider, embeddings domain.Embeddings) error
}
