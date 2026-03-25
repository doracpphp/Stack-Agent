package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	appctx "stack-agent/context"
	"stack-agent/executor"
	"stack-agent/llm"
	"stack-agent/pipeline"
)

// App holds the Wails application state.
type App struct {
	ctx context.Context
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// OpenPipeline opens a file dialog and returns the pipeline JSON content.
func (a *App) OpenPipeline() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "パイプラインを開く",
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON ファイル", Pattern: "*.json"},
		},
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}
	return string(data), nil
}

// SavePipeline opens a save dialog and writes the pipeline JSON to disk.
func (a *App) SavePipeline(content string) (string, error) {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "パイプラインを保存",
		DefaultFilename: "pipeline.json",
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON ファイル", Pattern: "*.json"},
		},
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("保存に失敗しました: %w", err)
	}
	return path, nil
}

// RunPipeline executes the pipeline from a JSON string.
// Log lines are emitted as "pipeline:log" events. On completion "pipeline:done"
// is emitted with an empty string (success) or an error message.
func (a *App) RunPipeline(pipelineJSON string) {
	go func() {
		emit := func(line string) {
			runtime.EventsEmit(a.ctx, "pipeline:log", line)
		}

		var p pipeline.Pipeline
		if err := json.Unmarshal([]byte(pipelineJSON), &p); err != nil {
			runtime.EventsEmit(a.ctx, "pipeline:done", fmt.Sprintf("JSON パースエラー: %v", err))
			return
		}

		if err := pipeline.Validate(&p); err != nil {
			runtime.EventsEmit(a.ctx, "pipeline:done", fmt.Sprintf("パイプライン検証エラー: %v", err))
			return
		}

		emit(fmt.Sprintf("▶ パイプライン開始: %s (v%s)", p.Name, p.Version))
		if p.LLM.Provider != "" {
			emit(fmt.Sprintf("  LLM: %s / %s", p.LLM.Provider, p.LLM.Model))
		}
		emit("")

		execCtx := appctx.New()

		if p.LLM.Provider != "" {
			provider, err := llm.NewProvider(
				p.LLM.Provider,
				p.LLM.Model,
				p.LLM.BaseURL,
				p.LLM.APIKey,
			)
			if err != nil {
				runtime.EventsEmit(a.ctx, "pipeline:done",
					fmt.Sprintf("LLM 初期化エラー: %v", err))
				return
			}
			execCtx.LLM = provider
		}

		// logWriter sends each written line as a Wails event.
		lw := &logWriter{emit: emit}
		ex := executor.NewWithWriter(execCtx, lw)

		if err := ex.Run(&p); err != nil {
			runtime.EventsEmit(a.ctx, "pipeline:done", err.Error())
			return
		}
		runtime.EventsEmit(a.ctx, "pipeline:done", "")
	}()
}

// logWriter implements io.Writer and emits each line via a callback.
type logWriter struct {
	emit func(string)
	buf  strings.Builder
}

func (lw *logWriter) Write(p []byte) (int, error) {
	lw.buf.Write(p)
	for {
		s := lw.buf.String()
		idx := strings.IndexByte(s, '\n')
		if idx < 0 {
			break
		}
		lw.emit(s[:idx])
		lw.buf.Reset()
		lw.buf.WriteString(s[idx+1:])
	}
	return len(p), nil
}

// SelectInputFile opens a file dialog and returns the selected file path.
// This is used by tasks like open_file that need a file path.
func (a *App) SelectInputFile(title string) (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// SelectInputFolder opens a folder dialog and returns the selected directory path.
func (a *App) SelectInputFolder(title string) (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// ensure io is used
var _ io.Writer = (*logWriter)(nil)
