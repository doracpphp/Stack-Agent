import { useState, createContext, useContext, type ReactNode } from 'react';

export type Locale = 'ja' | 'en';

export interface T {
  // App header
  appName: string;
  pipelineName: string;
  version: string;
  defaultPipelineName: string;
  errorsCount: (n: number) => string;
  openFile: string;
  fixErrorsToRun: string;
  running: string;
  run: string;
  hideLog: string;
  showLog: string;
  emptyTitle: string;
  emptySubtitle: string;
  importFailed: string;
  runError: string;
  errorPrefix: string;
  done: string;

  // JSON preview
  jsonPreview: string;
  copy: string;
  copied: string;
  save: string;
  saveFailed: string;
  fixErrorsBeforeSave: string;
  errorsTitle: (n: number) => string;
  moreErrors: (n: number) => string;

  // Log panel
  logTitle: string;
  logRunning: string;
  logEmpty: string;

  // Validation
  valPipelineName: string;
  valLLMProvider: string;
  valLLMModel: string;
  valTaskIdRequired: string;
  valTaskIdFormat: string;
  valTaskIdDuplicate: (id: string) => string;
  valFieldRequired: (label: string) => string;

  // Task palette / block
  taskBlocks: string;
  clickToAdd: string;
  catInput: string;
  catAi: string;
  catOutput: string;
  catExternal: string;
  taskIdPlaceholder: string;

  // LLM config
  llmProvider: string;
  ollamaLabel: string;
  modelName: string;
  baseUrl: string;
  apiKey: string;
  apiKeyHint: string;
  apiKeyPlaceholder: string;

  // Param field
  addItem: string;
  keyPlaceholder: string;
  valuePlaceholder: string;
}

const JA: T = {
  appName: 'stack-agent',
  pipelineName: 'パイプライン名',
  version: 'v',
  defaultPipelineName: 'マイパイプライン',
  errorsCount: (n) => `⚠ ${n}件のエラー`,
  openFile: '📂 開く',
  fixErrorsToRun: 'エラーを修正してから実行',
  running: '⏳ 実行中...',
  run: '▶ 実行',
  hideLog: '▾ ログ非表示',
  showLog: '▸ ログ表示',
  emptyTitle: '左のパレットからタスクをクリックして追加',
  emptySubtitle: 'ドラッグハンドル（⠿）で並び替え可能です',
  importFailed: 'JSONの読み込みに失敗しました',
  runError: '❌ 実行エラー:',
  errorPrefix: '❌ エラー:',
  done: '✅ 完了',

  jsonPreview: 'JSON プレビュー',
  copy: 'コピー',
  copied: '✓ コピー済み',
  save: '💾 保存',
  saveFailed: '保存に失敗しました',
  fixErrorsBeforeSave: 'エラーを修正してから保存',
  errorsTitle: (n) => `⚠ ${n} 件のエラー`,
  moreErrors: (n) => `他 ${n} 件...`,

  logTitle: 'ログ',
  logRunning: '⏳ 実行ログ',
  logEmpty: '（ログなし）',

  valPipelineName: 'パイプライン名を入力してください',
  valLLMProvider: 'LLMプロバイダーを選択してください',
  valLLMModel: 'モデル名を入力してください',
  valTaskIdRequired: 'タスクIDを入力してください',
  valTaskIdFormat: 'タスクIDは英数字とアンダースコアのみ使用できます',
  valTaskIdDuplicate: (id) => `タスクID "${id}" が重複しています`,
  valFieldRequired: (label) => `「${label}」は必須です`,

  taskBlocks: 'タスクブロック',
  clickToAdd: 'クリックで追加',
  catInput: '入力',
  catAi: 'AI処理',
  catOutput: '出力',
  catExternal: '外部連携',
  taskIdPlaceholder: 'task_id',

  llmProvider: 'プロバイダー',
  ollamaLabel: 'Ollama（ローカル）',
  modelName: 'モデル名',
  baseUrl: 'Base URL',
  apiKey: 'API キー',
  apiKeyHint: '（省略時は環境変数）',
  apiKeyPlaceholder: 'sk-...',

  addItem: '+ 追加',
  keyPlaceholder: 'キー',
  valuePlaceholder: '値',
};

const EN: T = {
  appName: 'stack-agent',
  pipelineName: 'Pipeline Name',
  version: 'v',
  defaultPipelineName: 'My Pipeline',
  errorsCount: (n) => `⚠ ${n} error${n > 1 ? 's' : ''}`,
  openFile: '📂 Open',
  fixErrorsToRun: 'Fix errors before running',
  running: '⏳ Running...',
  run: '▶ Run',
  hideLog: '▾ Hide Log',
  showLog: '▸ Show Log',
  emptyTitle: 'Click a task in the left palette to add it',
  emptySubtitle: 'Drag the ⠿ handle to reorder tasks',
  importFailed: 'Failed to load JSON',
  runError: '❌ Run error:',
  errorPrefix: '❌ Error:',
  done: '✅ Done',

  jsonPreview: 'JSON Preview',
  copy: 'Copy',
  copied: '✓ Copied',
  save: '💾 Save',
  saveFailed: 'Failed to save',
  fixErrorsBeforeSave: 'Fix errors before saving',
  errorsTitle: (n) => `⚠ ${n} error${n > 1 ? 's' : ''}`,
  moreErrors: (n) => `${n} more...`,

  logTitle: 'Log',
  logRunning: '⏳ Running Log',
  logEmpty: '(no output yet)',

  valPipelineName: 'Pipeline name is required',
  valLLMProvider: 'Select an LLM provider',
  valLLMModel: 'Model name is required',
  valTaskIdRequired: 'Task ID is required',
  valTaskIdFormat: 'Task ID must contain only letters, numbers, and underscores',
  valTaskIdDuplicate: (id) => `Duplicate task ID "${id}"`,
  valFieldRequired: (label) => `"${label}" is required`,

  taskBlocks: 'Task Blocks',
  clickToAdd: 'Click to add',
  catInput: 'Input',
  catAi: 'AI',
  catOutput: 'Output',
  catExternal: 'External',
  taskIdPlaceholder: 'task_id',

  llmProvider: 'Provider',
  ollamaLabel: 'Ollama (local)',
  modelName: 'Model Name',
  baseUrl: 'Base URL',
  apiKey: 'API Key',
  apiKeyHint: '(or use env var)',
  apiKeyPlaceholder: 'sk-...',

  addItem: '+ Add',
  keyPlaceholder: 'key',
  valuePlaceholder: 'value',
};

const TRANSLATIONS: Record<Locale, T> = { ja: JA, en: EN };

// ── React context ───────────────────────────────────────────────────────────

interface I18nContextValue {
  t: T;
  locale: Locale;
  setLocale: (l: Locale) => void;
}

export const I18nContext = createContext<I18nContextValue>({
  t: JA,
  locale: 'ja',
  setLocale: () => {},
});

export function useI18n(): I18nContextValue {
  return useContext(I18nContext);
}

export function I18nProvider({ children }: { children: ReactNode }) {
  const [locale, setLocale] = useState<Locale>('ja');
  const t = TRANSLATIONS[locale];
  return (
    <I18nContext.Provider value={{ t, locale, setLocale }}>
      {children}
    </I18nContext.Provider>
  );
}
