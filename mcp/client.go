package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync/atomic"
	"time"
)

const mcpRequestTimeout = 2 * time.Minute

// Client manages a single MCP server subprocess over stdio.
type Client struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	nextID atomic.Int64
}

type jsonRPCRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      int64          `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewClient spawns the MCP server process and performs the initialization handshake.
func NewClient(command string, args []string) (*Client, error) {
	cmd := exec.Command(command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("mcp: getting stdin pipe: %w", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("mcp: getting stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("mcp: starting server %q: %w", command, err)
	}

	c := &Client{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdoutPipe),
	}

	// MCP initialization handshake
	if err := c.initialize(); err != nil {
		c.Close()
		return nil, fmt.Errorf("mcp: initialization failed: %w", err)
	}

	return c, nil
}

func (c *Client) initialize() error {
	initReq := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      c.nextID.Add(1),
		Method:  "initialize",
		Params: map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo": map[string]any{
				"name":    "stack-agent",
				"version": "1.0.0",
			},
		},
	}

	if _, err := c.sendRequest(initReq); err != nil {
		return err
	}

	// Send initialized notification
	notif := map[string]any{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	}
	data, _ := json.Marshal(notif)
	data = append(data, '\n')
	_, err := c.stdin.Write(data)
	return err
}

// CallTool sends a tools/call request and returns the raw result.
func (c *Client) CallTool(toolName string, arguments map[string]any) (json.RawMessage, error) {
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      c.nextID.Add(1),
		Method:  "tools/call",
		Params: map[string]any{
			"name":      toolName,
			"arguments": arguments,
		},
	}

	resp, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func (c *Client) sendRequest(req jsonRPCRequest) (*jsonRPCResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("mcp: marshaling request: %w", err)
	}
	data = append(data, '\n')

	if _, err := c.stdin.Write(data); err != nil {
		return nil, fmt.Errorf("mcp: writing request: %w", err)
	}

	// Read the response with a timeout to avoid blocking forever on a hung server.
	type readResult struct {
		line string
		err  error
	}
	ch := make(chan readResult, 1)
	go func() {
		line, err := c.stdout.ReadString('\n')
		ch <- readResult{line, err}
	}()

	select {
	case r := <-ch:
		if r.err != nil {
			return nil, fmt.Errorf("mcp: reading response: %w", r.err)
		}
		var resp jsonRPCResponse
		if err := json.Unmarshal([]byte(r.line), &resp); err != nil {
			return nil, fmt.Errorf("mcp: parsing response: %w", err)
		}
		// Verify the response belongs to this request to prevent mismatched replies.
		if resp.ID != req.ID {
			return nil, fmt.Errorf("mcp: response ID mismatch (got %d, want %d)", resp.ID, req.ID)
		}
		if resp.Error != nil {
			return nil, fmt.Errorf("mcp: server error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return &resp, nil

	case <-time.After(mcpRequestTimeout):
		return nil, fmt.Errorf("mcp: request timed out after %s", mcpRequestTimeout)
	}
}

// Close terminates the subprocess.
func (c *Client) Close() error {
	c.stdin.Close()
	return c.cmd.Process.Kill()
}
