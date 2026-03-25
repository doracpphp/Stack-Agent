# stack-agent

> **For people who want to use AI automation but find the terminal intimidating — or simply don't know where to start.**

A visual AI pipeline tool where you build workflows by arranging blocks, then run them with a single click.
No coding required. Just chain together tasks like "read a file → send to AI → save the result" and let it run.

---

## Who is this for?

- You already use ChatGPT or Gemini, but you're still doing repetitive tasks manually
- You want to batch-process files or Excel data with AI — but don't know how to automate it
- You don't want to touch a terminal (the "black screen")
- You want to build a distributable tool your team can run by simply double-clicking

---

## Features

- **Double-click to launch**: Single executable with no installation required (.app / .exe)
- **Built-in visual editor**: Build pipelines by clicking and arranging task blocks
- **One-button execution**: Run pipelines and watch logs in real time — without leaving the editor
- **Offline support**: Works fully offline using [Ollama](https://ollama.com) with a local LLM
- **Multi-provider**: Switch between Ollama (local), OpenAI, and Gemini
- **Variable templates**: Pass data between tasks using `{{taskId}}` notation
- **Multimodal support**: Feed images to LLMs (llava / GPT-4o / Gemini)
- **Japanese / English UI**: Toggle the interface language at any time

---

## What you can automate

```
📄 Read Excel cells → 🤖 Analyze with AI → 📊 Write results back to Excel
📁 Read all PDFs in a folder → 🤖 Summarize with AI → 💾 Save to text file
📝 Read meeting notes + progress + issues → 🤖 Generate monthly report → 💾 Fill into a Word template
🖼️ Read an image → 🤖 Describe it with AI → 💬 Display the result
```

---

## Installation & Build

### GUI App (recommended)

Requires [Wails v2](https://wails.io):

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

Output: `build/bin/stack-agent.app` (macOS) / `build/bin/stack-agent.exe` (Windows)

> To distribute: just hand over the `.app` or `.exe` file. Recipients can double-click to run it — no additional installation needed.

### CLI (advanced users)

```bash
go build -o stack-agent ./cmd/cli
./stack-agent <pipeline-json-file>
```

---

## How to use

### 1. Launch the app

Double-click `stack-agent.app` (macOS) or `stack-agent.exe` (Windows).

### 2. Build a pipeline

Click task blocks in the left panel to add them, then fill in their parameters.

| Action | Description |
|--------|-------------|
| **Click a task** | Add it from the left palette |
| **⠿ Drag** | Reorder tasks |
| **✕ button** | Delete a task |
| **📂 Open** | Load a saved pipeline JSON file |
| **💾 Save** | Save the current pipeline to a JSON file |
| **▶ Run** | Execute the pipeline and watch logs in real time |

### 3. Set up API keys (if using cloud AI)

```bash
# Set environment variables (recommended)
export OPENAI_API_KEY="sk-..."
export GEMINI_API_KEY="AIza..."
```

Or enter your API key directly in the LLM config panel inside the app.

> Using Ollama? No API key needed. Just install [Ollama](https://ollama.com) and pull a model — it works entirely offline.

### Development server (hot reload)

```bash
wails dev
```

---

## Pipeline JSON format

### Basic structure

```json
{
  "name": "My Pipeline",
  "version": "1.0",
  "llm": {
    "provider": "ollama",
    "model": "gpt-oss:20b",
    "base_url": "http://localhost:11434",
    "api_key": ""
  },
  "tasks": [
    {
      "id": "taskId",
      "type": "taskType",
      "params": { ... }
    }
  ]
}
```

### LLM config (`llm`)

| Field | Type | Description |
|-------|------|-------------|
| `provider` | string | `"ollama"` / `"openai"` / `"gemini"` |
| `model` | string | Model name (e.g. `"gpt-oss:20b"`, `"gpt-4o"`, `"gemini-2.5-flash"`) |
| `base_url` | string | API endpoint (Ollama only; defaults to `http://localhost:11434`) |
| `api_key` | string | API key (falls back to `OPENAI_API_KEY` / `GEMINI_API_KEY` env vars) |

### Variable templates

Write `{{taskId}}` inside any parameter string to substitute the output of that task.

```json
{ "prompt": "Summarize the following:\n\n{{document}}" }
```

---

## Task types

### Input

| Type | Description |
|------|-------------|
| `open_file` | Read a single PDF / Word / Excel / text file |
| `open_files` | Read multiple files and concatenate them |
| `open_folder` | Auto-read all files in a folder |
| `read_stdin` | Accept keyboard input at runtime (CLI only) |
| `read_stdin_paths` | Accept multiple file paths from keyboard (CLI only) |
| `read_image` | Read an image file (for multimodal LLMs) |
| `read_excel_cells` | Read a specific cell or range (e.g. A1:C10) from an Excel sheet |

### AI processing

| Type | Description |
|------|-------------|
| `llm_prompt` | Send a prompt to the LLM and store the response |
| `fill_template` | Have the LLM fill a template file with provided data |

### Output

| Type | Description |
|------|-------------|
| `display_result` | Display a result in the log panel |
| `write_file` | Write text to a file |
| `write_excel` | Write a value to a specific Excel cell |
| `call_mcp` | Call a tool on an MCP server |

---

### `read_excel_cells` — Read Excel cells

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|:--------:|-------------|
| `path` | string | ✓ | Excel file path (`{{var}}` supported) |
| `sheet` | string | - | Sheet name (default: first sheet) |
| `range` | string | ✓ | Cell address or range (e.g. `"A1"`, `"A1:D20"`) |

- Single cell (`"A1"`): returns the cell's string value
- Range (`"A1:D20"`): returns tab-separated columns, one row per line

```json
{
  "id": "table_data",
  "type": "read_excel_cells",
  "params": {
    "path": "data.xlsx",
    "sheet": "Sales",
    "range": "A1:D20"
  }
}
```

For full specifications, see [format.md](./format.md).

---

## Example pipelines

| File | Description |
|------|-------------|
| `examples/summarize_document.json` | Summarize a file with Ollama and save to txt |
| `examples/excel_analysis.json` | Analyze Excel data with Ollama and write results back |
| `examples/image_describe.json` | Describe an image with llava (multimodal) |
| `examples/mcp_filesystem.json` | Read and summarize a file via MCP filesystem server |
| `examples/gemini_document_qa.json` | File read + Q&A with Gemini |
| `examples/gemini_image_analysis.json` | Detailed image analysis with Gemini, saved to txt |
| `examples/gemini_excel_report.json` | Analyze Excel data with Gemini and write back to Excel |
| `examples/gemini_multi_file_summary.json` | Summarize multiple files with Gemini |
| `examples/gemini_fill_report.json` | Generate monthly report from meeting notes + progress + issues |
| `examples/gemini_folder_summary.json` | Auto-read all files in a folder and summarize |
| `examples/gemini_folder_report.json` | Auto-read folder + fill a report template |

---

## Project structure

```
stack-agent/
├── main.go                          # Wails GUI entry point
├── app.go                           # GUI bindings (file dialogs, pipeline execution)
├── wails.json                       # Wails configuration
├── cmd/
│   └── cli/
│       └── main.go                  # CLI entry point
├── pipeline/
│   └── pipeline.go                  # JSON schema definition & loader
├── context/
│   └── context.go                   # Shared execution context & template engine
├── executor/
│   └── executor.go                  # Sequential task runner & handler registry
├── task/
│   ├── handler.go                   # Handler interface & utilities
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
│   ├── provider.go                  # LLMProvider interface & factory
│   ├── ollama.go                    # Ollama implementation
│   ├── openai.go                    # OpenAI implementation
│   └── gemini.go                    # Gemini implementation
├── office/
│   ├── pdf.go                       # PDF reader
│   ├── word.go                      # Word (.docx) reader
│   ├── excel.go                     # Excel read/write & cell reader
│   └── image.go                     # Image reader
├── mcp/
│   └── client.go                    # MCP client (stdio)
├── pipeline-editor/                 # GUI frontend (React + TypeScript)
│   ├── src/
│   │   ├── App.tsx                  # Main component
│   │   ├── i18n.tsx                 # Japanese / English i18n
│   │   ├── taskDefs.ts              # Task type definitions
│   │   ├── types.ts                 # Type definitions
│   │   ├── validate.ts              # Validation logic
│   │   └── components/              # UI components
│   └── dist/                        # Built frontend (embedded into Go binary)
└── examples/
    ├── templates/                   # Report templates
    ├── sample_data/                 # Sample data files
    └── *.json                       # Sample pipelines
```
