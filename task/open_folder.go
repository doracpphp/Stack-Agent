package task

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	appctx "stack-agent/context"
)

// OpenFolderHandler reads all files in a specified folder and concatenates their content.
// Subdirectories are skipped unless recursive is true.
//
// Params:
//
//	path        string    - folder path (supports {{var}} substitution)
//	extensions  []string  - file extensions to include, e.g. [".txt", ".pdf"]
//	                        if omitted, all supported file types are included
//	recursive   bool      - if true, also read files in subdirectories (default: false)
//	separator   string    - separator between files (default: "\n\n---\n\n")
type OpenFolderHandler struct{}

// supportedExts is the set of extensions readable by readFile()
var supportedExts = map[string]bool{
	".pdf":  true,
	".docx": true,
	".xlsx": true,
	".xls":  true,
	".txt":  true,
	".md":   true,
	".csv":  true,
	".json": true,
	".xml":  true,
	".html": true,
	".htm":  true,
}

func (h *OpenFolderHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	folderPath, err := paramString(params, "path", ctx)
	if err != nil {
		return fmt.Errorf("open_folder: rendering 'path': %w", err)
	}
	if folderPath == "" {
		return fmt.Errorf("open_folder: required param 'path' is missing")
	}

	info, err := os.Stat(folderPath)
	if err != nil {
		return fmt.Errorf("open_folder: folder not found: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("open_folder: %q is not a directory", folderPath)
	}

	// Build allowed extension set
	filterExts := buildExtFilter(paramStringSlice(params, "extensions"))

	recursive := false
	if v, ok := params["recursive"]; ok {
		if b, ok := v.(bool); ok {
			recursive = b
		}
	}

	separator := "\n\n---\n\n"
	if v, ok := params["separator"]; ok {
		if s, ok := v.(string); ok {
			separator = s
		}
	}

	// Collect file paths
	files, err := collectFiles(folderPath, filterExts, recursive)
	if err != nil {
		return fmt.Errorf("open_folder: scanning directory: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("open_folder: no readable files found in %q", folderPath)
	}

	// Read each file
	var parts []string
	for _, filePath := range files {
		text, err := readFile(filePath)
		if err != nil {
			fmt.Printf("  警告: %v (スキップします)\n", err)
			continue
		}
		rel, _ := filepath.Rel(folderPath, filePath)
		header := fmt.Sprintf("### ファイル: %s\n\n", rel)
		parts = append(parts, header+text)
		fmt.Printf("  → 読み込み完了: %s\n", rel)
	}

	if len(parts) == 0 {
		return fmt.Errorf("open_folder: すべてのファイルの読み込みに失敗しました")
	}

	combined := strings.Join(parts, separator)
	ctx.Set(taskID, &appctx.Content{Text: combined, MIME: "text/plain"})
	fmt.Printf("  → %d ファイルを読み込みました\n", len(parts))
	return nil
}

// collectFiles walks the directory and returns sorted file paths matching the filter.
func collectFiles(dir string, filterExts map[string]bool, recursive bool) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			if recursive {
				sub, err := collectFiles(fullPath, filterExts, true)
				if err != nil {
					return nil, err
				}
				files = append(files, sub...)
			}
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if filterExts[ext] {
			files = append(files, fullPath)
		}
	}

	sort.Strings(files)
	return files, nil
}

// buildExtFilter returns a set of lowercase extensions to include.
// If the user specified extensions, use those; otherwise use all supported types.
func buildExtFilter(specified []string) map[string]bool {
	if len(specified) == 0 {
		return supportedExts
	}
	m := make(map[string]bool, len(specified))
	for _, ext := range specified {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		m[strings.ToLower(ext)] = true
	}
	return m
}
