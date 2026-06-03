import client from './client';
import type { ApiResponse, PositionListData, PositionInfo } from '../types';

export function listPositions(params: {
  page?: number;
  page_size?: number;
  keyword?: string;
  location?: string;
}): Promise<ApiResponse<PositionListData>> {
  return client.get('/positions', { params });
}

export function getPosition(id: number): Promise<ApiResponse<PositionInfo>> {
  return client.get(`/positions/${id}`);
}