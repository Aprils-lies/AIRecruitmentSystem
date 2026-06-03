import { useState, useEffect, useRef } from 'react';
import type { HistoryItem, SessionInfo } from '../types';
import { chatStream, getHistory, listSessions, deleteSession, deleteMessage } from '../api/aiChat';
import { ChatBubble } from '../components/ChatBubble';
import { SessionList } from '../components/SessionList';
import { Loading } from '../components/Loading';
import { ToastContainer } from '../components/Toast';

interface Toast {
  id: number;
  message: string;
  type: 'success' | 'error' | 'warning';
}

export default function AIChat() {
  const [sessions, setSessions] = useState<SessionInfo[]>([]);
  const [activeSession, setActiveSession] = useState<string>('');
  const [history, setHistory] = useState<HistoryItem[]>([]);
  const [input, setInput] = useState('');
  const [isStreaming, setIsStreaming] = useState(false);
  const [streamingAnswer, setStreamingAnswer] = useState('');
  const [loading, setLoading] = useState(true);
  const [toasts, setToasts] = useState<Toast[]>([]);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const showToast = (message: string, type: 'success' | 'error' | 'warning') => {
    const id = Date.now();
    setToasts(prev => [...prev, { id, message, type }]);
  };

  const removeToast = (id: number) => {
    setToasts(prev => prev.filter(t => t.id !== id));
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [history, streamingAnswer]);

  const loadSessions = async () => {
    try {
      const resp = await listSessions();
      setSessions(resp.data.sessions || []);
      if ((resp.data.sessions || []).length > 0 && !activeSession) {
        setActiveSession(resp.data.sessions![0].session_id);
      } else if (!resp.data.sessions || resp.data.sessions.length === 0) {
        setLoading(false);
      }
    } catch (error) {
      showToast('加载会话列表失败', 'error');
      setLoading(false);
    }
  };

  const loadHistory = async (sessionId: string) => {
    if (!sessionId) return;
    setLoading(true);
    try {
      const resp = await getHistory({ session_id: sessionId, limit: 50 });
      setHistory(resp.data.items);
    } catch (error) {
      showToast('加载对话历史失败', 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadSessions();
  }, []);

  useEffect(() => {
    loadHistory(activeSession);
  }, [activeSession]);

  const handleSelectSession = (sessionId: string) => {
    setActiveSession(sessionId);
    setStreamingAnswer('');
  };

  const handleNewSession = () => {
    setActiveSession('');
    setHistory([]);
    setStreamingAnswer('');
    setLoading(false);
  };

  const handleDeleteSessionInternal = async (sessionId: string) => {
    try {
      await deleteSession(sessionId);
      showToast('删除成功', 'success');
      if (activeSession === sessionId) {
        setActiveSession('');
        setHistory([]);
      }
      loadSessions();
    } catch (error) {
      showToast(error instanceof Error ? error.message : '删除失败', 'error');
    }
  };

  const handleDeleteSession = async (sessionId: string) => {
    if (!confirm('确定要删除这个会话吗？')) return;
    await handleDeleteSessionInternal(sessionId);
  };

  const handleDeleteMessage = async (messageId: number) => {
    if (!confirm('确定要删除这条消息吗？')) return;
    try {
      await deleteMessage(messageId);
      showToast('删除成功', 'success');
      setHistory(prev => {
        const newHistory = prev.filter(item => item.id !== messageId);
        if (newHistory.length === 0 && activeSession) {
          handleDeleteSessionInternal(activeSession);
          return [];
        }
        return newHistory;
      });
    } catch (error) {
      showToast(error instanceof Error ? error.message : '删除失败', 'error');
    }
  };

  const handleSend = async () => {
    if (!input.trim() || isStreaming) return;

    const question = input.trim();
    setInput('');
    
    const newHistory: HistoryItem[] = [...history, {
      id: Date.now(),
      session_id: activeSession,
      role: 'user',
      content: question,
      created_at: new Date().toISOString(),
    }];
    setHistory(newHistory);
    setIsStreaming(true);
    setStreamingAnswer('');

    try {
      const resp = await chatStream({ question, session_id: activeSession || undefined });
      
      if (!resp.ok) {
        const errorText = await resp.text();
        throw new Error(errorText || '请求失败');
      }

      const reader = resp.body?.getReader();
      if (!reader) throw new Error('无法读取响应');

      let fullAnswer = '';
      let buffer = '';
      const decoder = new TextDecoder('utf-8');
      
      while (true) {
        const { done, value } = await reader.read();
        if (done) {
          break;
        }
        
        buffer += decoder.decode(value, { stream: false });
        
        while (true) {
          const lineEndIndex = buffer.indexOf('\n');
          if (lineEndIndex === -1) break;
          
          let line = buffer.substring(0, lineEndIndex);
          buffer = buffer.substring(lineEndIndex + 1);
          
          line = line.replace(/\r$/, '').trim();
          
          if (!line) continue;
          
          if (line.startsWith('data: ')) {
            const dataStr = line.substring(6);
            if (dataStr.trim() === '[DONE]') {
              break;
            }
            try {
              const json = JSON.parse(dataStr);
              if (json.chunk) {
                fullAnswer += json.chunk;
                setStreamingAnswer(fullAnswer);
              }
              if (json.done) {
                setHistory(prev => [...prev, {
                  id: json.message_id || Date.now(),
                  session_id: json.session_id || activeSession || 'new',
                  role: 'assistant',
                  content: fullAnswer,
                  created_at: new Date().toISOString(),
                }]);
                setActiveSession(json.session_id || activeSession || 'new');
                setStreamingAnswer('');
                loadSessions();
                return;
              }
            } catch (e) {
              console.log('JSON parse error:', e, 'data:', dataStr);
              continue;
            }
          }
        }
      }
    } catch (error) {
      showToast(error instanceof Error ? error.message : '发送失败', 'error');
      setHistory(prev => prev.slice(0, -1));
    } finally {
      setIsStreaming(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="h-[calc(100vh-4rem)] flex">
      <SessionList
        sessions={sessions}
        activeId={activeSession}
        onSelect={handleSelectSession}
        onNewSession={handleNewSession}
        onDelete={handleDeleteSession}
      />
      
      <div className="flex-1 flex flex-col">
        {loading ? (
          <div className="flex-1 flex items-center justify-center">
            <Loading />
          </div>
        ) : (
          <>
            <div className="flex-1 overflow-y-auto p-4 space-y-4">
              {(!history || history.length === 0) && !streamingAnswer ? (
                <div className="flex flex-col items-center justify-center h-full text-gray-400">
                  <div className="w-20 h-20 rounded-full bg-gray-100 flex items-center justify-center mb-4">
                    <svg className="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                    </svg>
                  </div>
                  <h3 className="text-lg font-medium text-gray-500">AI 智能助手</h3>
                  <p className="mt-1 text-sm">我可以帮您分析招聘数据、筛选候选人、生成面试问题等</p>
                  <div className="mt-6 text-left bg-gray-50 rounded-lg p-4 max-w-md">
                    <p className="text-xs text-gray-500 mb-2">试试这些问题：</p>
                    <div className="space-y-2">
                      <button
                        onClick={() => { setInput('本周有多少候选人投递？'); }}
                        className="block w-full text-left px-3 py-2 text-sm text-gray-600 hover:bg-white rounded-lg transition-colors"
                      >
                        • 本周有多少候选人投递？
                      </button>
                      <button
                        onClick={() => { setInput('帮我分析一下Java岗位的候选人情况'); }}
                        className="block w-full text-left px-3 py-2 text-sm text-gray-600 hover:bg-white rounded-lg transition-colors"
                      >
                        • 帮我分析一下Java岗位的候选人情况
                      </button>
                      <button
                        onClick={() => { setInput('生成几个前端开发岗位的面试问题'); }}
                        className="block w-full text-left px-3 py-2 text-sm text-gray-600 hover:bg-white rounded-lg transition-colors"
                      >
                        • 生成几个前端开发岗位的面试问题
                      </button>
                    </div>
                  </div>
                </div>
              ) : (
                <>
                  {(history || []).map((item) => (
                    <ChatBubble
                      key={item.id}
                      role={item.role as 'user' | 'assistant'}
                      content={item.content}
                      timestamp={item.created_at}
                      messageId={item.id}
                      onDelete={handleDeleteMessage}
                    />
                  ))}
                  {streamingAnswer && (
                    <ChatBubble
                      role="assistant"
                      content={streamingAnswer}
                    />
                  )}
                  <div ref={messagesEndRef} />
                </>
              )}
            </div>
            
            <div className="p-4 border-t border-gray-200">
              <div className="flex gap-3">
                <div className="flex-1 relative">
                  <textarea
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={handleKeyDown}
                    placeholder="输入您的问题..."
                    disabled={isStreaming}
                    rows={2}
                    className="w-full px-4 py-3 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none disabled:opacity-50"
                  />
                </div>
                <button
                  onClick={handleSend}
                  disabled={!input.trim() || isStreaming}
                  className="px-6 py-3 bg-primary-500 text-white rounded-btn hover:bg-primary-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                >
                  {isStreaming ? (
                    <span className="flex items-center gap-2">
                      <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                      发送中
                    </span>
                  ) : (
                    <>
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                      </svg>
                      发送
                    </>
                  )}
                </button>
              </div>
            </div>
          </>
        )}
      </div>

      <ToastContainer toasts={toasts} onRemove={removeToast} />
    </div>
  );
}
