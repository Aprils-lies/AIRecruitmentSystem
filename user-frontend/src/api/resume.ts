import client from './client';
import type { ApiResponse, ResumeInfo, UploadSignURLData, ConfirmUploadRequest, DownloadSignURLData } from '../types';

export function getUploadSignURL(fileName: string, contentType: string): Promise<ApiResponse<UploadSignURLData>> {
  return client.get('/candidate/resumes/upload-url', {
    params: { file_name: fileName, content_type: contentType },
  });
}

export function confirmUpload(data: ConfirmUploadRequest): Promise<ApiResponse<{ resume_id: number }>> {
  return client.post('/candidate/resumes/confirm', data);
}

export function listMyResumes(): Promise<ApiResponse<{ resumes: ResumeInfo[] }>> {
  return client.get('/candidate/resumes');
}

export function getDownloadSignURL(resumeId: number): Promise<ApiResponse<DownloadSignURLData>> {
  return client.get(`/candidate/resumes/${resumeId}/download-url`);
}

export function deleteResume(resumeId: number): Promise<ApiResponse<{ success: boolean }>> {
  return client.delete(`/candidate/resumes/${resumeId}`);
}