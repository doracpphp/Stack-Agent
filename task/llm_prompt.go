package task

import (
	"fmt"

	appctx "stack-agent/context"
)

// LLMPromptHandler sends a prompt to the configured LLM and stores the response.
//
// Params:
//
//	prompt   string    - prompt text (supports {{var}} substitution)
//	images   []string  - list of task IDs whose image content to include (optional)
//	system   string    - system prompt prepended to the prompt (optional)
type LLMPromptHandler struct{}

func (h *LLMPromptHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	if ctx.LLM == nil {
		return fmt.Errorf("llm_prompt: no LLM provider configured")
	}

	prompt, err := paramString(params, "prompt", ctx)
	if err != nil {
		return fmt.Errorf("llm_prompt: rendering 'prompt': %w", err)
	}
	if prompt == "" {
		return fmt.Errorf("llm_prompt: required param 'prompt' is missing")
	}

	// Prepend system prompt if provided
	if system := paramStringRaw(params, "system"); system != "" {
		prompt = system + "\n\n" + prompt
	}

	// Collect image context
	imageIDs := paramStringSlice(params, "images")
	images := ctx.CollectImages(imageIDs)

	fmt.Printf("  → Calling %s...\n", ctx.LLM.Name())

	response, err := ctx.LLM.Complete(prompt, images)
	if err != nil {
		return fmt.Errorf("llm_prompt: LLM call failed: %w", err)
	}

	ctx.Set(taskID, &appctx.Content{Text: response})
	return nil
}
