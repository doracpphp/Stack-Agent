import { useState } from 'react';
import type { Pipeline } from '../types';
import type { ValidationError } from '../types';
import { SavePipeline } from '../wailsjs/go/main/App';
import { useI18n } from '../i18n';

interface Props {
  pipeline: Pipeline;
  errors: ValidationError[];
}

function toExportJSON(pipeline: Pipeline): object {
  return {
    name: pipeline.name,
    version: pipeline.version,
    llm: pipeline.llm,
    tasks: pipeline.tasks.map(({ _key: _k, ...rest }) => rest),
  };
}

export function JsonPreview({ pipeline, errors }: Props) {
  const [copied, setCopied] = useState(false);
  const { t } = useI18n();

  const json = JSON.stringify(toExportJSON(pipeline), null, 2);

  function handleCopy() {
    navigator.clipboard.writeText(json);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  async function handleSave() {
    try {
      await SavePipeline(json);
    } catch (e) {
      alert(t.saveFailed + ': ' + e);
    }
  }

  return (
    <aside className="w-80 flex-shrink-0 bg-gray-900 flex flex-col">
      <div className="flex items-center justify-between px-3 py-2 border-b border-gray-700">
        <span className="text-xs font-bold text-gray-300 uppercase tracking-wide">{t.jsonPreview}</span>
        <div className="flex gap-1.5">
          <button
            className="text-xs px-2 py-1 rounded bg-gray-700 text-gray-300 hover:bg-gray-600 transition-colors"
            onClick={handleCopy}
          >
            {copied ? t.copied : t.copy}
          </button>
          <button
            className="text-xs px-2 py-1 rounded bg-blue-600 text-white hover:bg-blue-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            onClick={handleSave}
            disabled={errors.length > 0}
            title={errors.length > 0 ? t.fixErrorsBeforeSave : ''}
          >
            {t.save}
          </button>
        </div>
      </div>

      {errors.length > 0 && (
        <div className="bg-red-900/50 border-b border-red-800 px-3 py-2">
          <p className="text-xs font-bold text-red-300 mb-1">{t.errorsTitle(errors.length)}</p>
          {errors.slice(0, 4).map((e, i) => (
            <p key={i} className="text-xs text-red-400">• {e.message}</p>
          ))}
          {errors.length > 4 && (
            <p className="text-xs text-red-400">{t.moreErrors(errors.length - 4)}</p>
          )}
        </div>
      )}

      <pre className="flex-1 overflow-auto p-3 text-xs text-green-300 font-mono leading-relaxed">
        {json}
      </pre>
    </aside>
  );
}

export { toExportJSON };
