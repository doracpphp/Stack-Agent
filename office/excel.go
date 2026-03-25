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

// ReadExcelCells reads a single cell or a range (e.g. "A1" or "A1:C3") from an Excel sheet.
// If sheet is empty the first sheet is used.
// A single cell returns just that cell's value.
// A range returns tab-separated columns, one row per line.
func ReadExcelCells(path, sheet, cellOrRange string) (string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return "", fmt.Errorf("opening xlsx: %w", err)
	}
	defer f.Close()

	if sheet == "" {
		list := f.GetSheetList()
		if len(list) == 0 {
			return "", fmt.Errorf("no sheets found in %s", path)
		}
		sheet = list[0]
	}

	if !strings.Contains(cellOrRange, ":") {
		// Single cell
		val, err := f.GetCellValue(sheet, cellOrRange)
		if err != nil {
			return "", fmt.Errorf("reading cell %s!%s: %w", sheet, cellOrRange, err)
		}
		return val, nil
	}

	// Range e.g. "A1:C3"
	parts := strings.SplitN(cellOrRange, ":", 2)
	startCol, startRow, err := excelize.CellNameToCoordinates(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid cell address %q: %w", parts[0], err)
	}
	endCol, endRow, err := excelize.CellNameToCoordinates(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid cell address %q: %w", parts[1], err)
	}

	var sb strings.Builder
	for r := startRow; r <= endRow; r++ {
		for c := startCol; c <= endCol; c++ {
			if c > startCol {
				sb.WriteByte('\t')
			}
			name, _ := excelize.CoordinatesToCellName(c, r)
			val, _ := f.GetCellValue(sheet, name)
			sb.WriteString(val)
		}
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n"), nil
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
