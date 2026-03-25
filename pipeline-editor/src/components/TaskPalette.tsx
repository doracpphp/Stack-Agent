import { useI18n } from '../i18n';
import { getLocalizedTaskDefs, CATEGORY_STYLE, type Category } from '../taskDefs';

interface Props {
  onAdd: (type: string) => void;
}

const CATEGORY_ORDER: Category[] = ['input', 'ai', 'output', 'external'];

export function TaskPalette({ onAdd }: Props) {
  const { t, locale } = useI18n();

  const catLabels: Record<Category, string> = {
    input:    t.catInput,
    ai:       t.catAi,
    output:   t.catOutput,
    external: t.catExternal,
  };

  const defs = getLocalizedTaskDefs(locale);
  const grouped = CATEGORY_ORDER.map((cat) => ({
    cat,
    defs: defs.filter((d) => d.category === cat),
  }));

  return (
    <aside className="w-56 flex-shrink-0 bg-white border-r border-gray-200 overflow-y-auto">
      <div className="p-3 border-b border-gray-200">
        <h2 className="text-xs font-bold text-gray-500 uppercase tracking-wide">{t.taskBlocks}</h2>
        <p className="text-xs text-gray-400 mt-0.5">{t.clickToAdd}</p>
      </div>

      <div className="p-2 space-y-4">
        {grouped.map(({ cat, defs: catDefs }) => {
          const style = CATEGORY_STYLE[cat];
          return (
            <div key={cat}>
              <div className={`text-xs font-bold px-2 py-1 rounded mb-1 ${style.color} ${style.bg}`}>
                {catLabels[cat]}
              </div>
              <div className="space-y-1">
                {catDefs.map((def) => (
                  <button
                    key={def.type}
                    className={`w-full text-left px-2 py-1.5 rounded border text-xs hover:shadow-sm transition-all active:scale-95 ${style.bg} ${style.border} ${style.color}`}
                    onClick={() => onAdd(def.type)}
                    title={def.description}
                  >
                    <span className="mr-1">{def.icon}</span>
                    {def.label}
                  </button>
                ))}
              </div>
            </div>
          );
        })}
      </div>
    </aside>
  );
}
