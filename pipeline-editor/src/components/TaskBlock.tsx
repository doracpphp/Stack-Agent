import { useState } from 'react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import type { Task, ValidationError } from '../types';
import { getLocalizedTaskDef, CATEGORY_STYLE, type Category } from '../taskDefs';
import { useI18n } from '../i18n';
import { ParamField } from './ParamField';

interface Props {
  task: Task;
  index: number;
  availableIds: string[];
  errors: ValidationError[];
  onChange: (updated: Task) => void;
  onDelete: () => void;
}

export function TaskBlock({ task, index, availableIds, errors, onChange, onDelete }: Props) {
  const [open, setOpen] = useState(true);
  const { t, locale } = useI18n();

  const catLabels: Record<Category, string> = {
    input: t.catInput, ai: t.catAi, output: t.catOutput, external: t.catExternal,
  };

  const def = getLocalizedTaskDef(task.type, locale);
  const style = def ? CATEGORY_STYLE[def.category] : CATEGORY_STYLE.input;
  const catLabel = def ? catLabels[def.category] : '';

  const taskErrors = errors.filter((e) => e.taskKey === task._key);
  const hasError = taskErrors.length > 0;

  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({ id: task._key });

  const domStyle = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  function setParam(key: string, val: unknown) {
    onChange({ ...task, params: { ...task.params, [key]: val } });
  }

  return (
    <div
      ref={setNodeRef}
      style={domStyle}
      className={`rounded-lg border-2 bg-white shadow-sm ${hasError ? 'border-red-400' : style.border} mb-3`}
    >
      {/* Header */}
      <div className={`flex items-center gap-2 px-3 py-2 rounded-t-md ${style.bg} select-none`}>
        <button
          className="cursor-grab active:cursor-grabbing text-gray-400 hover:text-gray-600 touch-none"
          {...attributes}
          {...listeners}
        >
          ⠿
        </button>

        <span className="text-base">{def?.icon ?? '⚙️'}</span>

        <span className={`text-xs font-bold px-1.5 py-0.5 rounded text-white ${style.badge}`}>
          {catLabel}
        </span>

        <span className={`font-semibold text-sm ${style.color}`}>
          {def?.label ?? task.type}
        </span>

        <span className="text-gray-400 text-xs ml-1">#{index + 1}</span>

        <div className="flex-1" />

        {/* ID field */}
        <div className="flex items-center gap-1">
          <span className="text-xs text-gray-500">id:</span>
          <input
            className={`w-28 rounded border px-1.5 py-0.5 text-xs font-mono focus:outline-none focus:ring-1 focus:ring-blue-400 ${
              taskErrors.some((e) => e.field === 'id') ? 'border-red-400 bg-red-50' : 'border-gray-300 bg-white'
            }`}
            value={task.id}
            placeholder={t.taskIdPlaceholder}
            onClick={(e) => e.stopPropagation()}
            onChange={(e) => onChange({ ...task, id: e.target.value })}
          />
        </div>

        <button className="text-gray-400 hover:text-gray-600 text-sm w-5" onClick={() => setOpen(!open)}>
          {open ? '▲' : '▼'}
        </button>
        <button className="text-gray-400 hover:text-red-500 text-sm w-5" onClick={onDelete}>
          ✕
        </button>
      </div>

      {/* Body */}
      {open && def && (
        <div className="px-4 py-3 space-y-3">
          {def.params.map((param) => {
            const fieldError = taskErrors.some((e) => e.field === param.key);
            const val = task.params[param.key];
            const displayVal = val !== undefined ? val : param.defaultValue;

            return (
              <div key={param.key}>
                <div className="flex items-center gap-1 mb-1">
                  <label className="text-xs font-medium text-gray-600">
                    {param.label}
                    {param.required && <span className="text-red-500 ml-0.5">*</span>}
                  </label>
                  {param.hint && (
                    <span className="text-xs text-gray-400">— {param.hint}</span>
                  )}
                </div>
                <ParamField
                  def={param}
                  value={displayVal}
                  onChange={(v) => setParam(param.key, v)}
                  availableIds={availableIds}
                  hasError={fieldError}
                />
                {fieldError && (
                  <p className="text-xs text-red-500 mt-0.5">
                    {taskErrors.find((e) => e.field === param.key)?.message}
                  </p>
                )}
              </div>
            );
          })}

          {taskErrors.filter((e) => e.field === 'id').map((e, i) => (
            <p key={i} className="text-xs text-red-500">{e.message}</p>
          ))}
        </div>
      )}
    </div>
  );
}
