package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appctx "stack-agent/context"
	"stack-agent/office"
)

// OpenFileHandler opens a file (docx, xlsx, txt, etc.) and stores its text content.
//
// Params:
//
//	path  string  - file path (supports {{var}} substitution)
type OpenFileHandler struct{}

func (h *OpenFileHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	path, err := paramString(params, "path", ctx)
	if err != nil {
		return fmt.Errorf("open_file: rendering 'path': %w", err)
	}
	if path == "" {
		return fmt.Errorf("open_file: required param 'path' is missing")
	}

	text, err := readFile(path)
	if err != nil {
		return fmt.Errorf("open_file: %w", err)
	}

	ctx.Set(taskID, &appctx.Content{
		Text: text,
		MIME: "text/plain",
	})
	return nil
}

func readFile(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".pdf":
		return office.ReadPDF(path)
	case ".docx":
		return office.ReadWord(path)
	case ".xlsx", ".xls":
		return office.ReadExcel(path)
	case ".txt", ".md", ".csv", ".json", ".xml", ".html", ".htm":
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", path, err)
		}
		return string(data), nil
	default:
		// Attempt to read as text for unknown extensions
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", path, err)
		}
		return string(data), nil
	}
}
