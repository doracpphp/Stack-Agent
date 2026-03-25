package task

import (
	"fmt"

	appctx "stack-agent/context"
	"stack-agent/office"
)

// WriteExcelHandler writes a value to a specific cell in an Excel file.
//
// Params:
//
//	path   string  - Excel file path (supports {{var}} substitution)
//	sheet  string  - sheet name (default: "Sheet1")
//	cell   string  - cell address, e.g. "A1", "B2"
//	value  string  - value to write (supports {{var}} substitution)
type WriteExcelHandler struct{}

func (h *WriteExcelHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	path, err := paramString(params, "path", ctx)
	if err != nil {
		return fmt.Errorf("write_excel: rendering 'path': %w", err)
	}
	if path == "" {
		return fmt.Errorf("write_excel: required param 'path' is missing")
	}

	sheet := paramStringRaw(params, "sheet")
	if sheet == "" {
		sheet = "Sheet1"
	}

	cell := paramStringRaw(params, "cell")
	if cell == "" {
		return fmt.Errorf("write_excel: required param 'cell' is missing")
	}

	value, err := paramString(params, "value", ctx)
	if err != nil {
		return fmt.Errorf("write_excel: rendering 'value': %w", err)
	}

	if err := office.WriteExcel(path, sheet, cell, value); err != nil {
		return fmt.Errorf("write_excel: %w", err)
	}

	fmt.Printf("  → Wrote to %s!%s in %s\n", sheet, cell, path)
	ctx.Set(taskID, &appctx.Content{Text: value})
	return nil
}
