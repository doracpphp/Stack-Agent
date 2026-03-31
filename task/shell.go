package task

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	appctx "stack-agent/context"
)

const defaultShellTimeout = 60 * time.Second

// ShellHandler executes a shell command and stores stdout+stderr as the result.
//
// Params:
//
//	command   string  - shell command to run (supports {{var}} substitution)
//	timeout   string  - max execution time in seconds (default: 60)
//	workdir   string  - working directory (default: current dir)
type ShellHandler struct{}

func (h *ShellHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	command, err := paramString(params, "command", ctx)
	if err != nil {
		return fmt.Errorf("shell: rendering 'command': %w", err)
	}
	if command = strings.TrimSpace(command); command == "" {
		return fmt.Errorf("shell: required param 'command' is missing")
	}

	timeout := defaultShellTimeout
	if raw := paramStringRaw(params, "timeout"); raw != "" {
		secs, err := strconv.Atoi(raw)
		if err != nil || secs <= 0 {
			return fmt.Errorf("shell: invalid timeout %q (must be positive integer seconds)", raw)
		}
		timeout = time.Duration(secs) * time.Second
	}

	workdir, err := paramString(params, "workdir", ctx)
	if err != nil {
		return fmt.Errorf("shell: rendering 'workdir': %w", err)
	}

	return runShellCommand(command, workdir, timeout, ctx, taskID, "shell")
}

// runShellCommand runs a shell command and stores combined stdout+stderr in ctx.
func runShellCommand(command, workdir string, timeout time.Duration, ctx *appctx.ExecutionContext, taskID, prefix string) error {
	fmt.Printf("  → $ %s\n", command)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	if workdir != "" {
		cmd.Dir = workdir
	}

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%s: starting command: %w", prefix, err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		output := buf.String()
		if err != nil {
			// Include output even on non-zero exit so next tasks can inspect it.
			return fmt.Errorf("%s: command failed: %w\noutput:\n%s", prefix, err, output)
		}
		ctx.Set(taskID, &appctx.Content{Text: output})
		return nil
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		return fmt.Errorf("%s: command timed out after %s", prefix, timeout)
	}
}
