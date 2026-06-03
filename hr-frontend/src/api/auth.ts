import client from './client';
import type { ApiResponse, LoginRequest, LoginData, RegisterRequest } from '../types';

export function login(data: LoginRequest): Promise<ApiResponse<LoginData>> {
  return client.post('/auth/login', data);
}

export function register(data: RegisterRequest): Promise<ApiResponse<{ user_id: number }>> {
  return client.post('/auth/register', data);
}
