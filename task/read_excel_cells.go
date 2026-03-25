package task

import (
	"fmt"

	appctx "stack-agent/context"
	"stack-agent/office"
)

// ReadExcelCellsHandler reads a cell or range from an Excel sheet.
//
// Params:
//
//	path   string  - Excel file path (supports {{var}} substitution)
//	sheet  string  - sheet name (default: first sheet)
//	range  string  - cell address or range, e.g. "A1" or "A1:C10"
type ReadExcelCellsHandler struct{}

func (h *ReadExcelCellsHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	path, err := paramString(params, "path", ctx)
	if err != nil {
		return fmt.Errorf("read_excel_cells: rendering 'path': %w", err)
	}
	if path == "" {
		return fmt.Errorf("read_excel_cells: required param 'path' is missing")
	}

	sheet := paramStringRaw(params, "sheet")

	cellRange := paramStringRaw(params, "range")
	if cellRange == "" {
		return fmt.Errorf("read_excel_cells: required param 'range' is missing")
	}

	value, err := office.ReadExcelCells(path, sheet, cellRange)
	if err != nil {
		return fmt.Errorf("read_excel_cells: %w", err)
	}

	fmt.Printf("  → Read %s!%s from %s\n", sheet, cellRange, path)
	ctx.Set(taskID, &appctx.Content{Text: value})
	return nil
}
