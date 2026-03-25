package executor

import (
	"fmt"
	"io"
	"os"

	appctx "stack-agent/context"
	"stack-agent/pipeline"
	"stack-agent/task"
)

// Executor runs pipeline tasks sequentially.
type Executor struct {
	handlers map[string]task.Handler
	ctx      *appctx.ExecutionContext
	w        io.Writer
}

// New creates an Executor with all built-in task handlers registered.
// Logs are written to os.Stdout.
func New(ctx *appctx.ExecutionContext) *Executor {
	return NewWithWriter(ctx, os.Stdout)
}

// NewWithWriter creates an Executor that writes logs to the given writer.
func NewWithWriter(ctx *appctx.ExecutionContext, w io.Writer) *Executor {
	return &Executor{
		handlers: registerHandlers(),
		ctx:      ctx,
		w:        w,
	}
}

// Run executes all tasks in the pipeline in order.
func (e *Executor) Run(p *pipeline.Pipeline) error {
	total := len(p.Tasks)
	for i, t := range p.Tasks {
		h, ok := e.handlers[t.Type]
		if !ok {
			return fmt.Errorf("task[%d] %q: unknown task type %q", i+1, t.ID, t.Type)
		}

		fmt.Fprintf(e.w, "[%d/%d] %s (type: %s)\n", i+1, total, t.ID, t.Type)

		if err := h.Execute(t.Params, e.ctx, t.ID); err != nil {
			return fmt.Errorf("task[%d] %q failed: %w", i+1, t.ID, err)
		}
	}

	fmt.Fprintf(e.w, "\nパイプライン完了: %d タスクを実行しました\n", total)
	return nil
}

func registerHandlers() map[string]task.Handler {
	return map[string]task.Handler{
		"open_file":        &task.OpenFileHandler{},
		"open_files":       &task.OpenFilesHandler{},
		"open_folder":      &task.OpenFolderHandler{},
		"read_stdin":       &task.ReadStdinHandler{},
		"read_stdin_paths": &task.ReadStdinPathsHandler{},
		"read_image":       &task.ReadImageHandler{},
		"llm_prompt":       &task.LLMPromptHandler{},
		"fill_template":    &task.FillTemplateHandler{},
		"display_result":   &task.DisplayResultHandler{},
		"read_excel_cells": &task.ReadExcelCellsHandler{},
		"write_excel":      &task.WriteExcelHandler{},
		"write_file":       &task.WriteFileHandler{},
		"call_mcp":         &task.CallMCPHandler{},
	}
}
