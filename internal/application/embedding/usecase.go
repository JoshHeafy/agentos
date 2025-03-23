package embedding

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"agentos/internal/domain"
	"agentos/internal/domain/ports/app"
	"agentos/internal/domain/ports/out"
	"agentos/internal/infrastructure/outbound/ollama"
	"agentos/pkg/tokenizer"

	"github.com/pgvector/pgvector-go"
)

var systemPrompt = `
	Responde la pregunta del usuario basándote en las siguientes restricciones:
	- Debes responder en español.
	- No utilices tu conocimiento previo, debes responder ÚNICAMENTE con el contexto que se te envía en el mensaje o con el contexto.
	- No alucines ni inventes información adicional.
`

const (
	defaultChunkSize    = 512
	defaultChunkOverlap = 50
	defaultChunkModel   = "deepseek-r1"
)

type UseCase struct {
	logger   out.Logger
	repo     out.EmbeddingRepository
	provider map[domain.EmbeddingProvider]app.EmbeddingProvider
}

func NewUseCase(
	logger out.Logger,
	repo out.EmbeddingRepository,
	provider map[domain.EmbeddingProvider]app.EmbeddingProvider,
) UseCase {
	return UseCase{
		logger:   logger,
		repo:     repo,
		provider: provider,
	}
}

func (uc UseCase) Search(ctx context.Context, embeddingProvider domain.EmbeddingProvider, query string, limit int) (map[string]any, error) {
	var err error

	messagesAssist := ollama.Messages{{Role: "user", Content: systemPrompt}}

	// Improve question
	improvedQuestion, err := uc.improveQuestion(query)
	if err != nil {
		return nil, err
	}

	// Embed the question
	embeddedQuestion, err := uc.embedQuestion(embeddingProvider, improvedQuestion)
	if err != nil {
		return nil, err
	}

	contexts, err := uc.searchEmbedInDB(ctx, embeddingProvider, embeddedQuestion, limit)
	if err != nil {
		return nil, err
	}

	// Ask for the answer from IA
	response, err := uc.askToIA(messagesAssist, contexts, improvedQuestion)
	if err != nil {
		return nil, err
	}

	result := map[string]any{
		"provider":   embeddingProvider,
		"created_at": time.Now(),
		"question":   query,
		"message": map[string]any{
			"role":    "assistant",
			"content": response,
		},
	}

	return result, nil
}

func (uc UseCase) Upload(ctx context.Context, embeddingProvider domain.EmbeddingProvider, file *multipart.FileHeader) error {
	chunks, err := splitFileIntoChunks(file)
	if err != nil {
		return fmt.Errorf("error al dividir el archivo en fragmentos: %w", err)
	}

	uc.logger.Info(ctx, fmt.Sprintf("Generando embeddings con el servicio %s", embeddingProvider))

	if _, ok := uc.provider[embeddingProvider]; !ok {
		return fmt.Errorf("proveedor de embedding no válido: %s", embeddingProvider)
	}

	embeddings, err := uc.provider[embeddingProvider].GenerateEmbeddings(chunks)
	if err != nil {
		return err
	}

	if err := uc.repo.CreateBulk(ctx, embeddingProvider, embeddings); err != nil {
		return err
	}

	return nil
}

func splitFileIntoChunks(file *multipart.FileHeader) ([]string, error) {
	config := tokenizer.TokenChunkConfig{
		ChunkSize:    defaultChunkSize,
		ChunkOverlap: defaultChunkOverlap,
		ModelName:    defaultChunkModel,
	}

	fileBytes, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileBytes.Close()

	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(fileBytes); err != nil {
		return nil, err
	}

	chunks, err := tokenizer.SplitIntoTokenChunks(buffer.String(), config)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func float64ToFloat32Slice(input []float64) []float32 {
	output := make([]float32, len(input))
	for i, v := range input {
		output[i] = float32(v)
	}
	return output
}

// Fine tunning
func (uc UseCase) improveQuestion(question string) (string, error) {
	systemPrompt := `
		Eres un asistente de IA y tu objetivo es mejorar el prompt enviado por el usuario con el fin de obtener una respuesta más precisa ya que generaremos vectores de la pregunta para compararla con la base de datos de vectores que tenemos para nuestro RAG.

		Debes entregar solamente la pregunta mejorada, no la respuesta. No agregues información adicional. Responde solo en español.
	`

	firstPrompt := `
		El usuario ha enviado la siguiente pregunta:
		<Question>
			%s
		</Question>
	`

	parsedPrompt := fmt.Sprintf(firstPrompt, question)
	messagesImprover := ollama.Messages{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: parsedPrompt},
	}

	response, err := uc.provider[domain.Ollama].Chat(messagesImprover, domain.Gemma3, nil)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (uc UseCase) embedQuestion(embeddingProvider domain.EmbeddingProvider, improvedQuestion string) (pgvector.Vector, error) {
	if _, ok := uc.provider[embeddingProvider]; !ok {
		return pgvector.Vector{}, fmt.Errorf("proveedor de embedding no válido: %s", embeddingProvider)
	}

	embeddings, err := uc.provider[embeddingProvider].GenerateEmbeddings([]string{improvedQuestion})
	if err != nil {
		return pgvector.Vector{}, err
	}

	return pgvector.NewVector(embeddings[0].Embedding), nil
}

func (uc UseCase) searchEmbedInDB(ctx context.Context, embeddingProvider domain.EmbeddingProvider, embeddedQuestion pgvector.Vector, limit int) (string, error) {
	answer, err := uc.repo.Search(ctx, embeddingProvider, embeddedQuestion, limit)
	if err != nil {
		return "", err
	}

	var allAnswers bytes.Buffer
	for _, answer := range answer {
		allAnswers.WriteString(fmt.Sprintf("%s\n", answer.Paragraph))
	}

	return allAnswers.String(), nil
}

func (uc UseCase) askToIA(messagesAssist ollama.Messages, contexts, improvedQuestion string) (string, error) {
	prompt := `
		<Context>
			%s
		</Context>
		<Question>
			%s
		</Question>
	`

	parsedPrompt := fmt.Sprintf(prompt, contexts, improvedQuestion)
	messagesAssist = append(messagesAssist, ollama.Message{Role: "user", Content: parsedPrompt})

	response, err := uc.provider[domain.Ollama].Chat(
		messagesAssist,
		domain.DeepseekR1,
		&ollama.Options{
			Temperature: 0.7,
		})
	if err != nil {
		return "", err
	}

	return response, nil
}
