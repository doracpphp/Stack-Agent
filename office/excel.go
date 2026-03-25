package office

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ReadExcel extracts all text content from an Excel file.
func ReadExcel(path string) (string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return "", fmt.Errorf("opening xlsx: %w", err)
	}
	defer f.Close()

	var sb strings.Builder
	for _, sheetName := range f.GetSheetList() {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			return "", fmt.Errorf("reading sheet %q: %w", sheetName, err)
		}

		sb.WriteString(fmt.Sprintf("=== Sheet: %s ===\n", sheetName))
		for _, row := range rows {
			sb.WriteString(strings.Join(row, "\t"))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// WriteExcel writes a value to a specific cell in an Excel file.
// Creates the file if it doesn't exist.
func WriteExcel(path, sheet, cell, value string) error {
	var f *excelize.File
	var err error

	// Try to open existing file, create new if not found
	f, err = excelize.OpenFile(path)
	if err != nil {
		f = excelize.NewFile()
	}
	defer f.Close()

	// Create sheet if it doesn't exist
	if idx, _ := f.GetSheetIndex(sheet); idx == -1 {
		f.NewSheet(sheet)
	}

	if err := f.SetCellValue(sheet, cell, value); err != nil {
		return fmt.Errorf("setting cell %s!%s: %w", sheet, cell, err)
	}

	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("saving xlsx: %w", err)
	}

	return nil
}
