import client from './client';
import type { ApiResponse, UserProfile, UpdateProfileRequest } from '../types';

export function getProfile(): Promise<ApiResponse<UserProfile>> {
  return client.get('/candidate/profile');
}

export function updateProfile(data: UpdateProfileRequest): Promise<ApiResponse<null>> {
  return client.put('/candidate/profile', data);
}