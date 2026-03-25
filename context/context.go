package context

import (
	"fmt"
	"regexp"
	"strings"
)

// Content holds the payload produced by a task.
type Content struct {
	Text   string // plain text (document, LLM response, etc.)
	Binary []byte // raw bytes for images
	MIME   string // e.g. "image/png"
}

// LLMProvider is the interface the context uses to call the LLM.
// Defined here to avoid circular imports.
type LLMProvider interface {
	Complete(prompt string, images []*Content) (string, error)
	Name() string
}

// ExecutionContext is the shared mutable state threaded through all tasks.
type ExecutionContext struct {
	vars map[string]*Content
	LLM  LLMProvider
}

// New creates an empty ExecutionContext.
func New() *ExecutionContext {
	return &ExecutionContext{
		vars: make(map[string]*Content),
	}
}

// Set stores content under the given ID.
func (c *ExecutionContext) Set(id string, content *Content) {
	c.vars[id] = content
}

// Get retrieves content by ID.
func (c *ExecutionContext) Get(id string) (*Content, bool) {
	v, ok := c.vars[id]
	return v, ok
}

var templateRe = regexp.MustCompile(`\{\{(\w+)\}\}`)

// Render substitutes {{var}} placeholders with their text values.
func (c *ExecutionContext) Render(template string) (string, error) {
	var missing []string
	result := templateRe.ReplaceAllStringFunc(template, func(match string) string {
		key := match[2 : len(match)-2]
		if content, ok := c.vars[key]; ok {
			return content.Text
		}
		missing = append(missing, key)
		return match
	})
	if len(missing) > 0 {
		return "", fmt.Errorf("undefined variables in template: %s", strings.Join(missing, ", "))
	}
	return result, nil
}

// CollectImages returns Content items that carry binary image data for the given IDs.
func (c *ExecutionContext) CollectImages(ids []string) []*Content {
	var out []*Content
	for _, id := range ids {
		if v, ok := c.vars[id]; ok && len(v.Binary) > 0 {
			out = append(out, v)
		}
	}
	return out
}
