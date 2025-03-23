package tokenizer

import (
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

type TokenChunkConfig struct {
	ChunkSize    int    // Number of tokens per chunk
	ChunkOverlap int    // Number of tokens to overlap between chunks
	ModelName    string // Name of the model to use for tokenization
}

func SplitIntoTokenChunks(text string, config TokenChunkConfig) ([]string, error) {
	// Available encondings: https://github.com/pkoukk/tiktoken-go?tab=readme-ov-file#available-encodings
	tikToken, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return nil, err
	}

	tokens := tikToken.Encode(text, nil, nil)

	var chunks []string

	if len(tokens) <= config.ChunkSize {
		return []string{text}, nil
	}

	for i := 0; i < len(tokens); i += config.ChunkSize - config.ChunkOverlap {
		end := i + config.ChunkSize
		if end > len(tokens) {
			end = len(tokens)
		}

		chunkTokens := tokens[i:end]
		chunkText := tikToken.Decode(chunkTokens)
		chunkText = strings.TrimSpace(string(chunkText))

		if chunkText != "" {
			chunks = append(chunks, chunkText)
		}

		if end == len(tokens) {
			break
		}
	}

	return chunks, nil
}
