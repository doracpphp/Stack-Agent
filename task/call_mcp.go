package task

import (
	"encoding/json"
	"fmt"

	appctx "stack-agent/context"
	"stack-agent/mcp"
)

// CallMCPHandler calls a tool on an MCP server.
//
// Params:
//
//	command    string    - MCP server command (e.g. "npx")
//	args       []string  - command arguments (e.g. ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"])
//	tool       string    - tool name to call
//	arguments  object    - tool arguments (string values support {{var}} substitution)
type CallMCPHandler struct{}

func (h *CallMCPHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	command := paramStringRaw(params, "command")
	if command == "" {
		return fmt.Errorf("call_mcp: required param 'command' is missing")
	}

	args := paramStringSlice(params, "args")

	toolName := paramStringRaw(params, "tool")
	if toolName == "" {
		return fmt.Errorf("call_mcp: required param 'tool' is missing")
	}

	// Render template variables in arguments
	arguments, err := renderArguments(params["arguments"], ctx)
	if err != nil {
		return fmt.Errorf("call_mcp: rendering arguments: %w", err)
	}

	fmt.Printf("  → Calling MCP tool %q on %s\n", toolName, command)

	client, err := mcp.NewClient(command, args)
	if err != nil {
		return fmt.Errorf("call_mcp: starting MCP server: %w", err)
	}
	defer client.Close()

	result, err := client.CallTool(toolName, arguments)
	if err != nil {
		return fmt.Errorf("call_mcp: tool call failed: %w", err)
	}

	// Extract text from MCP result
	text := extractMCPText(result)
	ctx.Set(taskID, &appctx.Content{Text: text})
	return nil
}

// renderArguments recursively renders {{var}} templates in the arguments map.
func renderArguments(raw any, ctx *appctx.ExecutionContext) (map[string]any, error) {
	if raw == nil {
		return map[string]any{}, nil
	}

	m, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("arguments must be an object")
	}

	result := make(map[string]any, len(m))
	for k, v := range m {
		if s, ok := v.(string); ok {
			rendered, err := ctx.Render(s)
			if err != nil {
				return nil, fmt.Errorf("key %q: %w", k, err)
			}
			result[k] = rendered
		} else {
			result[k] = v
		}
	}
	return result, nil
}

// extractMCPText tries to extract human-readable text from an MCP tool result.
func extractMCPText(raw json.RawMessage) string {
	// MCP tool results are typically: {"content": [{"type": "text", "text": "..."}]}
	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(raw, &result); err == nil && len(result.Content) > 0 {
		for _, c := range result.Content {
			if c.Type == "text" && c.Text != "" {
				return c.Text
			}
		}
	}

	// Fall back to raw JSON string
	return string(raw)
}
