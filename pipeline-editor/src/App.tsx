import { useState, useCallback, useEffect, useRef } from 'react';
import {
  DndContext, closestCenter, KeyboardSensor, PointerSensor,
  useSensor, useSensors, type DragEndEvent,
} from '@dnd-kit/core';
import {
  SortableContext, sortableKeyboardCoordinates,
  verticalListSortingStrategy, arrayMove,
} from '@dnd-kit/sortable';

import type { Pipeline, Task } from './types';
import { getTaskDef } from './taskDefs';
import { validate } from './validate';
import { toExportJSON } from './components/JsonPreview';
import { OpenPipeline, RunPipeline } from './wailsjs/go/main/App';
import { EventsOn } from './wailsjs/runtime/runtime';
import { useI18n, type Locale } from './i18n';

import { TaskPalette } from './components/TaskPalette';
import { TaskBlock } from './components/TaskBlock';
import { LLMConfigPanel } from './components/LLMConfig';
import { JsonPreview } from './components/JsonPreview';
import { LogPanel } from './components/LogPanel';

let keyCounter = 0;
function newKey() { return `task-${++keyCounter}`; }

function makeDefaultPipeline(name: string): Pipeline {
  return {
    name,
    version: '1.0',
    llm: { provider: 'ollama', model: 'gpt-oss:20b', base_url: 'http://localhost:11434', api_key: '' },
    tasks: [],
  };
}

export default function App() {
  const { t, locale, setLocale } = useI18n();
  const [pipeline, setPipeline] = useState<Pipeline>(() => makeDefaultPipeline(t.defaultPipelineName));
  const [logs, setLogs] = useState<string[]>([]);
  const [running, setRunning] = useState(false);
  const [showLogs, setShowLogs] = useState(false);
  const logsRef = useRef<string[]>([]);

  const errors = validate(pipeline, t, locale);

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  );

  useEffect(() => {
    const offLog = EventsOn('pipeline:log', (data) => {
      const line = String(data ?? '');
      logsRef.current = [...logsRef.current, line];
      setLogs([...logsRef.current]);
    });
    const offDone = EventsOn('pipeline:done', (data) => {
      const errMsg = String(data ?? '');
      if (errMsg) {
        logsRef.current = [...logsRef.current, '', `${t.errorPrefix} ${errMsg}`];
      } else {
        logsRef.current = [...logsRef.current, '', t.done];
      }
      setLogs([...logsRef.current]);
      setRunning(false);
    });
    return () => { offLog(); offDone(); };
  }, [t]);

  function updateTask(key: string, updated: Task) {
    setPipeline((p) => ({ ...p, tasks: p.tasks.map((task) => task._key === key ? updated : task) }));
  }
  function deleteTask(key: string) {
    setPipeline((p) => ({ ...p, tasks: p.tasks.filter((task) => task._key !== key) }));
  }
  function addTask(type: string) {
    const def = getTaskDef(type);
    if (!def) return;
    const defaults: Record<string, unknown> = {};
    def.params.forEach((p) => { if (p.defaultValue !== undefined) defaults[p.key] = p.defaultValue; });
    const task: Task = { _key: newKey(), id: type + '_' + (pipeline.tasks.length + 1), type, params: defaults };
    setPipeline((p) => ({ ...p, tasks: [...p.tasks, task] }));
  }

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (over && active.id !== over.id) {
      setPipeline((p) => {
        const from = p.tasks.findIndex((task) => task._key === active.id);
        const to = p.tasks.findIndex((task) => task._key === over.id);
        return { ...p, tasks: arrayMove(p.tasks, from, to) };
      });
    }
  }

  async function handleImport() {
    try {
      const json = await OpenPipeline();
      if (!json) return;
      const data = JSON.parse(json);
      const tasks: Task[] = (data.tasks ?? []).map((task: Omit<Task, '_key'>) => ({
        ...task, _key: newKey(), params: task.params ?? {},
      }));
      setPipeline({ ...makeDefaultPipeline(t.defaultPipelineName), ...data, tasks });
    } catch (e) {
      alert(t.importFailed + ': ' + e);
    }
  }

  async function handleRun() {
    if (running) return;
    const json = JSON.stringify(toExportJSON(pipeline), null, 2);
    logsRef.current = [];
    setLogs([]);
    setRunning(true);
    setShowLogs(true);
    try {
      await RunPipeline(json);
    } catch (e) {
      logsRef.current = [...logsRef.current, `${t.runError} ${e}`];
      setLogs([...logsRef.current]);
      setRunning(false);
    }
  }

  const getAvailableIds = useCallback(
    (taskKey: string) => {
      const idx = pipeline.tasks.findIndex((task) => task._key === taskKey);
      return pipeline.tasks.slice(0, idx).map((task) => task.id).filter(Boolean);
    },
    [pipeline.tasks],
  );

  const LOCALES: { value: Locale; label: string }[] = [
    { value: 'ja', label: '日本語' },
    { value: 'en', label: 'English' },
  ];

  return (
    <div className="flex flex-col h-screen overflow-hidden">
      {/* ── Header ── */}
      <header className="bg-gray-800 text-white px-4 py-2 flex items-center gap-3 flex-shrink-0">
        <span className="font-bold text-sm text-blue-400">⚡ {t.appName}</span>
        <span className="text-gray-600">|</span>
        <span className="text-xs text-gray-400">{t.pipelineName}:</span>
        <input
          className="rounded px-2 py-0.5 text-sm bg-gray-700 border border-gray-600 text-white focus:outline-none focus:border-blue-400 w-48"
          value={pipeline.name}
          onChange={(e) => setPipeline((p) => ({ ...p, name: e.target.value }))}
        />
        <span className="text-xs text-gray-400">{t.version}</span>
        <input
          className="rounded px-2 py-0.5 text-sm bg-gray-700 border border-gray-600 text-white focus:outline-none focus:border-blue-400 w-16"
          value={pipeline.version}
          onChange={(e) => setPipeline((p) => ({ ...p, version: e.target.value }))}
        />

        <div className="flex-1" />

        {errors.length > 0 && (
          <span className="text-xs text-red-400 bg-red-900/40 px-2 py-0.5 rounded">
            {t.errorsCount(errors.length)}
          </span>
        )}

        {/* Language toggle */}
        <div className="flex rounded overflow-hidden border border-gray-600">
          {LOCALES.map(({ value, label }) => (
            <button
              key={value}
              className={`text-xs px-2 py-1 transition-colors ${
                locale === value
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
              }`}
              onClick={() => setLocale(value)}
            >
              {label}
            </button>
          ))}
        </div>

        <button
          className="text-xs px-3 py-1 rounded bg-gray-700 hover:bg-gray-600 text-gray-300"
          onClick={handleImport}
        >
          {t.openFile}
        </button>
        <button
          className={`text-xs px-4 py-1 rounded font-semibold transition-colors ${
            running ? 'bg-yellow-600 text-white cursor-not-allowed'
            : errors.length > 0 ? 'bg-gray-600 text-gray-400 cursor-not-allowed'
            : 'bg-green-600 hover:bg-green-500 text-white'
          }`}
          onClick={handleRun}
          disabled={running || errors.length > 0}
          title={errors.length > 0 ? t.fixErrorsToRun : ''}
        >
          {running ? t.running : t.run}
        </button>
        {(logs.length > 0 || running) && (
          <button
            className="text-xs px-2 py-1 rounded bg-gray-700 hover:bg-gray-600 text-gray-400"
            onClick={() => setShowLogs((v) => !v)}
          >
            {showLogs ? t.hideLog : t.showLog}
          </button>
        )}
      </header>

      <LLMConfigPanel config={pipeline.llm} onChange={(llm) => setPipeline((p) => ({ ...p, llm }))} />

      {showLogs && (
        <LogPanel logs={logs} running={running} onClose={() => setShowLogs(false)} />
      )}

      <div className="flex flex-1 overflow-hidden">
        <TaskPalette onAdd={addTask} />

        <main className="flex-1 overflow-y-auto p-4 bg-gray-100">
          {pipeline.tasks.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-gray-400 gap-3">
              <span className="text-5xl">🧩</span>
              <p className="text-sm font-medium">{t.emptyTitle}</p>
              <p className="text-xs">{t.emptySubtitle}</p>
            </div>
          ) : (
            <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
              <SortableContext items={pipeline.tasks.map((task) => task._key)} strategy={verticalListSortingStrategy}>
                {pipeline.tasks.map((task, idx) => (
                  <TaskBlock
                    key={task._key}
                    task={task}
                    index={idx}
                    availableIds={getAvailableIds(task._key)}
                    errors={errors}
                    onChange={(updated) => updateTask(task._key, updated)}
                    onDelete={() => deleteTask(task._key)}
                  />
                ))}
              </SortableContext>
            </DndContext>
          )}
        </main>

        <JsonPreview pipeline={pipeline} errors={errors} />
      </div>
    </div>
  );
}
