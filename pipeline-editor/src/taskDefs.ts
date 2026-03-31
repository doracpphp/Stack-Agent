import type { Locale } from './i18n';

export type Category = 'input' | 'ai' | 'output' | 'external';
export type ParamType = 'string' | 'textarea' | 'bool' | 'string[]' | 'kv-object';

export interface ParamDef {
  key: string;
  type: ParamType;
  required: boolean;
  label: string;       // ja
  labelEn: string;     // en
  placeholder?: string;
  defaultValue?: unknown;
  hint?: string;       // ja
  hintEn?: string;     // en
}

export interface TaskTypeDef {
  type: string;
  label: string;       // ja
  labelEn: string;     // en
  category: Category;
  description: string;    // ja
  descriptionEn: string;  // en
  icon: string;
  params: ParamDef[];
}

// Returns the label/hint/description in the given locale
export function localize<T extends { label: string; labelEn: string }>(
  def: T,
  locale: Locale,
): T & { label: string } {
  return locale === 'en' ? { ...def, label: def.labelEn } : def;
}

export interface LocalizedParamDef extends Omit<ParamDef, 'labelEn' | 'hintEn'> {
  hint?: string;
}
export interface LocalizedTaskDef extends Omit<TaskTypeDef, 'labelEn' | 'descriptionEn' | 'params'> {
  description: string;
  params: LocalizedParamDef[];
}

export function getLocalizedTaskDef(type: string, locale: Locale): LocalizedTaskDef | undefined {
  const def = TASK_DEFS.find((d) => d.type === type);
  if (!def) return undefined;
  return {
    ...def,
    label: locale === 'en' ? def.labelEn : def.label,
    description: locale === 'en' ? def.descriptionEn : def.description,
    params: def.params.map((p) => ({
      ...p,
      label: locale === 'en' ? p.labelEn : p.label,
      hint: locale === 'en' ? (p.hintEn ?? p.hint) : p.hint,
    })),
  };
}

export function getLocalizedTaskDefs(locale: Locale): LocalizedTaskDef[] {
  return TASK_DEFS.map((def) => getLocalizedTaskDef(def.type, locale)!);
}

// ── Category metadata ──────────────────────────────────────────────────────
// Labels are injected at runtime from i18n; only visual tokens here.
export const CATEGORY_STYLE: Record<Category, { color: string; bg: string; border: string; badge: string }> = {
  input:    { color: 'text-blue-700',   bg: 'bg-blue-50',   border: 'border-blue-300',  badge: 'bg-blue-500'   },
  ai:       { color: 'text-purple-700', bg: 'bg-purple-50', border: 'border-purple-300',badge: 'bg-purple-500' },
  output:   { color: 'text-green-700',  bg: 'bg-green-50',  border: 'border-green-300', badge: 'bg-green-500'  },
  external: { color: 'text-orange-700', bg: 'bg-orange-50', border: 'border-orange-300',badge: 'bg-orange-500' },
};

// Keep backward compat for components that haven't been migrated yet
export const CATEGORY_META: Record<Category, { label: string; color: string; bg: string; border: string; badge: string }> = {
  input:    { label: '入力',     ...CATEGORY_STYLE.input    },
  ai:       { label: 'AI処理',   ...CATEGORY_STYLE.ai       },
  output:   { label: '出力',     ...CATEGORY_STYLE.output   },
  external: { label: '外部連携', ...CATEGORY_STYLE.external },
};

// ── Task definitions ───────────────────────────────────────────────────────
export const TASK_DEFS: TaskTypeDef[] = [
  {
    type: 'open_file',
    label: 'ファイル読み込み',     labelEn: 'Open File',
    category: 'input',
    description: 'PDF・Word・Excel・テキストファイルを読み込む',
    descriptionEn: 'Read PDF, Word, Excel, or text files',
    icon: '📄',
    params: [
      { key: 'path', type: 'string', required: true,
        label: 'ファイルパス',       labelEn: 'File Path',
        placeholder: 'report.pdf',
        hint: '{{変数}} 使用可',     hintEn: '{{var}} supported' },
    ],
  },
  {
    type: 'open_files',
    label: '複数ファイル読み込み',  labelEn: 'Open Multiple Files',
    category: 'input',
    description: '複数ファイルを連結して読み込む',
    descriptionEn: 'Read and concatenate multiple files',
    icon: '📂',
    params: [
      { key: 'paths', type: 'string[]', required: true,
        label: 'ファイルパス一覧',   labelEn: 'File Paths',
        placeholder: 'file.txt',
        hint: '各パスに {{変数}} 使用可', hintEn: '{{var}} supported in each path' },
      { key: 'separator', type: 'string', required: false,
        label: '区切り文字',         labelEn: 'Separator',
        placeholder: '\\n\\n---\\n\\n' },
    ],
  },
  {
    type: 'open_folder',
    label: 'フォルダ読み込み',      labelEn: 'Open Folder',
    category: 'input',
    description: 'フォルダ内の全ファイルを自動で読み込む',
    descriptionEn: 'Auto-read all files in a folder (no user input needed)',
    icon: '🗂️',
    params: [
      { key: 'path', type: 'string', required: true,
        label: 'フォルダパス',       labelEn: 'Folder Path',
        placeholder: 'input/reports' },
      { key: 'extensions', type: 'string[]', required: false,
        label: '対象拡張子（省略時: 全形式）', labelEn: 'Extensions (leave empty for all)',
        placeholder: '.pdf',
        hint: '.pdf .docx .txt .xlsx など', hintEn: 'e.g. .pdf .docx .txt .xlsx' },
      { key: 'recursive', type: 'bool', required: false,
        label: 'サブフォルダも再帰的に読み込む', labelEn: 'Include subfolders recursively',
        defaultValue: false },
      { key: 'separator', type: 'string', required: false,
        label: '区切り文字',         labelEn: 'Separator' },
    ],
  },
  {
    type: 'read_stdin',
    label: 'キーボード入力',        labelEn: 'Keyboard Input',
    category: 'input',
    description: '実行時にユーザーから入力を受け取る',
    descriptionEn: 'Prompt the user for input at runtime',
    icon: '⌨️',
    params: [
      { key: 'prompt', type: 'string', required: false,
        label: 'プロンプトメッセージ', labelEn: 'Prompt Message',
        placeholder: 'ファイルパスを入力してください: ' },
      { key: 'read_file', type: 'bool', required: false,
        label: '入力をファイルとして読み込む', labelEn: 'Treat input as file path and read it',
        defaultValue: true },
    ],
  },
  {
    type: 'read_stdin_paths',
    label: '複数パス入力',          labelEn: 'Multiple Path Input',
    category: 'input',
    description: '複数ファイルパスを順番に入力（空行で終了）',
    descriptionEn: 'Enter multiple file paths one by one (empty line to finish)',
    icon: '📋',
    params: [
      { key: 'prompt', type: 'string', required: false,
        label: 'プロンプトメッセージ', labelEn: 'Prompt Message',
        placeholder: 'ファイルパスを入力してください（空行で終了）:' },
      { key: 'separator', type: 'string', required: false,
        label: '区切り文字',         labelEn: 'Separator' },
    ],
  },
  {
    type: 'read_image',
    label: '画像読み込み',          labelEn: 'Read Image',
    category: 'input',
    description: '画像ファイルを読み込む（マルチモーダルLLM用）',
    descriptionEn: 'Read an image file for multimodal LLM',
    icon: '🖼️',
    params: [
      { key: 'path', type: 'string', required: true,
        label: '画像ファイルパス',   labelEn: 'Image Path',
        placeholder: 'photo.png',
        hint: '.jpg .png .gif .webp 対応', hintEn: 'Supports .jpg .png .gif .webp' },
    ],
  },
  {
    type: 'llm_prompt',
    label: 'LLMプロンプト',         labelEn: 'LLM Prompt',
    category: 'ai',
    description: 'LLMにプロンプトを送信してレスポンスを取得',
    descriptionEn: 'Send a prompt to the LLM and get a response',
    icon: '🤖',
    params: [
      { key: 'prompt', type: 'textarea', required: true,
        label: 'プロンプト',         labelEn: 'Prompt',
        placeholder: '以下を要約してください:\n\n{{document}}',
        hint: '{{タスクID}} で他タスクの出力を参照', hintEn: 'Use {{taskId}} to reference prior task output' },
      { key: 'system', type: 'textarea', required: false,
        label: 'システムプロンプト', labelEn: 'System Prompt',
        placeholder: 'あなたは優秀なアシスタントです。' },
      { key: 'images', type: 'string[]', required: false,
        label: '画像タスクID（read_imageのID）', labelEn: 'Image Task IDs (from read_image)',
        placeholder: 'photo' },
    ],
  },
  {
    type: 'fill_template',
    label: 'テンプレート流し込み',   labelEn: 'Fill Template',
    category: 'ai',
    description: 'テンプレートファイルにデータをLLMで流し込む',
    descriptionEn: 'Use LLM to fill a template file with provided data',
    icon: '📝',
    params: [
      { key: 'template', type: 'string', required: true,
        label: 'テンプレートファイルパス', labelEn: 'Template File Path',
        placeholder: 'templates/report.txt',
        hint: '{{変数}} 使用可', hintEn: '{{var}} supported' },
      { key: 'data', type: 'string[]', required: true,
        label: '参照タスクID一覧',   labelEn: 'Source Task IDs',
        placeholder: 'document',
        hint: '参照したいタスクのIDを列挙', hintEn: 'List the task IDs to use as data' },
      { key: 'instruction', type: 'textarea', required: false,
        label: 'LLMへの追加指示',   labelEn: 'Additional Instruction',
        placeholder: 'テンプレートの[...]を提供データで埋めてください。' },
    ],
  },
  {
    type: 'display_result',
    label: '結果表示',              labelEn: 'Display Result',
    category: 'output',
    description: '結果をターミナルに表示する',
    descriptionEn: 'Print the result to the terminal',
    icon: '💬',
    params: [
      { key: 'value', type: 'string', required: true,
        label: '表示する値',         labelEn: 'Value to Display',
        placeholder: '{{summary}}',
        hint: '{{タスクID}} を指定', hintEn: 'Specify {{taskId}}' },
      { key: 'label', type: 'string', required: false,
        label: '見出しラベル',       labelEn: 'Header Label',
        placeholder: '要約結果' },
    ],
  },
  {
    type: 'write_file',
    label: 'ファイル書き込み',      labelEn: 'Write File',
    category: 'output',
    description: 'テキストをファイルに書き込む',
    descriptionEn: 'Write text content to a file',
    icon: '💾',
    params: [
      { key: 'path', type: 'string', required: true,
        label: '出力ファイルパス',   labelEn: 'Output File Path',
        placeholder: 'output.txt',
        hint: '{{変数}} 使用可', hintEn: '{{var}} supported' },
      { key: 'value', type: 'string', required: true,
        label: '書き込む内容',       labelEn: 'Content to Write',
        placeholder: '{{result}}' },
      { key: 'append', type: 'bool', required: false,
        label: '追記モード（上書きしない）', labelEn: 'Append mode (do not overwrite)',
        defaultValue: false },
    ],
  },
  {
    type: 'get_current_time',
    label: '現在時刻取得',           labelEn: 'Get Current Time',
    category: 'input',
    description: '現在の日時を取得する',
    descriptionEn: 'Get the current date and time',
    icon: '🕐',
    params: [
      { key: 'format', type: 'string', required: false,
        label: '日時フォーマット',   labelEn: 'Date/Time Format',
        placeholder: '2006-01-02 15:04:05',
        hint: 'Go形式。省略時: 2006-01-02 15:04:05', hintEn: 'Go time format. Default: 2006-01-02 15:04:05' },
      { key: 'timezone', type: 'string', required: false,
        label: 'タイムゾーン',       labelEn: 'Timezone',
        placeholder: 'Asia/Tokyo',
        hint: 'IANA形式。省略時: ローカル', hintEn: 'IANA name. Default: local' },
    ],
  },
  {
    type: 'read_excel_cells',
    label: 'Excelセル読み込み',     labelEn: 'Read Excel Cells',
    category: 'input',
    description: 'Excelの特定セルまたは範囲を読み込む',
    descriptionEn: 'Read a specific cell or range from an Excel sheet',
    icon: '📋',
    params: [
      { key: 'path', type: 'string', required: true,
        label: 'Excelファイルパス',  labelEn: 'Excel File Path',
        placeholder: 'data.xlsx',
        hint: '{{変数}} 使用可', hintEn: '{{var}} supported' },
      { key: 'sheet', type: 'string', required: false,
        label: 'シート名（省略時: 先頭シート）', labelEn: 'Sheet Name (default: first sheet)',
        placeholder: 'Sheet1' },
      { key: 'range', type: 'string', required: true,
        label: 'セルまたは範囲',     labelEn: 'Cell or Range',
        placeholder: 'A1  または  A1:C10',
        hint: '単一セル: A1 / 範囲: A1:C10', hintEn: 'Single: A1 / Range: A1:C10' },
    ],
  },
  {
    type: 'write_excel',
    label: 'Excel書き込み',         labelEn: 'Write Excel',
    category: 'output',
    description: 'ExcelセルにLLM結果を書き込む',
    descriptionEn: 'Write LLM result to a specific Excel cell',
    icon: '📊',
    params: [
      { key: 'path', type: 'string', required: true,
        label: 'Excelファイルパス',  labelEn: 'Excel File Path',
        placeholder: 'report.xlsx' },
      { key: 'sheet', type: 'string', required: false,
        label: 'シート名',           labelEn: 'Sheet Name',
        placeholder: 'Sheet1' },
      { key: 'cell', type: 'string', required: true,
        label: 'セル番地',           labelEn: 'Cell Address',
        placeholder: 'A1' },
      { key: 'value', type: 'string', required: true,
        label: '書き込む値',         labelEn: 'Value to Write',
        placeholder: '{{result}}' },
    ],
  },
  {
    type: 'shell',
    label: 'シェルコマンド実行',     labelEn: 'Run Shell Command',
    category: 'external',
    description: 'OSコマンドを実行して結果を取得する',
    descriptionEn: 'Execute a shell command and capture its output',
    icon: '💻',
    params: [
      { key: 'command', type: 'string', required: true,
        label: 'コマンド',           labelEn: 'Command',
        placeholder: 'ls -la {{workdir}}',
        hint: '{{変数}} 使用可', hintEn: '{{var}} supported' },
      { key: 'workdir', type: 'string', required: false,
        label: '作業ディレクトリ',   labelEn: 'Working Directory',
        placeholder: '/path/to/dir' },
      { key: 'timeout', type: 'string', required: false,
        label: 'タイムアウト（秒）', labelEn: 'Timeout (seconds)',
        placeholder: '60' },
    ],
  },
  {
    type: 'llm_shell',
    label: 'LLMコマンド生成・実行',  labelEn: 'LLM Generate & Run Command',
    category: 'external',
    description: 'LLMにコマンドを生成させて実行する（Claude Code的な動作）',
    descriptionEn: 'Ask the LLM to generate a shell command, then execute it',
    icon: '🤖💻',
    params: [
      { key: 'prompt', type: 'textarea', required: true,
        label: 'LLMへの指示',        labelEn: 'Instruction to LLM',
        placeholder: '{{folder}} 内の .log ファイルを全て削除するコマンドを生成してください',
        hint: '{{変数}} 使用可。LLMはコマンドのみ返す', hintEn: '{{var}} supported. LLM returns only the command' },
      { key: 'system', type: 'textarea', required: false,
        label: '追加システムプロンプト', labelEn: 'Additional System Prompt',
        placeholder: 'macOS環境です。GNU coreutils は使えません。' },
      { key: 'workdir', type: 'string', required: false,
        label: '作業ディレクトリ',   labelEn: 'Working Directory',
        placeholder: '/path/to/dir' },
      { key: 'timeout', type: 'string', required: false,
        label: 'タイムアウト（秒）', labelEn: 'Timeout (seconds)',
        placeholder: '60' },
    ],
  },
  {
    type: 'call_mcp',
    label: 'MCPツール呼び出し',     labelEn: 'Call MCP Tool',
    category: 'external',
    description: 'MCPサーバーのツールを呼び出す',
    descriptionEn: 'Launch an MCP server and call a tool',
    icon: '🔌',
    params: [
      { key: 'command', type: 'string', required: true,
        label: 'コマンド',           labelEn: 'Command',
        placeholder: 'npx' },
      { key: 'args', type: 'string[]', required: false,
        label: 'コマンド引数',       labelEn: 'Command Args',
        placeholder: '-y @modelcontextprotocol/server-filesystem /tmp' },
      { key: 'tool', type: 'string', required: true,
        label: 'ツール名',           labelEn: 'Tool Name',
        placeholder: 'read_file' },
      { key: 'arguments', type: 'kv-object', required: false,
        label: 'ツール引数',         labelEn: 'Tool Arguments',
        hint: '値に {{変数}} 使用可', hintEn: '{{var}} supported in values' },
    ],
  },
];

export function getTaskDef(type: string): TaskTypeDef | undefined {
  return TASK_DEFS.find((d) => d.type === type);
}
