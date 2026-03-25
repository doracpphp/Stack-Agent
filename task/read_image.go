package task

import (
	"fmt"
	"path/filepath"

	appctx "stack-agent/context"
	"stack-agent/office"
)

// ReadImageHandler reads an image file and stores it for use in multimodal LLM prompts.
//
// Params:
//
//	path  string  - image file path (supports {{var}} substitution)
type ReadImageHandler struct{}

func (h *ReadImageHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	path, err := paramString(params, "path", ctx)
	if err != nil {
		return fmt.Errorf("read_image: rendering 'path': %w", err)
	}
	if path == "" {
		return fmt.Errorf("read_image: required param 'path' is missing")
	}

	data, mime, err := office.ReadImage(path)
	if err != nil {
		return fmt.Errorf("read_image: %w", err)
	}

	ctx.Set(taskID, &appctx.Content{
		Text:   fmt.Sprintf("[image: %s]", filepath.Base(path)),
		Binary: data,
		MIME:   mime,
	})
	return nil
}
