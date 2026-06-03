import client from './client';
import type { ApiResponse, CandidateBrief, CandidateDetail } from '../types';

export function listCandidates(params: {
  positionId: number;
  page?: number;
  page_size?: number;
}): Promise<ApiResponse<{ candidates: CandidateBrief[]; total: number; page: number }>> {
  return client.get(`/hr/positions/${params.positionId}/candidates`, {
    params: { page: params.page, page_size: params.page_size },
  });
}

export function getCandidateDetail(params: {
  candidateId: number;
  applicationId: number;
}): Promise<ApiResponse<CandidateDetail>> {
  return client.get(`/hr/candidates/${params.candidateId}`, {
    params: { application_id: params.applicationId },
  });
}

export function updateApplicationStatus(
  applicationId: number,
  newStatus: string
): Promise<ApiResponse<void>> {
  return client.put(`/hr/applications/${applicationId}/status`, {
    new_status: newStatus,
  });
}
