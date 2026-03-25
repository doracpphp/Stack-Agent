package office

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadImage reads an image file and returns its raw bytes and MIME type.
func ReadImage(path string) ([]byte, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("reading image: %w", err)
	}

	mime := mimeFromExt(filepath.Ext(path))
	return data, mime, nil
}

func mimeFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}
