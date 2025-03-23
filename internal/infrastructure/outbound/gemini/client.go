package gemini

import (
	"agentos/internal/domain"
	"agentos/internal/infrastructure/outbound/ollama"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Url    string
	ApiKey string
	Model  string
}

func NewClient(
	Url string,
	ApiKey string,
	Model string,
) Client {
	return Client{
		ApiKey: ApiKey,
		Model:  Model,
	}
}

func (c Client) GenerateEmbeddings(chunks []string) (domain.Embeddings, error) {
	url := fmt.Sprintf("%s/models/%s:embedContent?key=%s", c.Url, c.Model, c.ApiKey)

	var embeddings []domain.Embedding
	for _, chunk := range chunks {
		payload := Payload{}
		payload.Content.Parts = append(payload.Content.Parts, Part{Text: chunk})

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("error creando payload: %w", err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("GEMINI Client: error en la API: %s", string(body))
		}

		var response struct {
			Embedding struct {
				Values []float32 `json:"values"`
			} `json:"embedding"`
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		if len(response.Embedding.Values) == 0 {
			return nil, fmt.Errorf("la API no retornó embeddings")
		}

		embeddings = append(embeddings, domain.Embedding{
			Paragraph: chunk,
			Embedding: response.Embedding.Values,
		})
	}

	return embeddings, nil
}

func (c Client) Chat(messages ollama.Messages, model domain.OllamaModel, options *ollama.Options) (string, error) {
	return "", nil
}
