package embedding

import (
	"context"
	"fmt"

	"agentos/internal/domain"
	"agentos/internal/domain/ports/out"

	"github.com/pgvector/pgvector-go"
)

const table = "embeddings"

const (
	columnOllama = "embedding_ollama"
	columnOpenAI = "embedding_openai"
	columnGemini = "embedding_gemini"
)

type Repository struct {
	db out.Database
}

func NewRepository(db out.Database) Repository {
	return Repository{db: db}
}

func (r Repository) Search(ctx context.Context, provider domain.EmbeddingProvider, queryEmbedding pgvector.Vector, limit int) (domain.Embeddings, error) {
	var column string
	switch provider {
	case domain.Ollama:
		column = columnOllama
	case domain.OpenAI:
		column = columnOpenAI
	case domain.Gemini:
		column = columnGemini
	default:
		return nil, fmt.Errorf("proveedor de embedding no valido: %s", provider)
	}

	sqlQuery := fmt.Sprintf(`
		SELECT id, paragraph, 1 - (%s <=> $1) AS similarity
		FROM %s
		ORDER BY %s <=> $1
		LIMIT $2
	`, column, table, column)
	rows, err := r.db.Query(ctx, sqlQuery, queryEmbedding, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var embeddings domain.Embeddings
	for rows.Next() {
		var embedding domain.Embedding
		if err := rows.Scan(
			&embedding.ID,
			&embedding.Paragraph,
			&embedding.Similarity,
		); err != nil {
			return nil, err
		}
		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}

func (r Repository) CreateBulk(ctx context.Context, provider domain.EmbeddingProvider, embeddings domain.Embeddings) error {
	var column string

	switch provider {
	case domain.Ollama:
		column = columnOllama
	case domain.OpenAI:
		column = columnOpenAI
	case domain.Gemini:
		column = columnGemini
	default:
		return fmt.Errorf("proveedor de embedding no valido: %s", provider)
	}

	for _, embedding := range embeddings {
		query := fmt.Sprintf("INSERT INTO %s (paragraph, %s) VALUES ($1, $2)", table, column)

		if _, err := r.db.Exec(ctx, query, embedding.Paragraph, pgvector.NewVector(embedding.Embedding)); err != nil {
			return fmt.Errorf("error insertando curso: %w", err)
		}
	}

	return nil
}
