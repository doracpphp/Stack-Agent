package task

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	appctx "stack-agent/context"
)

// ReadStdinHandler prompts the user for input via stdin.
// If the input looks like a file path, it reads that file automatically.
// Otherwise it stores the raw input string.
//
// Params:
//
//	prompt      string  - message displayed to user (optional)
//	read_file   bool    - if true, treat input as a file path and read it (default: true)
type ReadStdinHandler struct{}

func (h *ReadStdinHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	prompt := paramStringRaw(params, "prompt")
	if prompt == "" {
		prompt = "入力してください: "
	}

	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read_stdin: reading input: %w", err)
	}
	input := strings.TrimSpace(line)
	// シングル/ダブルクォートで囲まれている場合は除去
	if len(input) >= 2 {
		if (input[0] == '\'' && input[len(input)-1] == '\'') ||
			(input[0] == '"' && input[len(input)-1] == '"') {
			input = input[1 : len(input)-1]
		}
	}

	// Check if we should treat the input as a file path
	readFile := true
	if v, ok := params["read_file"]; ok {
		if b, ok := v.(bool); ok {
			readFile = b
		}
	}

	if readFile {
		text, err := readFile2(input)
		if err != nil {
			return fmt.Errorf("read_stdin: ファイルの読み込みに失敗しました %q: %w", input, err)
		}
		ctx.Set(taskID, &appctx.Content{Text: text, MIME: "text/plain"})
	} else {
		ctx.Set(taskID, &appctx.Content{Text: input})
	}

	return nil
}

// readFile2 reads a file given its path (avoids naming collision with open_file.go).
func readFile2(path string) (string, error) {
	return readFile(path)
}
