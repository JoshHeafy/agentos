package openai

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
		Url:    Url,
		ApiKey: ApiKey,
		Model:  Model,
	}
}

func (c Client) GenerateEmbeddings(chunks []string) (domain.Embeddings, error) {
	apiURL := fmt.Sprintf("%s/embeddings", c.Url)

	payload := map[string]any{
		"input": chunks,
		"model": c.Model,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error creando payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en la solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la API: %s", string(bodyBytes))
	}

	var response struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parseando respuesta de la API: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("la API no retornó embeddings")
	}

	var result domain.Embeddings
	for i, embedding := range response.Data {
		result = append(result, domain.Embedding{
			Paragraph: chunks[i],
			Embedding: embedding.Embedding,
		})
	}

	return result, nil
}

func (c Client) Chat(messages ollama.Messages, model domain.OllamaModel, options *ollama.Options) (string, error) {
	return "", nil
}
