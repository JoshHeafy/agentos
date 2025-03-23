package domain

type OllamaModel string

const (
	Llama3     OllamaModel = "llama-3.2:latest"
	DeepseekR1 OllamaModel = "deepseek-r1:latest"
	Gemma3     OllamaModel = "gemma3:1b"
)
