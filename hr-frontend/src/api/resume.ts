import client from './client';
import type { ApiResponse, DownloadSignURLData } from '../types';

export function getDownloadSignURL(resumeId: number): Promise<ApiResponse<DownloadSignURLData>> {
  return client.get(`/hr/resumes/${resumeId}/download-url`);
}
