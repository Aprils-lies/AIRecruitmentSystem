import { formatDateTime } from '../utils/format';

interface ChatBubbleProps {
  role: 'user' | 'assistant';
  content: string;
  timestamp?: string;
  messageId?: number;
  onDelete?: (messageId: number) => void;
}

export function ChatBubble({ role, content, timestamp, messageId, onDelete }: ChatBubbleProps) {
  const isUser = role === 'user';

  return (
    <div className={`flex gap-3 ${isUser ? 'flex-row-reverse' : ''} group`}>
      <div className={`w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 ${
        isUser ? 'bg-primary-500 text-white' : 'bg-gray-200 text-gray-700'
      }`}>
        {isUser ? (
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        ) : (
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        )}
      </div>
      <div className={`max-w-[70%] ${isUser ? 'text-right' : ''}`}>
        <div className={`inline-block px-4 py-2 rounded-2xl relative ${
          isUser 
            ? 'bg-primary-500 text-white rounded-tr-md' 
            : 'bg-gray-100 text-gray-800 rounded-tl-md'
        }`}>
          <p className="text-sm whitespace-pre-wrap">{content}</p>
          {messageId && onDelete && (
            <button
              onClick={() => onDelete(messageId)}
              className="absolute top-1 right-1 p-1 text-gray-400 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity"
              title="删除消息"
            >
              <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          )}
        </div>
        {timestamp && (
          <p className="text-xs text-gray-400 mt-1">{formatDateTime(timestamp)}</p>
        )}
      </div>
    </div>
  );
}
