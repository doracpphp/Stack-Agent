package llm

import (
	"context"
	"encoding/base64"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	appctx "stack-agent/context"
)

// OpenAIProvider calls the OpenAI API (or compatible endpoint).
type OpenAIProvider struct {
	model  string
	client *openai.Client
}

// NewOpenAI creates an OpenAI provider. baseURL is optional for custom endpoints.
func NewOpenAI(model, baseURL, apiKey string) *OpenAIProvider {
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	return &OpenAIProvider{
		model:  model,
		client: openai.NewClientWithConfig(cfg),
	}
}

func (o *OpenAIProvider) Name() string {
	return fmt.Sprintf("openai/%s", o.model)
}

func (o *OpenAIProvider) Complete(prompt string, images []*appctx.Content) (string, error) {
	ctx := context.Background()

	// Build message content
	var parts []openai.ChatMessagePart

	if len(images) > 0 {
		// Multimodal: add text part first
		parts = append(parts, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeText,
			Text: prompt,
		})
		for _, img := range images {
			encoded := base64.StdEncoding.EncodeToString(img.Binary)
			dataURL := fmt.Sprintf("data:%s;base64,%s", img.MIME, encoded)
			parts = append(parts, openai.ChatMessagePart{
				Type: openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{
					URL: dataURL,
				},
			})
		}
	}

	var messages []openai.ChatCompletionMessage
	if len(parts) > 0 {
		messages = []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, MultiContent: parts},
		}
	} else {
		messages = []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		}
	}

	req := openai.ChatCompletionRequest{
		Model:    o.model,
		Messages: messages,
	}

	resp, err := o.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("openai: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai: empty response")
	}

	return resp.Choices[0].Message.Content, nil
}
