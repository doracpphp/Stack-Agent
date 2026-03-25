package task

import (
	appctx "stack-agent/context"
)

// Handler is implemented by every task type.
type Handler interface {
	Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error
}

// paramString safely extracts a string param, applying template rendering.
func paramString(params map[string]any, key string, ctx *appctx.ExecutionContext) (string, error) {
	v, ok := params[key]
	if !ok {
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", nil
	}
	return ctx.Render(s)
}

// paramStringRaw extracts a string param without rendering.
func paramStringRaw(params map[string]any, key string) string {
	v, ok := params[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

// paramStringSlice extracts a []string param.
func paramStringSlice(params map[string]any, key string) []string {
	v, ok := params[key]
	if !ok {
		return nil
	}
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
