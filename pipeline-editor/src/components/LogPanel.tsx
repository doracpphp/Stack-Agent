import { useEffect, useRef } from 'react';
import { useI18n } from '../i18n';

interface Props {
  logs: string[];
  running: boolean;
  onClose: () => void;
}

export function LogPanel({ logs, running, onClose }: Props) {
  const bottomRef = useRef<HTMLDivElement>(null);
  const { t } = useI18n();

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [logs]);

  return (
    <div className="bg-gray-950 border-b border-gray-700 flex flex-col" style={{ height: '220px' }}>
      <div className="flex items-center justify-between px-3 py-1.5 border-b border-gray-800 flex-shrink-0">
        <span className="text-xs font-bold text-gray-400 uppercase tracking-wide">
          {running ? t.logRunning : t.logTitle}
        </span>
        <button className="text-xs text-gray-500 hover:text-gray-300 px-2" onClick={onClose}>
          ✕
        </button>
      </div>
      <div className="flex-1 overflow-y-auto px-3 py-2 font-mono text-xs leading-relaxed">
        {logs.length === 0 ? (
          <span className="text-gray-600">{t.logEmpty}</span>
        ) : (
          logs.map((line, i) => (
            <div
              key={i}
              className={
                line.startsWith('❌') ? 'text-red-400'
                : line.startsWith('✅') ? 'text-green-400'
                : line.startsWith('▶') ? 'text-blue-400'
                : 'text-gray-300'
              }
            >
              {line || '\u00a0'}
            </div>
          ))
        )}
        {running && <div className="text-yellow-400 animate-pulse">▌</div>}
        <div ref={bottomRef} />
      </div>
    </div>
  );
}
