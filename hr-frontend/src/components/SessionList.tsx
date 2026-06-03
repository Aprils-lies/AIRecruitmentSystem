import type { SessionInfo } from '../types';
import { formatDateTime } from '../utils/format';

interface SessionListProps {
  sessions: SessionInfo[];
  activeId: string;
  onSelect: (sessionId: string) => void;
  onNewSession: () => void;
  onDelete: (sessionId: string) => void;
}

export function SessionList({ sessions, activeId, onSelect, onNewSession, onDelete }: SessionListProps) {
  return (
    <div className="w-72 bg-gray-50 border-r border-gray-200 flex flex-col">
      <div className="p-4 border-b border-gray-200">
        <button
          onClick={onNewSession}
          className="w-full px-4 py-2 bg-primary-500 text-white rounded-btn hover:bg-primary-600 transition-colors"
        >
          + 新建会话
        </button>
      </div>
      <div className="flex-1 overflow-y-auto">
        {sessions.length === 0 ? (
          <div className="p-4 text-center text-gray-400">
            <p>暂无会话</p>
            <p className="text-sm mt-1">点击上方按钮创建新会话</p>
          </div>
        ) : (
          <ul className="p-2">
            {sessions.map((session) => (
              <li key={session.session_id} className="relative group">
                <button
                  onClick={() => onSelect(session.session_id)}
                  className={`w-full text-left p-3 rounded-lg transition-colors ${
                    activeId === session.session_id
                      ? 'bg-primary-100 text-primary-700'
                      : 'hover:bg-gray-100 text-gray-700'
                  }`}
                >
                  <p className="text-sm font-medium truncate">{session.last_message || '新会话'}</p>
                  <p className="text-xs text-gray-400 mt-1">{formatDateTime(session.created_at)}</p>
                </button>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onDelete(session.session_id);
                  }}
                  className="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-gray-400 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity"
                  title="删除会话"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
