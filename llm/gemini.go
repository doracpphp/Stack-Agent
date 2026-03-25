package llm

import (
	"context"
	"fmt"

	"google.golang.org/genai"

	appctx "stack-agent/context"
)

// GeminiProvider calls the Google Gemini API.
type GeminiProvider struct {
	model  string
	client *genai.Client
}

// NewGemini creates a Gemini provider.
func NewGemini(model, apiKey string) (*GeminiProvider, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini: creating client: %w", err)
	}
	return &GeminiProvider{model: model, client: client}, nil
}

func (g *GeminiProvider) Name() string {
	return fmt.Sprintf("gemini/%s", g.model)
}

func (g *GeminiProvider) Complete(prompt string, images []*appctx.Content) (string, error) {
	ctx := context.Background()

	var parts []*genai.Part
	parts = append(parts, genai.NewPartFromText(prompt))

	for _, img := range images {
		parts = append(parts, genai.NewPartFromBytes(img.Binary, img.MIME))
	}

	contents := []*genai.Content{
		{Parts: parts, Role: "user"},
	}

	resp, err := g.client.Models.GenerateContent(ctx, g.model, contents, nil)
	if err != nil {
		return "", fmt.Errorf("gemini: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("gemini: empty response")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("gemini: no content in response")
	}

	return candidate.Content.Parts[0].Text, nil
}
