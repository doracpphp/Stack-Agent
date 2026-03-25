package main

import (
	"flag"
	"fmt"
	"os"

	appctx "stack-agent/context"
	"stack-agent/executor"
	"stack-agent/llm"
	"stack-agent/pipeline"
)

func main() {
	pipelineFile := flag.String("pipeline", "pipeline.json", "パイプラインJSONファイルのパス")
	flag.Parse()

	// Also support positional argument
	if flag.NArg() > 0 {
		*pipelineFile = flag.Arg(0)
	}

	p, err := pipeline.Load(*pipelineFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: パイプラインの読み込みに失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("パイプライン: %s (v%s)\n", p.Name, p.Version)
	fmt.Printf("LLM: %s / %s\n\n", p.LLM.Provider, p.LLM.Model)

	ctx := appctx.New()

	// Only init LLM provider if LLM is configured
	if p.LLM.Provider != "" {
		provider, err := llm.NewProvider(
			p.LLM.Provider,
			p.LLM.Model,
			p.LLM.BaseURL,
			p.LLM.APIKey,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: LLMプロバイダーの初期化に失敗しました: %v\n", err)
			os.Exit(1)
		}
		ctx.LLM = provider
	}

	ex := executor.New(ctx)
	if err := ex.Run(p); err != nil {
		fmt.Fprintf(os.Stderr, "\nエラー: %v\n", err)
		os.Exit(1)
	}
}
