package ollama

import (
	"agentos/internal/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Url   string
	Model string
}

func NewClient(
	Url string,
	Model string,
) Client {
	return Client{
		Url:   Url,
		Model: Model,
	}
}

func (c Client) GenerateEmbeddings(chunks []string) (domain.Embeddings, error) {
	url := c.Url + "/api/embed"

	payload := map[string]any{
		"model": c.Model,
		"input": chunks,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error creando payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	var response struct {
		Embeddings [][]float32 `json:"embeddings"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %w", err)
	}

	if len(response.Embeddings) == 0 {
		return nil, fmt.Errorf("error al generar el embedding")
	}

	var result domain.Embeddings

	for i, chunk := range chunks {
		result = append(result, domain.Embedding{
			Paragraph: chunk,
			Embedding: response.Embeddings[i],
		})
	}

	return result, nil
}

func (c Client) Chat(messages Messages, model domain.OllamaModel, options *Options) (string, error) {
	url := c.Url + "/api/chat"

	if model == "" {
		model = domain.OllamaModel(c.Model)
	}

	payload := map[string]any{
		"model":    model,
		"messages": messages,
		"stream":   false,
		"options":  options,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error creando payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error creando solicitud: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error en solicitud HTTP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error en la solicitud HTTP: %s", string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error parseando respuesta: %w", err)
	}

	return response.Message.Content, nil
}
