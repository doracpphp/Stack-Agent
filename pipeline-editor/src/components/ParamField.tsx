import { useI18n } from '../i18n';
import type { LocalizedParamDef } from '../taskDefs';

interface Props {
  def: LocalizedParamDef;
  value: unknown;
  onChange: (val: unknown) => void;
  availableIds: string[];
  hasError: boolean;
}

export function ParamField({ def, value, onChange, availableIds, hasError }: Props) {
  const { t } = useI18n();

  const base =
    'w-full rounded border px-2 py-1 text-sm focus:outline-none focus:ring-1 focus:ring-blue-400 ' +
    (hasError ? 'border-red-400 bg-red-50' : 'border-gray-300 bg-white');

  /* ── bool ── */
  if (def.type === 'bool') {
    return (
      <label className="flex items-center gap-2 cursor-pointer select-none">
        <div
          className={`relative w-9 h-5 rounded-full transition-colors ${value ? 'bg-blue-500' : 'bg-gray-300'}`}
          onClick={() => onChange(!value)}
        >
          <span
            className={`absolute top-0.5 w-4 h-4 rounded-full bg-white shadow transition-transform ${value ? 'translate-x-4' : 'translate-x-0.5'}`}
          />
        </div>
        <span className="text-xs text-gray-600">{def.label}</span>
      </label>
    );
  }

  /* ── string[] ── */
  if (def.type === 'string[]') {
    const arr = Array.isArray(value) ? (value as string[]) : [];
    return (
      <div className="space-y-1">
        {arr.map((item, i) => (
          <div key={i} className="flex gap-1">
            <input
              className={base + ' flex-1'}
              value={item}
              placeholder={def.placeholder}
              onChange={(e) => {
                const next = [...arr];
                next[i] = e.target.value;
                onChange(next);
              }}
            />
            <button
              className="px-2 rounded border border-gray-300 text-gray-500 hover:bg-red-50 hover:border-red-300 hover:text-red-500 text-xs"
              onClick={() => onChange(arr.filter((_, j) => j !== i))}
            >
              ✕
            </button>
          </div>
        ))}
        <button
          className="text-xs text-blue-600 hover:text-blue-800 flex items-center gap-0.5"
          onClick={() => onChange([...arr, ''])}
        >
          {t.addItem}
        </button>
      </div>
    );
  }

  /* ── kv-object ── */
  if (def.type === 'kv-object') {
    const obj = (value && typeof value === 'object' && !Array.isArray(value))
      ? (value as Record<string, string>)
      : {};
    const entries = Object.entries(obj);
    return (
      <div className="space-y-1">
        {entries.map(([k, v], i) => (
          <div key={i} className="flex gap-1 items-center">
            <input
              className="w-28 rounded border border-gray-300 bg-white px-2 py-1 text-sm focus:outline-none focus:ring-1 focus:ring-blue-400"
              value={k}
              placeholder={t.keyPlaceholder}
              onChange={(e) => {
                const next: Record<string, string> = {};
                entries.forEach(([ek, ev], j) => { next[j === i ? e.target.value : ek] = ev; });
                onChange(next);
              }}
            />
            <span className="text-gray-400 text-xs">:</span>
            <input
              className="flex-1 rounded border border-gray-300 bg-white px-2 py-1 text-sm focus:outline-none focus:ring-1 focus:ring-blue-400"
              value={v}
              placeholder={t.valuePlaceholder}
              onChange={(e) => onChange({ ...obj, [k]: e.target.value })}
            />
            <button
              className="px-2 rounded border border-gray-300 text-gray-500 hover:bg-red-50 hover:border-red-300 hover:text-red-500 text-xs"
              onClick={() => { const next = { ...obj }; delete next[k]; onChange(next); }}
            >
              ✕
            </button>
          </div>
        ))}
        <button
          className="text-xs text-blue-600 hover:text-blue-800 flex items-center gap-0.5"
          onClick={() => onChange({ ...obj, '': '' })}
        >
          {t.addItem}
        </button>
      </div>
    );
  }

  /* ── textarea ── */
  if (def.type === 'textarea') {
    return (
      <div className="relative">
        <textarea
          className={base + ' min-h-[72px] resize-y font-mono text-xs'}
          value={(value as string) ?? ''}
          placeholder={def.placeholder}
          onChange={(e) => onChange(e.target.value)}
        />
        {availableIds.length > 0 && (
          <div className="flex flex-wrap gap-1 mt-1">
            {availableIds.map((id) => (
              <button
                key={id}
                className="text-xs px-1.5 py-0.5 rounded bg-gray-100 border border-gray-300 text-gray-600 hover:bg-blue-50 hover:border-blue-300 font-mono"
                onClick={() => onChange(((value as string) ?? '') + `{{${id}}}`)}
              >
                {`{{${id}}}`}
              </button>
            ))}
          </div>
        )}
      </div>
    );
  }

  /* ── string (default) ── */
  return (
    <div className="relative">
      <input
        className={base}
        value={(value as string) ?? ''}
        placeholder={def.placeholder}
        onChange={(e) => onChange(e.target.value)}
      />
      {availableIds.length > 0 && def.placeholder?.includes('{{') && (
        <div className="flex flex-wrap gap-1 mt-1">
          {availableIds.map((id) => (
            <button
              key={id}
              className="text-xs px-1.5 py-0.5 rounded bg-gray-100 border border-gray-300 text-gray-600 hover:bg-blue-50 hover:border-blue-300 font-mono"
              onClick={() => onChange(((value as string) ?? '') + `{{${id}}}`)}
            >
              {`{{${id}}}`}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
