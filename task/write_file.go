package task

import (
	"fmt"
	"os"
	"path/filepath"

	appctx "stack-agent/context"
)

// WriteFileHandler writes text content to a file.
//
// Params:
//
//	path    string  - output file path (supports {{var}} substitution)
//	value   string  - content to write (supports {{var}} substitution)
//	append  bool    - if true, append instead of overwrite (default: false)
type WriteFileHandler struct{}

func (h *WriteFileHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	path, err := paramString(params, "path", ctx)
	if err != nil {
		return fmt.Errorf("write_file: rendering 'path': %w", err)
	}
	if path == "" {
		return fmt.Errorf("write_file: required param 'path' is missing")
	}

	value, err := paramString(params, "value", ctx)
	if err != nil {
		return fmt.Errorf("write_file: rendering 'value': %w", err)
	}

	// Ensure parent directory exists
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("write_file: creating directory %s: %w", dir, err)
		}
	}

	appendMode := false
	if v, ok := params["append"]; ok {
		if b, ok := v.(bool); ok {
			appendMode = b
		}
	}

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if appendMode {
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}

	f, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return fmt.Errorf("write_file: opening %s: %w", path, err)
	}
	defer f.Close()

	if _, err := f.WriteString(value); err != nil {
		return fmt.Errorf("write_file: writing to %s: %w", path, err)
	}

	fmt.Printf("  → Wrote to %s\n", path)
	ctx.Set(taskID, &appctx.Content{Text: value})
	return nil
}
