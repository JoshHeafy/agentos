package domain

type Configuration struct {
	AllowedOrigins []string
	AllowedMethods []string
	PortHTTP       uint
	Database       Database
	Server         string
	Ollama         OllamaConfig
	Openai         OpenaiConfig
	Gemini         GeminiConfig
}

type Database struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     uint
	Name     string
	SSLMode  string
}

type OllamaConfig struct {
	Url   string
	Model string
}

type OpenaiConfig struct {
	Url    string
	ApiKey string
	Model  string
}

type GeminiConfig struct {
	Url    string
	ApiKey string
	Model  string
}
