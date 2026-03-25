package office

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// wordDocument represents the minimal XML structure of word/document.xml
type wordDocument struct {
	XMLName xml.Name  `xml:"document"`
	Body    wordBody  `xml:"body"`
}

type wordBody struct {
	Paragraphs []wordParagraph `xml:"p"`
}

type wordParagraph struct {
	Runs []wordRun `xml:"r"`
}

type wordRun struct {
	Text wordText `xml:"t"`
}

type wordText struct {
	Value string `xml:",chardata"`
}

// ReadWord extracts plain text from a .docx file.
func ReadWord(path string) (string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return "", fmt.Errorf("opening docx: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name != "word/document.xml" {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("opening document.xml: %w", err)
		}
		defer rc.Close()

		data, err := io.ReadAll(rc)
		if err != nil {
			return "", fmt.Errorf("reading document.xml: %w", err)
		}

		return parseWordXML(data)
	}

	return "", fmt.Errorf("word/document.xml not found in %s", path)
}

func parseWordXML(data []byte) (string, error) {
	var doc wordDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("parsing document.xml: %w", err)
	}

	var sb strings.Builder
	for _, para := range doc.Body.Paragraphs {
		var paraText strings.Builder
		for _, run := range para.Runs {
			paraText.WriteString(run.Text.Value)
		}
		text := paraText.String()
		if text != "" {
			sb.WriteString(text)
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}
