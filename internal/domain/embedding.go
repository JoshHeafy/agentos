package domain

import "github.com/pgvector/pgvector-go"

type EmbeddingProvider string

const (
	OpenAI EmbeddingProvider = "openai"
	Ollama EmbeddingProvider = "ollama"
	Gemini EmbeddingProvider = "gemini"
)

type Embedding struct {
	ID              uint            `json:"id"`
	Paragraph       string          `json:"paragraph"`
	Embedding       []float32       `json:"-"`
	EmbeddingVector pgvector.Vector `json:"embedding"`
	Similarity      float64         `json:"similarity"`
}

type Embeddings []Embedding
