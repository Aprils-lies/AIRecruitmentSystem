import client from './client';
import type { ApiResponse, MyApplication } from '../types';

export function apply(positionId: number, resumeId: number): Promise<ApiResponse<{ application_id: number }>> {
  return client.post(`/candidate/positions/${positionId}/apply`, { resume_id: resumeId });
}

export function listMyApplications(params: {
  page?: number;
  page_size?: number;
}): Promise<ApiResponse<{ applications: MyApplication[]; total: number; page: number }>> {
  return client.get('/candidate/applications', { params });
}

export function withdrawApplication(applicationId: number): Promise<ApiResponse<void>> {
  return client.delete(`/candidate/applications/${applicationId}`);
}