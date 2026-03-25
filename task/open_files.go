package task

import (
	"fmt"
	"path/filepath"
	"strings"

	appctx "stack-agent/context"
)

// OpenFilesHandler opens multiple files and concatenates their content.
//
// Params:
//
//	paths      []string  - list of file paths (each supports {{var}} substitution)
//	separator  string    - separator between files (default: "\n\n---\n\n")
type OpenFilesHandler struct{}

func (h *OpenFilesHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	rawPaths := paramStringSlice(params, "paths")
	if len(rawPaths) == 0 {
		return fmt.Errorf("open_files: required param 'paths' is missing or empty")
	}

	separator := "\n\n---\n\n"
	if v, ok := params["separator"]; ok {
		if s, ok := v.(string); ok {
			separator = s
		}
	}

	var parts []string
	for _, rawPath := range rawPaths {
		// Apply template substitution to each path
		path, err := ctx.Render(rawPath)
		if err != nil {
			return fmt.Errorf("open_files: rendering path %q: %w", rawPath, err)
		}

		text, err := readFile(path)
		if err != nil {
			return fmt.Errorf("open_files: %w", err)
		}

		header := fmt.Sprintf("### ファイル: %s\n\n", filepath.Base(path))
		parts = append(parts, header+text)
		fmt.Printf("  → 読み込み完了: %s\n", path)
	}

	combined := strings.Join(parts, separator)
	ctx.Set(taskID, &appctx.Content{Text: combined, MIME: "text/plain"})
	return nil
}
