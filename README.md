# stack-agent

> **「AIを使った自動化に興味はあるけど、黒い画面（ターミナル）は難しい」「何から始めればいいかわからない」という方のためのツールです。**

タスクをビジュアルエディタで組み立てて、ボタン1つで実行できる AI パイプラインツールです。
プログラムを書く必要はありません。JSON のブロックを並べるだけで、ファイル読み込み・LLM処理・結果の保存まで自動化できます。

---

## こんな方に

- ChatGPT や Gemini は使っているが、毎回同じ操作を手動でやっている
- 大量のファイルや Excel データを AI に一括処理させたい
- ターミナル（黒い画面）は使いたくない
- 社内に配布できる「ダブルクリックで起動するツール」を作りたい

---

## 特徴

- **ダブルクリックで起動**: インストール不要の単一実行ファイル（.app / .exe）
- **ビジュアルエディタ内蔵**: ブロックを並べるだけでパイプラインを組み立て
- **ボタン1つで実行**: エディタから離れずにそのまま実行・ログ確認
- **オフライン動作対応**: Ollama を使えばインターネット接続なしで動作
- **マルチプロバイダー**: Ollama（ローカル）/ OpenAI / Gemini を切り替え可能
- **変数テンプレート**: `{{タスクID}}` 記法でタスク間のデータ受け渡し
- **マルチモーダル対応**: 画像を LLM に渡せる（llava / GPT-4o / Gemini）
- **日本語 / English 対応**: UI の表示言語を切り替え可能

---

## できること（パイプライン例）

```
📄 Excelのセルを読み込む → 🤖 AIで分析 → 📊 結果をExcelに書き戻す
📁 フォルダ内のPDFを全部読む → 🤖 AIで要約 → 💾 テキストに保存
📝 議事録・進捗・課題を読む → 🤖 AIで月次報告書を作成 → 💾 Wordテンプレートに流し込む
🖼️ 画像を読み込む → 🤖 AIで内容を説明 → 💬 画面に表示
```

---

## インストール・ビルド

### GUIアプリ（推奨）

[Wails v2](https://wails.io) が必要です。

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

```bash
git clone https://github.com/yourname/stack-agent
cd stack-agent

# macOS (.app)
wails build

# Windows (.exe)
wails build --platform windows/amd64
```

生成物: `build/bin/stack-agent.app`（macOS）/ `build/bin/stack-agent.exe`（Windows）

> 配布する場合は生成された `.app` / `.exe` ファイルを渡すだけです。受け取った側は追加インストール不要でダブルクリックで起動できます。

### CLI ツール（上級者向け）

```bash
go build -o stack-agent ./cmd/cli
./stack-agent <パイプラインJSONファイル>
```

---

## 使い方

### 1. アプリを起動する

`stack-agent.app`（macOS）または `stack-agent.exe`（Windows）をダブルクリックします。

### 2. パイプラインを組み立てる

左のパネルからタスクブロックをクリックして追加し、パラメータを入力します。

| 操作 | 説明 |
|------|------|
| **タスクをクリック** | 左パレットからクリックして追加 |
| **⠿ ドラッグ** | タスクの順番を並び替え |
| **✕ ボタン** | タスクを削除 |
| **📂 開く** | 保存済みのパイプライン JSON を読み込む |
| **💾 保存** | 現在のパイプラインを JSON ファイルに保存 |
| **▶ 実行** | パイプラインを実行してログをリアルタイム表示 |

### 3. APIキーを設定する（クラウド AI を使う場合）

```bash
# 環境変数で設定（推奨）
export OPENAI_API_KEY="sk-..."
export GEMINI_API_KEY="AIza..."
```

または、LLM設定パネルの「API キー」欄に直接入力することもできます。

> Ollama を使う場合はAPIキー不要です。[Ollama](https://ollama.com) をインストールしてモデルを pull するだけで動作します。

### 開発サーバー（ホットリロード）

```bash
wails dev
```

---

## パイプラインJSON仕様

### 基本構造

```json
{
  "name": "パイプライン名",
  "version": "1.0",
  "llm": {
    "provider": "ollama",
    "model": "gpt-oss:20b",
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
| `model` | string | モデル名（例: `"gpt-oss:20b"`, `"gpt-4o"`, `"gemini-2.5-flash"`）|
| `base_url` | string | APIエンドポイント（Ollama用。省略時は `http://localhost:11434`）|
| `api_key` | string | APIキー（省略時は環境変数 `OPENAI_API_KEY` / `GEMINI_API_KEY` を使用）|

### 変数テンプレート

タスクのパラメータ内で `{{タスクID}}` と書くと、そのタスクの出力内容に置き換えられます。

```json
{ "prompt": "以下を要約してください:\n\n{{document}}" }
```

---

## タスクタイプ一覧

### 入力系

| タイプ | 説明 |
|--------|------|
| `open_file` | PDF / Word / Excel / テキストファイルを1つ読み込む |
| `open_files` | 複数ファイルをまとめて読み込む |
| `open_folder` | フォルダ内のファイルを全自動で読み込む |
| `read_stdin` | 実行時にキーボードで入力（CLIのみ）|
| `read_stdin_paths` | 実行時にファイルパスを複数入力（CLIのみ）|
| `read_image` | 画像ファイルを読み込む（マルチモーダル用）|
| `read_excel_cells` | Excelの特定セルまたは範囲（A1:C10等）を読み込む |

### AI処理系

| タイプ | 説明 |
|--------|------|
| `llm_prompt` | LLM にプロンプトを送信してレスポンスを取得 |
| `fill_template` | テンプレートファイルにデータを LLM で流し込む |

### 出力系

| タイプ | 説明 |
|--------|------|
| `display_result` | 結果を画面（ログパネル）に表示 |
| `write_file` | テキストをファイルに書き込む |
| `write_excel` | Excelの指定セルに値を書き込む |
| `call_mcp` | MCP サーバーのツールを呼び出す |

---

### `read_excel_cells` — Excelセル読み込み

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| `path` | string | ✓ | Excelファイルパス（`{{変数}}` 使用可）|
| `sheet` | string | - | シート名（省略時: 先頭シート）|
| `range` | string | ✓ | セル番地または範囲（例: `"A1"`, `"A1:D20"`）|

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

詳細な仕様は [format.md](./format.md) を参照してください。

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
| `examples/gemini_fill_report.json` | 議事録・進捗・課題からテンプレートに流し込み月次報告書生成 |
| `examples/gemini_folder_summary.json` | フォルダ内ファイルを自動読み込みして一括要約 |
| `examples/gemini_folder_report.json` | フォルダ内ファイルを自動読み込みしてテンプレートで報告書生成 |

---

## ディレクトリ構成

```
stack-agent/
├── main.go                          # Wails GUIエントリポイント
├── app.go                           # GUIバインディング（ファイルダイアログ・実行）
├── wails.json                       # Wails設定
├── cmd/
│   └── cli/
│       └── main.go                  # CLIエントリポイント
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
│   ├── read_excel_cells.go          # read_excel_cells
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
│   ├── excel.go                     # Excel読み書き・セル読み込み
│   └── image.go                     # 画像読み込み
├── mcp/
│   └── client.go                    # MCPクライアント（stdio）
├── pipeline-editor/                 # GUIフロントエンド（React + TypeScript）
│   ├── src/
│   │   ├── App.tsx                  # メインコンポーネント
│   │   ├── i18n.tsx                 # 日本語/English 切り替え
│   │   ├── taskDefs.ts              # タスクタイプ定義
│   │   ├── types.ts                 # 型定義
│   │   ├── validate.ts              # バリデーション
│   │   └── components/              # UIコンポーネント
│   └── dist/                        # ビルド済みフロントエンド（Goに埋め込み）
└── examples/
    ├── templates/                   # 報告書テンプレート
    ├── sample_data/                 # サンプルデータ
    └── *.json                       # サンプルパイプライン
```
