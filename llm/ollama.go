package llm

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	appctx "stack-agent/context"
)

// OllamaProvider calls a local Ollama instance.
type OllamaProvider struct {
	model   string
	baseURL string
	client  *http.Client
}

// NewOllama creates an Ollama provider.
func NewOllama(model, baseURL string) *OllamaProvider {
	return &OllamaProvider{
		model:   model,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (o *OllamaProvider) Name() string {
	return fmt.Sprintf("ollama/%s", o.model)
}

type ollamaRequest struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Images []string `json:"images,omitempty"` // base64-encoded
	Stream bool     `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

func (o *OllamaProvider) Complete(prompt string, images []*appctx.Content) (string, error) {
	req := ollamaRequest{
		Model:  o.model,
		Prompt: prompt,
		Stream: false,
	}

	for _, img := range images {
		req.Images = append(req.Images, base64.StdEncoding.EncodeToString(img.Binary))
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("ollama: marshaling request: %w", err)
	}

	resp, err := o.client.Post(o.baseURL+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ollama: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ollama: reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama: HTTP %d: %s", resp.StatusCode, string(data))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(data, &ollamaResp); err != nil {
		return "", fmt.Errorf("ollama: parsing response: %w", err)
	}

	if ollamaResp.Error != "" {
		return "", fmt.Errorf("ollama: %s", ollamaResp.Error)
	}

	return ollamaResp.Response, nil
}
