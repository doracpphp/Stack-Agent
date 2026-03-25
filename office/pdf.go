package office

import (
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ReadPDF extracts plain text from a PDF file.
func ReadPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening pdf: %w", err)
	}
	defer f.Close()

	var sb strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("reading page %d: %w", i, err)
		}

		if text != "" {
			fmt.Fprintf(&sb, "--- ページ %d ---\n%s\n", i, text)
		}
	}

	return sb.String(), nil
}
