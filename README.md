# stack-agent

複数のAIプロバイダー（Ollama・OpenAI・Gemini）に対応したスタックベースのタスク自動実行エージェントです。
タスクをJSONファイルに定義し、順番に実行することでファイル読み込み・LLM処理・結果出力を自動化します。

## 特徴

- **オフライン動作対応**: Ollamaを使用してローカルLLMで完全オフライン実行が可能
- **マルチプロバイダー**: Ollama / OpenAI / Gemini を切り替え可能
- **スタックベース実行**: タスクをJSONで定義し、順番に実行
- **変数テンプレート**: `{{タスクID}}` 記法でタスク間のデータ受け渡しが可能
- **マルチモーダル対応**: 画像をLLMに渡せる（llava / GPT-4o / Gemini）

---

## インストール

```bash
git clone https://github.com/yourname/stack-agent
cd stack-agent
go build -o stack-agent .
```

---

## 使い方

```bash
./stack-agent <パイプラインJSONファイル>

# または -pipeline フラグで指定
./stack-agent -pipeline examples/gemini_fill_report.json
```

### APIキーの設定

```bash
export OPENAI_API_KEY="sk-..."
export GEMINI_API_KEY="AIza..."
```

JSONの `api_key` フィールドに直接記載することも可能です。

---

## パイプラインJSON仕様

### 基本構造

```json
{
  "name": "パイプライン名",
  "version": "1.0",
  "llm": {
    "provider": "ollama",
    "model": "llama3.2",
    "base_url": "http://localhost:11434",
    "api_key": ""
  },
  "tasks": [
    {
      "id": "タスクID",
      "type": "タスクタイプ",
      "params": { ... }
    }
  ]
}
```

### LLM設定 (`llm`)

| フィールド | 型 | 説明 |
|-----------|-----|------|
| `provider` | string | `"ollama"` / `"openai"` / `"gemini"` |
| `model` | string | モデル名（例: `"llama3.2"`, `"gpt-4o"`, `"gemini-2.5-flash"`）|
| `base_url` | string | APIエンドポイント（Ollama用。省略時は `http://localhost:11434`）|
| `api_key` | string | APIキー（省略時は環境変数 `OPENAI_API_KEY` / `GEMINI_API_KEY` を使用）|

### 変数テンプレート

タスクのパラメータ内で `{{タスクID}}` と書くと、そのタスクの出力内容に置き換えられます。

```json
{ "prompt": "以下を要約してください:\n\n{{document}}" }
```

---

## タスクタイプ一覧

### `open_file` — ファイル読み込み

1つのファイルを読み込み、テキストとしてコンテキストに保存します。

**対応フォーマット:**

| 拡張子 | 形式 |
|--------|------|
| `.pdf` | PDFファイル（ページ単位でテキスト抽出）|
| `.docx` | Word文書 |
| `.xlsx` / `.xls` | Excelファイル（全シートをタブ区切りで抽出）|
| `.txt` / `.md` / `.csv` / `.json` / `.xml` / `.html` | テキスト系ファイル |

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `path` | string | ✓ | ファイルパス（`{{変数}}` 使用可）|

```json
{
  "id": "doc",
  "type": "open_file",
  "params": { "path": "report.pdf" }
}
```

---

### `open_folder` — フォルダ内ファイル一括読み込み

指定フォルダ内のファイルをすべて自動で読み込みます。ユーザーの入力は不要です。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `path` | string | ✓ | フォルダパス（`{{変数}}` 使用可）|
| `extensions` | []string | - | 読み込む拡張子（省略時: 全対応形式）|
| `recursive` | bool | - | サブフォルダも再帰的に読み込む（省略時: `false`）|
| `separator` | string | - | ファイル間の区切り文字（省略時: `\n\n---\n\n`）|

```json
{
  "id": "all_docs",
  "type": "open_folder",
  "params": {
    "path": "input/reports",
    "extensions": [".pdf", ".txt"],
    "recursive": true
  }
}
```

---

### `open_files` — 複数ファイル読み込み

複数のファイルをまとめて読み込み、連結してコンテキストに保存します。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `paths` | []string | ✓ | ファイルパスのリスト（各要素に `{{変数}}` 使用可）|
| `separator` | string | - | ファイル間の区切り文字（省略時: `\n\n---\n\n`）|

```json
{
  "id": "docs",
  "type": "open_files",
  "params": {
    "paths": ["{{path_a}}", "{{path_b}}", "notes.txt"]
  }
}
```

---

### `read_stdin` — 標準入力から読み込み

ユーザーに入力を促し、入力されたテキストまたはファイルの内容を保存します。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `prompt` | string | - | 表示するプロンプト文字列（省略時: `"入力してください: "`）|
| `read_file` | bool | - | `true` の場合、入力をファイルパスとして読み込む（省略時: `true`）|

> シングルクォート・ダブルクォートで囲まれたパスは自動的に除去されます。

```json
{
  "id": "user_input",
  "type": "read_stdin",
  "params": {
    "prompt": "ファイルパスを入力してください: ",
    "read_file": false
  }
}
```

---

### `read_stdin_paths` — 複数ファイルパスを標準入力から読み込み

ファイルパスを1行ずつ入力させ（空行で終了）、全ファイルを読み込んで連結します。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `prompt` | string | - | 最初に表示するメッセージ |
| `separator` | string | - | ファイル間の区切り文字（省略時: `\n\n---\n\n`）|

```json
{
  "id": "files",
  "type": "read_stdin_paths",
  "params": {
    "prompt": "参照ファイルを入力してください（空行で終了）:"
  }
}
```

---

### `read_image` — 画像読み込み

画像ファイルを読み込み、マルチモーダルLLMへの入力として保存します。

**対応フォーマット:** `.jpg` / `.jpeg` / `.png` / `.gif` / `.webp`

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `path` | string | ✓ | 画像ファイルパス（`{{変数}}` 使用可）|

```json
{
  "id": "photo",
  "type": "read_image",
  "params": { "path": "diagram.png" }
}
```

---

### `llm_prompt` — LLMへのプロンプト送信

LLMにプロンプトを送信し、レスポンスをコンテキストに保存します。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `prompt` | string | ✓ | プロンプト本文（`{{変数}}` 使用可）|
| `system` | string | - | システムプロンプト（プロンプトの前に付加）|
| `images` | []string | - | 画像として渡すタスクIDのリスト（`read_image` の出力を指定）|

```json
{
  "id": "summary",
  "type": "llm_prompt",
  "params": {
    "system": "あなたは文書要約の専門家です。",
    "prompt": "以下を3行で要約してください:\n\n{{doc}}",
    "images": ["photo"]
  }
}
```

---

### `fill_template` — テンプレートへのデータ流し込み

テンプレートファイルを読み込み、指定したデータを参照してLLMにテンプレートを埋めさせます。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `template` | string | ✓ | テンプレートファイルパス（`{{変数}}` 使用可）|
| `data` | []string | ✓ | 参照するタスクIDのリスト |
| `instruction` | string | - | LLMへの指示（省略時はデフォルト指示を使用）|

```json
{
  "id": "report",
  "type": "fill_template",
  "params": {
    "template": "templates/monthly_report.txt",
    "data": ["minutes", "progress", "issues"],
    "instruction": "テンプレートの各項目を提供データで埋めてください。"
  }
}
```

---

### `display_result` — 結果の表示

コンテキストの内容を標準出力に表示します。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `value` | string | ✓ | 表示する内容（`{{変数}}` 使用可）|
| `label` | string | - | 見出しラベル |

```json
{
  "id": "show",
  "type": "display_result",
  "params": {
    "label": "要約結果",
    "value": "{{summary}}"
  }
}
```

---

### `write_file` — ファイルへの書き込み

テキストをファイルに書き込みます。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `path` | string | ✓ | 出力ファイルパス（`{{変数}}` 使用可）|
| `value` | string | ✓ | 書き込む内容（`{{変数}}` 使用可）|
| `append` | bool | - | `true` の場合、追記モード（省略時: `false`）|

```json
{
  "id": "save",
  "type": "write_file",
  "params": {
    "path": "output.txt",
    "value": "{{summary}}"
  }
}
```

---

### `write_excel` — Excelセルへの書き込み

Excelファイルの指定セルに値を書き込みます。ファイルが存在しない場合は新規作成します。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `path` | string | ✓ | Excelファイルパス（`{{変数}}` 使用可）|
| `sheet` | string | - | シート名（省略時: `"Sheet1"`）|
| `cell` | string | ✓ | セル番地（例: `"A1"`, `"B3"`）|
| `value` | string | ✓ | 書き込む値（`{{変数}}` 使用可）|

```json
{
  "id": "excel_out",
  "type": "write_excel",
  "params": {
    "path": "report.xlsx",
    "sheet": "分析結果",
    "cell": "A1",
    "value": "{{analysis}}"
  }
}
```

---

### `call_mcp` — MCPツール呼び出し

MCPサーバーをサブプロセスとして起動し、指定ツールを呼び出します。

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `command` | string | ✓ | MCPサーバーの起動コマンド（例: `"npx"`）|
| `args` | []string | - | コマンド引数 |
| `tool` | string | ✓ | 呼び出すツール名 |
| `arguments` | object | - | ツールへの引数（文字列値に `{{変数}}` 使用可）|

```json
{
  "id": "mcp_result",
  "type": "call_mcp",
  "params": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
    "tool": "read_file",
    "arguments": { "path": "{{target_path}}" }
  }
}
```

---

## サンプル一覧

| ファイル | 内容 |
|---------|------|
| `examples/summarize_document.json` | ファイルをOllamaで要約してtxtに保存 |
| `examples/excel_analysis.json` | ExcelデータをOllamaで分析してシートに書き込み |
| `examples/image_describe.json` | 画像をllavaで説明（マルチモーダル）|
| `examples/mcp_filesystem.json` | MCPファイルシステムサーバー経由でファイルを読み要約 |
| `examples/gemini_document_qa.json` | ファイル読み込み＋Geminiへの質問応答 |
| `examples/gemini_image_analysis.json` | 画像をGeminiで詳細分析してtxtに保存 |
| `examples/gemini_excel_report.json` | ExcelデータをGeminiで分析してExcelに書き戻し |
| `examples/gemini_multi_file_summary.json` | 複数ファイルを対話入力してGeminiで要約 |
| `examples/gemini_multi_file_compare.json` | 2ファイルをGeminiで比較レポート生成 |
| `examples/gemini_fixed_files_report.json` | 固定4ファイルから年次報告書を生成 |
| `examples/gemini_fill_report.json` | 議事録・進捗・課題からテンプレートに流し込み月次報告書生成 |
| `examples/gemini_fill_report_interactive.json` | テンプレートとデータファイルを対話入力して報告書生成 |
| `examples/gemini_folder_summary.json` | フォルダ内ファイルを自動読み込みして一括要約 |
| `examples/gemini_folder_report.json` | フォルダ内ファイルを自動読み込みしてテンプレートで報告書生成 |

---

## ディレクトリ構成

```
stack-agent/
├── main.go                          # エントリーポイント
├── pipeline/
│   └── pipeline.go                  # JSONスキーマ定義・ロード
├── context/
│   └── context.go                   # タスク間共有コンテキスト・テンプレートエンジン
├── executor/
│   └── executor.go                  # タスクの順次実行・ハンドラー登録
├── task/
│   ├── handler.go                   # Handler インターフェース・ユーティリティ
│   ├── open_file.go                 # open_file
│   ├── open_files.go                # open_files
│   ├── open_folder.go               # open_folder
│   ├── read_stdin.go                # read_stdin
│   ├── read_stdin_paths.go          # read_stdin_paths
│   ├── read_image.go                # read_image
│   ├── llm_prompt.go                # llm_prompt
│   ├── fill_template.go             # fill_template
│   ├── display_result.go            # display_result
│   ├── write_file.go                # write_file
│   ├── write_excel.go               # write_excel
│   └── call_mcp.go                  # call_mcp
├── llm/
│   ├── provider.go                  # LLMProvider インターフェース・ファクトリ
│   ├── ollama.go                    # Ollama実装
│   ├── openai.go                    # OpenAI実装
│   └── gemini.go                    # Gemini実装
├── office/
│   ├── pdf.go                       # PDF読み込み
│   ├── word.go                      # Word(.docx)読み込み
│   ├── excel.go                     # Excel読み書き
│   └── image.go                     # 画像読み込み
├── mcp/
│   └── client.go                    # MCPクライアント（stdio）
└── examples/
    ├── templates/                   # 報告書テンプレート
    ├── sample_data/                 # サンプルデータ
    └── *.json                       # サンプルパイプライン
```
