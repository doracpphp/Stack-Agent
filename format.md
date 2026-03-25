# stack-agent パイプライン JSON フォーマット仕様

このドキュメントはGUIエディタ（Wails デスクトップアプリ）およびCLIでパイプラインを組み立てるために必要な、JSONフォーマットの完全仕様です。

---

## 1. トップレベル構造

```json
{
  "name": "string",
  "version": "string",
  "llm": { ... },
  "tasks": [ ... ]
}
```

| フィールド | 型 | 必須 | 説明 |
|-----------|-----|:----:|------|
| `name` | string | ✓ | パイプラインの表示名 |
| `version` | string | - | バージョン番号（例: `"1.0"`）|
| `llm` | object | - | LLM設定。`llm_prompt` / `fill_template` を使う場合は必須 |
| `tasks` | array | ✓ | タスクの配列。上から順番に実行される |

---

## 2. LLM設定オブジェクト (`llm`)

```json
{
  "provider": "ollama",
  "model": "gpt-oss:20b",
  "base_url": "http://localhost:11434",
  "api_key": ""
}
```

| フィールド | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `provider` | string | ✓ | LLMプロバイダー。下記の選択肢から選ぶ | - |
| `model` | string | ✓ | モデル名 | - |
| `base_url` | string | - | APIエンドポイントURL | Ollamaのみ: `http://localhost:11434` |
| `api_key` | string | - | APIキー。空の場合は環境変数を使用 | `""` |

### `provider` の選択肢

| 値 | 説明 | 必要なもの |
|----|------|-----------|
| `"ollama"` | ローカルLLM。オフライン動作 | Ollamaのインストール |
| `"openai"` | OpenAI API | `OPENAI_API_KEY` または `api_key` |
| `"gemini"` | Google Gemini API | `GEMINI_API_KEY` または `api_key` |

### モデル名の例

| provider | model 例 |
|----------|---------|
| `ollama` | `"gpt-oss:20b"`, `"llama3.2"`, `"llava"`, `"mistral"`, `"gemma3"` |
| `openai` | `"gpt-4o"`, `"gpt-4o-mini"`, `"gpt-4-turbo"` |
| `gemini` | `"gemini-2.5-flash"`, `"gemini-2.0-flash"`, `"gemini-1.5-pro"` |

---

## 3. タスクオブジェクト

```json
{
  "id": "string",
  "type": "string",
  "params": { ... }
}
```

| フィールド | 型 | 必須 | 説明 |
|-----------|-----|:----:|------|
| `id` | string | ✓ | タスクの一意なID。他のタスクから `{{id}}` で参照できる |
| `type` | string | ✓ | タスクの種類。下記タスクタイプ一覧から選ぶ |
| `params` | object | - | タスクの設定パラメータ。タスクタイプごとに異なる |

### ID のルール
- 英数字とアンダースコアのみ使用可（例: `my_task`, `step1`）
- パイプライン内で重複不可
- 他のタスクのパラメータから `{{このID}}` として参照される

### 変数テンプレート記法
パラメータ内の文字列値に `{{タスクID}}` と書くと、そのタスクの出力テキストに置換される。

```json
{ "prompt": "以下を要約してください:\n\n{{document}}" }
```

---

## 4. タスクタイプ一覧

### 4-1. `open_file` — ファイル読み込み

ファイルを読み込み、テキストとして保存する。

**カテゴリ:** 入力
**出力:** ファイルの内容（テキスト）

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `path` | string | ✓ | ファイルパス（`{{変数}}` 使用可）| - |

**対応ファイル形式:**

| 拡張子 | 形式 | 備考 |
|--------|------|------|
| `.pdf` | PDF | ページ番号付きでテキスト抽出 |
| `.docx` | Word文書 | 段落テキストを抽出 |
| `.xlsx`, `.xls` | Excel | 全シートをタブ区切りテキストで抽出 |
| `.txt` | テキスト | そのまま読み込み |
| `.md` | Markdown | そのまま読み込み |
| `.csv` | CSV | そのまま読み込み |
| `.json` | JSON | そのまま読み込み |
| `.xml` | XML | そのまま読み込み |
| `.html`, `.htm` | HTML | そのまま読み込み |
| その他 | - | テキストとして読み込みを試みる |

```json
{
  "id": "document",
  "type": "open_file",
  "params": {
    "path": "report.pdf"
  }
}
```

---

### 4-2. `open_files` — 複数ファイル読み込み

複数のファイルを読み込み、連結して保存する。各ファイルの先頭にファイル名のヘッダーが付く。

**カテゴリ:** 入力
**出力:** 全ファイルの内容を連結したテキスト

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `paths` | string[] | ✓ | ファイルパスの配列（各要素に `{{変数}}` 使用可）| - |
| `separator` | string | - | ファイル間の区切り文字列 | `"\n\n---\n\n"` |

```json
{
  "id": "all_docs",
  "type": "open_files",
  "params": {
    "paths": ["{{path_a}}", "{{path_b}}", "fixed_file.txt"],
    "separator": "\n\n===\n\n"
  }
}
```

---

### 4-3. `open_folder` — フォルダ内ファイル一括読み込み

指定フォルダ内のファイルをすべて読み込み、連結して保存する。ユーザーの入力なしで実行される。

**カテゴリ:** 入力
**出力:** フォルダ内全ファイルの内容を連結したテキスト

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `path` | string | ✓ | フォルダパス（`{{変数}}` 使用可）| - |
| `extensions` | string[] | - | 読み込む拡張子のリスト（例: `[".txt", ".pdf"]`）| 全対応形式 |
| `recursive` | bool | - | `true` の場合、サブフォルダも再帰的に読み込む | `false` |
| `separator` | string | - | ファイル間の区切り文字列 | `"\n\n---\n\n"` |

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

### 4-4. `read_stdin` — キーボード入力待ち

実行時にユーザーに入力を求める。ファイルパスとして読み込むか、テキストとして保存するかを選べる。

**カテゴリ:** 入力
**出力:** 入力されたテキスト、またはファイルの内容

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `prompt` | string | - | ユーザーに表示するプロンプト文字列 | `"入力してください: "` |
| `read_file` | bool | - | `true`: 入力をファイルパスとして扱いファイルを読み込む。`false`: 入力テキストをそのまま保存 | `true` |

> シングルクォート（`'`）またはダブルクォート（`"`）で囲まれたパスは自動的に除去される。

```json
{
  "id": "user_question",
  "type": "read_stdin",
  "params": {
    "prompt": "質問を入力してください: ",
    "read_file": false
  }
}
```

---

### 4-5. `read_stdin_paths` — 複数ファイルパスをキーボード入力

複数のファイルパスを1行ずつ入力させ（空行で終了）、全ファイルを読み込んで連結する。

**カテゴリ:** 入力
**出力:** 全ファイルの内容を連結したテキスト

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `prompt` | string | - | 最初に表示するメッセージ | `"ファイルパスを入力してください（空行で終了）:"` |
| `separator` | string | - | ファイル間の区切り文字列 | `"\n\n---\n\n"` |

```json
{
  "id": "input_files",
  "type": "read_stdin_paths",
  "params": {
    "prompt": "参照するファイルを入力してください（空行で終了）:"
  }
}
```

---

### 4-6. `read_image` — 画像読み込み

画像ファイルを読み込む。`llm_prompt` の `images` パラメータと組み合わせてマルチモーダルLLMに渡せる。

**カテゴリ:** 入力
**出力:** 画像バイナリデータ

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `path` | string | ✓ | 画像ファイルパス（`{{変数}}` 使用可）| - |

**対応形式:** `.jpg` / `.jpeg` / `.png` / `.gif` / `.webp`

```json
{
  "id": "photo",
  "type": "read_image",
  "params": {
    "path": "{{image_path}}"
  }
}
```

---

### 4-7. `read_excel_cells` — Excelセル読み込み

Excelファイルの特定セルまたはセル範囲を読み込む。

**カテゴリ:** 入力
**出力:** セルの値（単一セル）またはタブ区切り・改行区切りのテーブルテキスト（範囲）

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `path` | string | ✓ | Excelファイルパス（`{{変数}}` 使用可）| - |
| `sheet` | string | - | シート名 | 先頭シート |
| `range` | string | ✓ | セル番地または範囲（例: `"A1"`, `"A1:C10"`）| - |

**出力フォーマット:**

- 単一セル（`"A1"`）: そのセルの文字列値
- 範囲（`"A1:C3"`）:
  ```
  値A1\t値B1\t値C1
  値A2\t値B2\t値C2
  値A3\t値B3\t値C3
  ```

```json
{
  "id": "table_data",
  "type": "read_excel_cells",
  "params": {
    "path": "data.xlsx",
    "sheet": "売上",
    "range": "A1:D20"
  }
}
```

---

### 4-8. `llm_prompt` — LLMへのプロンプト送信

LLMにプロンプトを送信し、レスポンスを保存する。

**カテゴリ:** AI処理
**出力:** LLMのレスポンステキスト
**LLM設定:** 必須

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `prompt` | string | ✓ | プロンプト本文（`{{変数}}` 使用可）| - |
| `system` | string | - | システムプロンプト（プロンプトの前に付加）| - |
| `images` | string[] | - | 画像として渡す `read_image` タスクのIDリスト | `[]` |

```json
{
  "id": "answer",
  "type": "llm_prompt",
  "params": {
    "system": "あなたは優秀なアシスタントです。",
    "prompt": "以下のデータを分析してください:\n\n{{table_data}}",
    "images": ["diagram"]
  }
}
```

---

### 4-9. `fill_template` — テンプレートへのデータ流し込み

テンプレートファイルを読み込み、指定データをもとにLLMがテンプレートの各項目を埋める。

**カテゴリ:** AI処理
**出力:** テンプレートが埋められたテキスト
**LLM設定:** 必須

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `template` | string | ✓ | テンプレートファイルパス（`{{変数}}` 使用可）| - |
| `data` | string[] | ✓ | 参照するタスクIDのリスト | - |
| `instruction` | string | - | LLMへの追加指示 | デフォルト指示文を使用 |

```json
{
  "id": "report",
  "type": "fill_template",
  "params": {
    "template": "templates/monthly_report.txt",
    "data": ["minutes", "progress", "issues"],
    "instruction": "テンプレートの[...]部分を提供データで埋めてください。"
  }
}
```

---

### 4-10. `display_result` — 結果の表示

コンテキストの内容をターミナル（GUIではログパネル）に表示する。

**カテゴリ:** 出力
**出力:** 表示した値（そのままコンテキストにも保存）

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `value` | string | ✓ | 表示する値（`{{変数}}` 使用可）| - |
| `label` | string | - | 見出しラベル | `"結果"` |

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

### 4-11. `write_file` — ファイルへの書き込み

テキストをファイルに書き込む。親ディレクトリが存在しない場合は自動作成。

**カテゴリ:** 出力
**出力:** 書き込んだ値（そのままコンテキストにも保存）

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `path` | string | ✓ | 出力ファイルパス（`{{変数}}` 使用可）| - |
| `value` | string | ✓ | 書き込む内容（`{{変数}}` 使用可）| - |
| `append` | bool | - | `true` の場合、上書きせず末尾に追記 | `false` |

```json
{
  "id": "save",
  "type": "write_file",
  "params": {
    "path": "output/result.md",
    "value": "{{report}}",
    "append": false
  }
}
```

---

### 4-12. `write_excel` — Excelセルへの書き込み

Excelファイルの指定セルに値を書き込む。ファイル・シートが存在しない場合は自動作成。

**カテゴリ:** 出力
**出力:** 書き込んだ値（そのままコンテキストにも保存）

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `path` | string | ✓ | Excelファイルパス（`{{変数}}` 使用可）| - |
| `sheet` | string | - | シート名 | `"Sheet1"` |
| `cell` | string | ✓ | セル番地（例: `"A1"`, `"C3"`）| - |
| `value` | string | ✓ | 書き込む値（`{{変数}}` 使用可）| - |

```json
{
  "id": "excel_out",
  "type": "write_excel",
  "params": {
    "path": "report.xlsx",
    "sheet": "AI分析",
    "cell": "B2",
    "value": "{{analysis}}"
  }
}
```

---

### 4-13. `call_mcp` — MCPツール呼び出し

MCPサーバーをサブプロセスとして起動し、指定ツールを呼び出す。呼び出し後にサーバーは終了する。

**カテゴリ:** 外部連携
**出力:** ツールの返却値（テキスト）

| パラメータ | 型 | 必須 | 説明 | デフォルト |
|-----------|-----|:----:|------|-----------|
| `command` | string | ✓ | MCPサーバーの起動コマンド（例: `"npx"`, `"python"`）| - |
| `args` | string[] | - | コマンドの引数リスト | `[]` |
| `tool` | string | ✓ | 呼び出すツール名 | - |
| `arguments` | object | - | ツールへの引数（文字列値に `{{変数}}` 使用可）| `{}` |

```json
{
  "id": "mcp_result",
  "type": "call_mcp",
  "params": {
    "command": "npx",
    "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
    "tool": "read_file",
    "arguments": {
      "path": "{{target_path}}"
    }
  }
}
```

---

## 5. タスクタイプ早見表

| タイプ | カテゴリ | LLM必要 | 主な入力 | 主な出力 |
|--------|---------|:-------:|---------|---------|
| `open_file` | 入力 | - | ファイルパス | ファイルのテキスト |
| `open_files` | 入力 | - | ファイルパス[] | 連結テキスト |
| `open_folder` | 入力 | - | フォルダパス（固定）| フォルダ内全ファイルの連結テキスト |
| `read_stdin` | 入力 | - | キーボード入力 | テキスト or ファイル内容 |
| `read_stdin_paths` | 入力 | - | キーボード入力（複数行）| 連結テキスト |
| `read_image` | 入力 | - | 画像ファイルパス | 画像バイナリ |
| `read_excel_cells` | 入力 | - | Excelパス + セル/範囲 | セル値 or テーブルテキスト |
| `llm_prompt` | AI処理 | ✓ | プロンプト + `{{変数}}` | LLMレスポンス |
| `fill_template` | AI処理 | ✓ | テンプレートファイル + データ変数 | 埋め込み済みテキスト |
| `display_result` | 出力 | - | `{{変数}}` | ターミナル/ログパネル表示 |
| `write_file` | 出力 | - | `{{変数}}` | ファイル |
| `write_excel` | 出力 | - | `{{変数}}` | Excelセル |
| `call_mcp` | 外部連携 | - | ツール名 + 引数 | ツール返却値 |

---

## 6. データフロー

タスクは上から順に実行され、各タスクの出力は `id` をキーとしてコンテキストに蓄積される。後続タスクは `{{id}}` でその値を参照できる。

```
tasks[0] → コンテキスト["id0"] = "..."
tasks[1] → コンテキスト["id1"] = "...{{id0}}..." → 展開して処理
tasks[2] → コンテキスト["id2"] = "...{{id0}}...{{id1}}..."
```

### フロー図の例（Excelデータ分析）

```
[read_excel_cells] table ──→ [llm_prompt] analysis → [write_excel] result
```

### フロー図の例（月次報告書生成）

```
[open_file] minutes   ─┐
[open_file] progress  ─┼→ [fill_template] report → [display_result] → [write_file]
[open_file] issues    ─┘                                               [write_excel]
```

---

## 7. バリデーションルール

| ルール | 説明 |
|--------|------|
| `tasks` は1件以上 | 空のパイプラインは無効 |
| `id` は必須 | 空文字不可 |
| `id` は一意 | パイプライン内で重複不可 |
| `id` は英数字・アンダースコアのみ | スペース・記号不可 |
| `type` は必須 | 未定義のタスクタイプは実行時エラー |
| `{{変数}}` の参照先は存在すること | 未定義IDへの参照は実行時エラー |
| `llm_prompt` / `fill_template` 使用時は `llm` 設定が必要 | プロバイダー・モデルは必須 |
| `read_image` の出力を `llm_prompt.images` に渡すこと | テキスト系タスクの出力を `images` に指定してもバイナリが存在しない |

---

## 8. 完全なパイプライン例

### 例1: Excelセルを読み込んでLLMで分析

```json
{
  "name": "売上データ分析",
  "version": "1.0",
  "llm": {
    "provider": "ollama",
    "model": "gpt-oss:20b",
    "base_url": "http://localhost:11434",
    "api_key": ""
  },
  "tasks": [
    {
      "id": "table_data",
      "type": "read_excel_cells",
      "params": {
        "path": "data/sales.xlsx",
        "sheet": "売上",
        "range": "A1:D13"
      }
    },
    {
      "id": "analysis",
      "type": "llm_prompt",
      "params": {
        "system": "あなたはデータ分析の専門家です。",
        "prompt": "以下の売上データを分析し、傾向と改善点を日本語で述べてください:\n\n{{table_data}}"
      }
    },
    {
      "id": "show",
      "type": "display_result",
      "params": {
        "label": "分析結果",
        "value": "{{analysis}}"
      }
    },
    {
      "id": "save",
      "type": "write_excel",
      "params": {
        "path": "data/sales.xlsx",
        "sheet": "分析",
        "cell": "A1",
        "value": "{{analysis}}"
      }
    }
  ]
}
```

### 例2: 月次業務報告書自動生成

```json
{
  "name": "月次業務報告書自動生成",
  "version": "1.0",
  "llm": {
    "provider": "gemini",
    "model": "gemini-2.5-flash",
    "api_key": ""
  },
  "tasks": [
    {
      "id": "minutes",
      "type": "open_file",
      "params": {
        "path": "data/meeting_minutes.txt"
      }
    },
    {
      "id": "progress",
      "type": "open_file",
      "params": {
        "path": "data/task_progress.txt"
      }
    },
    {
      "id": "issues",
      "type": "open_file",
      "params": {
        "path": "data/issues.txt"
      }
    },
    {
      "id": "report",
      "type": "fill_template",
      "params": {
        "template": "templates/monthly_report.txt",
        "data": ["minutes", "progress", "issues"]
      }
    },
    {
      "id": "show",
      "type": "display_result",
      "params": {
        "label": "月次報告書",
        "value": "{{report}}"
      }
    },
    {
      "id": "save_md",
      "type": "write_file",
      "params": {
        "path": "output/monthly_report.md",
        "value": "{{report}}"
      }
    },
    {
      "id": "save_xlsx",
      "type": "write_excel",
      "params": {
        "path": "output/monthly_report.xlsx",
        "sheet": "報告書",
        "cell": "A1",
        "value": "{{report}}"
      }
    }
  ]
}
```
