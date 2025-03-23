package localconfig

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"agentos/internal/domain"
)

var envVars []string

func clearPreviousEnv() {
	for _, key := range envVars {
		os.Unsetenv(key)
	}
	envVars = []string{}
}

func loadEnvFile(path string) error {
	env := os.Getenv("ENV")
	if env != "local" {
		return nil
	}

	clearPreviousEnv()

	envMap, err := godotenv.Read(path)
	if err != nil {
		return err
	}

	for key, value := range envMap {
		os.Setenv(key, value)
		envVars = append(envVars, key)
	}

	return nil
}

func getEnvOrFail(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing %s environment variable", key)
	}

	return value, nil
}

func parseUint(value string, key string) (uint64, error) {
	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %w", key, err)
	}

	return parsed, nil
}

func NewLocalConfig(path string) (domain.Configuration, error) {
	_ = loadEnvFile(path)

	portHTTP, err := parseUint(os.Getenv("PORT_HTTP"), "PORT_HTTP")
	if err != nil {
		return domain.Configuration{}, err
	}

	dbPort, err := parseUint(os.Getenv("DB_PORT"), "DB_PORT")
	if err != nil {
		return domain.Configuration{}, err
	}

	dbEngine, err := getEnvOrFail("DB_ENGINE")
	if err != nil {
		return domain.Configuration{}, err
	}

	config := domain.Configuration{
		AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
		AllowedMethods: strings.Split(os.Getenv("ALLOWED_METHODS"), ","),
		PortHTTP:       uint(portHTTP),
		Database: domain.Database{
			Driver:   dbEngine,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_SERVER"),
			Port:     uint(dbPort),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Server: os.Getenv("SERVER"),
		Ollama: domain.OllamaConfig{
			Url:   os.Getenv("OLLAMA_URL"),
			Model: os.Getenv("OLLAMA_MODEL"),
		},
		Openai: domain.OpenaiConfig{
			Url:    os.Getenv("OPENAI_URL"),
			ApiKey: os.Getenv("OPENAI_API_KEY"),
			Model:  os.Getenv("OPENAI_MODEL"),
		},
		Gemini: domain.GeminiConfig{
			Url:    os.Getenv("GEMINI_URL"),
			ApiKey: os.Getenv("GEMINI_API_KEY"),
			Model:  os.Getenv("GEMINI_MODEL"),
		},
	}

	return config, nil
}
