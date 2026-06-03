export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginData {
  token: string;
  user_id: number;
  role: 'hr' | 'candidate';
  username: string;
}

export interface RegisterRequest {
  username: string;
  password: string;
  role: 'hr' | 'candidate';
}

export interface PositionInfo {
  id: number;
  hr_id: number;
  title: string;
  description: string;
  requirements: string;
  salary_min: number;
  salary_max: number;
  location: string;
  status: 'published' | 'offline';
  created_at: string;
  updated_at: string;
}

export interface PositionListData {
  positions: PositionInfo[];
  total: number;
  page: number;
}

export interface CreatePositionRequest {
  title: string;
  description: string;
  requirements?: string;
  salary_min?: number;
  salary_max?: number;
  location?: string;
}

export type ApplicationStatus = 'pending' | 'reviewed' | 'rejected' | 'accepted';

export interface CandidateBrief {
  user_id: number;
  username: string;
  real_name: string;
  phone: string;
  education: string;
  school: string;
  skills: string;
  resume_id: number;
  resume_name: string;
  application_id: number;
  applied_at: string;
  status: ApplicationStatus;
}

export interface CandidateDetail {
  user_id: number;
  username: string;
  real_name: string;
  phone: string;
  education: string;
  school: string;
  experience: string;
  skills: string;
  resume_id: number;
  resume_name: string;
  resume_type: string;
  applied_at: string;
}

export interface DownloadSignURLData {
  download_url: string;
  file_name: string;
  expire_sec: number;
}

export interface ChatRequest {
  question: string;
  session_id?: string;
}

export interface ChatData {
  answer: string;
  message_id: number;
  session_id: string;
}

export interface StreamChunk {
  chunk: string;
  done: boolean;
  message_id: number;
  session_id: string;
}

export interface HistoryItem {
  id: number;
  session_id: string;
  role: 'user' | 'assistant';
  content: string;
  created_at: string;
}

export interface SessionInfo {
  session_id: string;
  last_message: string;
  created_at: string;
}

export interface StatsData {
  summary: string;
}
