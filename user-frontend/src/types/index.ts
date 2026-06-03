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

export interface UserProfile {
  user_id: number;
  username: string;
  role: string;
  real_name: string;
  phone: string;
  education: string;
  school: string;
  experience: string;
  skills: string;
}

export interface UpdateProfileRequest {
  real_name?: string;
  phone?: string;
  education?: string;
  school?: string;
  experience?: string;
  skills?: string;
}

export type ApplicationStatus = 'pending' | 'reviewed' | 'rejected' | 'accepted';

export interface MyApplication {
  application_id: number;
  position_id: number;
  position_title: string;
  status: ApplicationStatus;
  applied_at: string;
}

export interface ResumeInfo {
  id: number;
  file_name: string;
  file_type: 'pdf' | 'doc' | 'docx';
  file_size: number;
  uploaded_at: string;
}

export interface UploadSignURLData {
  upload_url: string;
  oss_key: string;
  expire_sec: number;
}

export interface ConfirmUploadRequest {
  oss_key: string;
  file_name: string;
  file_type: string;
  file_size: number;
}

export interface DownloadSignURLData {
  download_url: string;
  file_name: string;
  expire_sec: number;
}