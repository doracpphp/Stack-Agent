package llm

import (
	"fmt"
	"os"

	appctx "stack-agent/context"
)

// Provider is the unified LLM interface.
type Provider interface {
	Complete(prompt string, images []*appctx.Content) (string, error)
	Name() string
}

// NewProvider constructs the appropriate provider from config.
func NewProvider(providerName, model, baseURL, apiKey string) (Provider, error) {
	// Resolve API key from env if not specified in config
	switch providerName {
	case "openai":
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		return NewOpenAI(model, baseURL, apiKey), nil

	case "gemini":
		if apiKey == "" {
			apiKey = os.Getenv("GEMINI_API_KEY")
		}
		return NewGemini(model, apiKey)

	case "ollama":
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		return NewOllama(model, baseURL), nil

	default:
		return nil, fmt.Errorf("unknown LLM provider: %q (supported: openai, gemini, ollama)", providerName)
	}
}
