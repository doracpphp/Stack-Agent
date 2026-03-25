package office

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// maxImageBytes is the maximum image size accepted (50 MB).
// LLM APIs typically have their own limits, but this prevents OOM from
// accidentally pointing at a huge binary file.
const maxImageBytes = 50 * 1024 * 1024

// ReadImage reads an image file and returns its raw bytes and MIME type.
func ReadImage(path string) ([]byte, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, "", fmt.Errorf("reading image: %w", err)
	}
	if info.Size() > maxImageBytes {
		return nil, "", fmt.Errorf("image file too large: %d bytes (max %d MB)", info.Size(), maxImageBytes/1024/1024)
	}

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
