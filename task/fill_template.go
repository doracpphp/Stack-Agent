package task

import (
	"fmt"
	"os"

	appctx "stack-agent/context"
)

// FillTemplateHandler reads a report template file and asks the LLM to fill it
// using the provided source data.
//
// Params:
//
//	template  string    - path to the template file
//	data      []string  - list of task IDs whose content is used as source data
//	                      (if omitted, all current context variables are used)
//	instruction string  - additional instruction for the LLM (optional)
type FillTemplateHandler struct{}

func (h *FillTemplateHandler) Execute(params map[string]any, ctx *appctx.ExecutionContext, taskID string) error {
	if ctx.LLM == nil {
		return fmt.Errorf("fill_template: no LLM provider configured")
	}

	templatePath, err := paramString(params, "template", ctx)
	if err != nil {
		return fmt.Errorf("fill_template: rendering 'template': %w", err)
	}
	if templatePath == "" {
		return fmt.Errorf("fill_template: required param 'template' is missing")
	}

	tmplBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("fill_template: reading template file: %w", err)
	}
	templateContent := string(tmplBytes)

	// Collect source data
	dataIDs := paramStringSlice(params, "data")
	var dataSections string
	for _, id := range dataIDs {
		c, ok := ctx.Get(id)
		if !ok {
			return fmt.Errorf("fill_template: variable %q not found", id)
		}
		dataSections += fmt.Sprintf("### %s\n\n%s\n\n---\n\n", id, c.Text)
	}

	instruction := paramStringRaw(params, "instruction")
	if instruction == "" {
		instruction = "テンプレートの各項目に、提供されたデータから適切な情報を日本語で埋めてください。データに該当する情報がない項目は「情報なし」と記載してください。テンプレートの構造・フォーマットは厳守し、余分な説明は加えないでください。"
	}

	prompt := fmt.Sprintf(
		"## 指示\n\n%s\n\n## 報告書テンプレート\n\n%s\n\n## 参照データ\n\n%s",
		instruction,
		templateContent,
		dataSections,
	)

	fmt.Printf("  → テンプレートにデータを流し込み中 (%s)...\n", ctx.LLM.Name())

	result, err := ctx.LLM.Complete(prompt, nil)
	if err != nil {
		return fmt.Errorf("fill_template: LLM call failed: %w", err)
	}

	ctx.Set(taskID, &appctx.Content{Text: result})
	return nil
}
