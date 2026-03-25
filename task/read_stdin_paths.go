package task

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appctx "stack-agent/context"
)

// ReadStdinPathsHandler prompts the user to enter multiple file paths (one per line).
// Input ends when the user presses Enter on an empty line.
// All files are read and their content is concatenated.
//
// Params:
//
//	prompt     string  - initial message (optional)
//	separator  string  - separator between files (default: "\n\n---\n\n")
type ReadStdinPathsHandler struct{}

func (h *ReadStdinPathsHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	prompt := paramStringRaw(params, "prompt")
	if prompt == "" {
		prompt = "ファイルパスを入力してください（空行で終了）:"
	}

	separator := "\n\n---\n\n"
	if v, ok := params["separator"]; ok {
		if s, ok := v.(string); ok {
			separator = s
		}
	}

	fmt.Println(prompt)

	reader := bufio.NewReader(os.Stdin)
	var parts []string
	lineNum := 1

	for {
		fmt.Printf("  [%d] パス: ", lineNum)
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		path := strings.TrimSpace(line)

		// シングル/ダブルクォートで囲まれている場合は除去
		if len(path) >= 2 {
			if (path[0] == '\'' && path[len(path)-1] == '\'') ||
				(path[0] == '"' && path[len(path)-1] == '"') {
				path = path[1 : len(path)-1]
			}
		}

		// 空行で終了
		if path == "" {
			break
		}

		text, err := readFile(path)
		if err != nil {
			fmt.Printf("  警告: %v (スキップします)\n", err)
			continue
		}

		header := fmt.Sprintf("### ファイル: %s\n\n", filepath.Base(path))
		parts = append(parts, header+text)
		fmt.Printf("  ✓ 読み込み完了: %s\n", path)
		lineNum++
	}

	if len(parts) == 0 {
		return fmt.Errorf("read_stdin_paths: ファイルが1つも読み込まれませんでした")
	}

	combined := strings.Join(parts, separator)
	ctx.Set(taskID, &appctx.Content{Text: combined, MIME: "text/plain"})
	fmt.Printf("  → %d ファイルを結合しました\n", len(parts))
	return nil
}
