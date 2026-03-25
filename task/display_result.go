package task

import (
	"fmt"

	appctx "stack-agent/context"
)

// DisplayResultHandler prints content to stdout.
//
// Params:
//
//	value   string  - value to display (supports {{var}} substitution)
//	label   string  - optional label displayed before the value
type DisplayResultHandler struct{}

func (h *DisplayResultHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	value, err := paramString(params, "value", ctx)
	if err != nil {
		return fmt.Errorf("display_result: rendering 'value': %w", err)
	}

	label := paramStringRaw(params, "label")
	if label != "" {
		fmt.Printf("\n=== %s ===\n", label)
	} else {
		fmt.Println("\n--- 結果 ---")
	}

	fmt.Println(value)
	fmt.Println("---")

	// Store displayed value under task ID for potential reuse
	ctx.Set(taskID, &appctx.Content{Text: value})
	return nil
}
