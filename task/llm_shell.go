package task

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	appctx "stack-agent/context"
)

// llmShellSystemPrompt instructs the LLM to output only a shell command.
const llmShellSystemPrompt = `You are a shell command generator.
Output ONLY the shell command to execute — no explanation, no markdown, no code fences.
A single line. Nothing else.`

// codeBlockRe matches optional markdown code fences (e.g. ` ```sh\n...\n``` `).
var codeBlockRe = regexp.MustCompile("(?s)```[a-zA-Z]*\\n?(.*?)```")

// LLMShellHandler asks the LLM to generate a shell command, then executes it.
// The LLM response is stored under "<taskID>_cmd" and the command output under taskID.
//
// Params:
//
//	prompt   string    - what to ask the LLM (supports {{var}} substitution)
//	system   string    - additional system instructions appended to the built-in system prompt
//	timeout  string    - command execution timeout in seconds (default: 60)
//	workdir  string    - working directory (default: current dir)
type LLMShellHandler struct{}

func (h *LLMShellHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	if ctx.LLM == nil {
		return fmt.Errorf("llm_shell: no LLM provider configured")
	}

	prompt, err := paramString(params, "prompt", ctx)
	if err != nil {
		return fmt.Errorf("llm_shell: rendering 'prompt': %w", err)
	}
	if prompt = strings.TrimSpace(prompt); prompt == "" {
		return fmt.Errorf("llm_shell: required param 'prompt' is missing")
	}

	sysPrompt := llmShellSystemPrompt
	if extra := paramStringRaw(params, "system"); extra != "" {
		sysPrompt = sysPrompt + "\n" + extra
	}
	fullPrompt := sysPrompt + "\n\n" + prompt

	timeout := defaultShellTimeout
	if raw := paramStringRaw(params, "timeout"); raw != "" {
		var secs int
		if _, err := fmt.Sscanf(raw, "%d", &secs); err != nil || secs <= 0 {
			return fmt.Errorf("llm_shell: invalid timeout %q", raw)
		}
		timeout = time.Duration(secs) * time.Second
	}

	workdir, err := paramString(params, "workdir", ctx)
	if err != nil {
		return fmt.Errorf("llm_shell: rendering 'workdir': %w", err)
	}

	fmt.Printf("  → Asking %s to generate command...\n", ctx.LLM.Name())

	raw, err := ctx.LLM.Complete(fullPrompt, nil)
	if err != nil {
		return fmt.Errorf("llm_shell: LLM call failed: %w", err)
	}

	command := extractCommand(raw)
	if command == "" {
		return fmt.Errorf("llm_shell: LLM returned an empty command (raw: %q)", raw)
	}

	// Store the generated command so subsequent tasks can inspect it.
	ctx.Set(taskID+"_cmd", &appctx.Content{Text: command})

	return runShellCommand(command, workdir, timeout, ctx, taskID, "llm_shell")
}

// extractCommand strips markdown code fences and surrounding whitespace from the LLM response.
func extractCommand(s string) string {
	s = strings.TrimSpace(s)
	if m := codeBlockRe.FindStringSubmatch(s); len(m) == 2 {
		s = strings.TrimSpace(m[1])
	}
	// Take only the first non-empty line in case the model emitted multiple.
	for _, line := range strings.Split(s, "\n") {
		if t := strings.TrimSpace(line); t != "" {
			return t
		}
	}
	return ""
}
