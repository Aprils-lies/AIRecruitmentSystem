import client from './client';
import type { ApiResponse, PositionListData, PositionInfo, CreatePositionRequest } from '../types';

export function listMyPositions(params: {
  page?: number;
  page_size?: number;
}): Promise<ApiResponse<PositionListData>> {
  return client.get('/hr/my-positions', { params });
}

export function createPosition(data: CreatePositionRequest): Promise<ApiResponse<{ position_id: number }>> {
  return client.post('/hr/positions', data);
}

export function updatePosition(id: number, data: Partial<CreatePositionRequest>): Promise<ApiResponse<null>> {
  return client.put(`/hr/positions/${id}`, data);
}

export function offlinePosition(id: number): Promise<ApiResponse<null>> {
  return client.post(`/hr/positions/${id}/offline`);
}

export function onlinePosition(id: number): Promise<ApiResponse<null>> {
  return client.post(`/hr/positions/${id}/online`);
}

export function getPosition(id: number): Promise<ApiResponse<PositionInfo>> {
  return client.get(`/positions/${id}`);
}
