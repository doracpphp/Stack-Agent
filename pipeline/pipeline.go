package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
)

// Pipeline is the top-level JSON document.
type Pipeline struct {
	Name    string    `json:"name"`
	Version string    `json:"version"`
	LLM     LLMConfig `json:"llm"`
	Tasks   []Task    `json:"tasks"`
}

// LLMConfig selects which provider and model to use.
type LLMConfig struct {
	Provider string `json:"provider"` // "openai" | "gemini" | "ollama"
	Model    string `json:"model"`
	BaseURL  string `json:"base_url"` // optional override, required for Ollama
	APIKey   string `json:"api_key"`  // optional, falls back to env vars
}

// Task is a single unit of work in the pipeline.
type Task struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"`
	Params map[string]any `json:"params"`
}

// Load reads and parses a pipeline JSON file.
func Load(path string) (*Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading pipeline file: %w", err)
	}

	var p Pipeline
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parsing pipeline JSON: %w", err)
	}

	if err := Validate(&p); err != nil {
		return nil, fmt.Errorf("invalid pipeline: %w", err)
	}

	return &p, nil
}

// Validate checks a Pipeline for structural errors. Exported so callers that
// build a Pipeline from sources other than a JSON file (e.g. the GUI) can
// reuse the same rules without going through Load.
func Validate(p *Pipeline) error {
	if len(p.Tasks) == 0 {
		return fmt.Errorf("pipeline has no tasks")
	}

	seen := make(map[string]bool)
	for i, t := range p.Tasks {
		if t.ID == "" {
			return fmt.Errorf("task[%d] missing required field 'id'", i)
		}
		if t.Type == "" {
			return fmt.Errorf("task[%d] %q missing required field 'type'", i, t.ID)
		}
		if seen[t.ID] {
			return fmt.Errorf("task[%d] duplicate id %q", i, t.ID)
		}
		seen[t.ID] = true
	}
	return nil
}
