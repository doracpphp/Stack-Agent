import type { LLMConfig, LLMProvider } from '../types';
import { useI18n } from '../i18n';

interface Props {
  config: LLMConfig;
  onChange: (c: LLMConfig) => void;
}

const PROVIDERS: { value: LLMProvider; labelKey: 'ollamaLabel' | null; label?: string }[] = [
  { value: 'ollama', labelKey: 'ollamaLabel' },
  { value: 'openai', labelKey: null, label: 'OpenAI' },
  { value: 'gemini', labelKey: null, label: 'Gemini' },
];

const MODEL_SUGGESTIONS: Record<LLMProvider, string[]> = {
  ollama:  ['gpt-oss:20b', 'llama3.2', 'llava', 'mistral', 'gemma3'],
  openai:  ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo'],
  gemini:  ['gemini-2.5-flash', 'gemini-2.0-flash', 'gemini-1.5-pro'],
};

const inputCls = 'w-full rounded border border-gray-300 px-2 py-1 text-sm focus:outline-none focus:ring-1 focus:ring-blue-400';

export function LLMConfigPanel({ config, onChange }: Props) {
  const { t } = useI18n();

  return (
    <div className="bg-white border-b border-gray-200 px-4 py-3">
      <div className="flex flex-wrap items-end gap-3">
        <div>
          <label className="block text-xs font-medium text-gray-500 mb-1">{t.llmProvider}</label>
          <select
            className={inputCls + ' w-44'}
            value={config.provider}
            onChange={(e) => onChange({ ...config, provider: e.target.value as LLMProvider })}
          >
            {PROVIDERS.map((p) => (
              <option key={p.value} value={p.value}>
                {p.labelKey ? t[p.labelKey] : p.label}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-xs font-medium text-gray-500 mb-1">{t.modelName}</label>
          <input
            list="model-suggestions"
            className={inputCls + ' w-48'}
            value={config.model}
            placeholder={MODEL_SUGGESTIONS[config.provider][0]}
            onChange={(e) => onChange({ ...config, model: e.target.value })}
          />
          <datalist id="model-suggestions">
            {MODEL_SUGGESTIONS[config.provider].map((m) => (
              <option key={m} value={m} />
            ))}
          </datalist>
        </div>

        {config.provider === 'ollama' && (
          <div>
            <label className="block text-xs font-medium text-gray-500 mb-1">{t.baseUrl}</label>
            <input
              className={inputCls + ' w-52'}
              value={config.base_url}
              placeholder="http://localhost:11434"
              onChange={(e) => onChange({ ...config, base_url: e.target.value })}
            />
          </div>
        )}

        {config.provider !== 'ollama' && (
          <div>
            <label className="block text-xs font-medium text-gray-500 mb-1">
              {t.apiKey}
              <span className="text-gray-400 font-normal ml-1">{t.apiKeyHint}</span>
            </label>
            <input
              type="password"
              className={inputCls + ' w-56'}
              value={config.api_key}
              placeholder={t.apiKeyPlaceholder}
              onChange={(e) => onChange({ ...config, api_key: e.target.value })}
            />
          </div>
        )}
      </div>
    </div>
  );
}
