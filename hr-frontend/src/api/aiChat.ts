import client from './client';
import type { ApiResponse, ChatRequest, HistoryItem, SessionInfo, StatsData } from '../types';
import { getToken } from '../utils/token';

export function chatStream(data: ChatRequest): Promise<Response> {
  const token = getToken();
  return fetch('/api/hr/ai/chat/stream', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify(data),
  });
}

export function getHistory(params: {
  session_id: string;
  limit?: number;
  offset?: number;
}): Promise<ApiResponse<{ items: HistoryItem[]; total: number }>> {
  return client.get('/hr/ai/history', { params });
}

export function listSessions(): Promise<ApiResponse<{ sessions: SessionInfo[] }>> {
  return client.get('/hr/ai/sessions');
}

export function getStats(): Promise<ApiResponse<StatsData>> {
  return client.get('/hr/ai/stats');
}

export function deleteSession(sessionId: string): Promise<ApiResponse<null>> {
  return client.delete(`/hr/ai/sessions/${sessionId}`);
}

export function deleteMessage(messageId: number): Promise<ApiResponse<null>> {
  return client.delete(`/hr/ai/messages/${messageId}`);
}
