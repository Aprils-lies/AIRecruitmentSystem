export function safeString(value: unknown): string {
  if (value === null || value === undefined) return '';
  return String(value);
}

export function safeNumber(value: unknown): number {
  if (value === null || value === undefined) return 0;
  const n = Number(value);
  return Number.isNaN(n) ? 0 : n;
}

export function safeArray<T = unknown>(value: unknown): T[] {
  if (value === null || value === undefined) return [];
  return Array.isArray(value) ? (value as T[]) : [];
}

export function safeObject<T extends Record<string, unknown>>(value: unknown, fallback: T = {} as T): T {
  if (value === null || value === undefined) return fallback;
  return typeof value === 'object' && !Array.isArray(value) ? (value as T) : fallback;
}

export function safeData<T>(response: unknown, fallback: T): T {
  if (response === null || response === undefined) return fallback;
  const r = response as { data?: T };
  return r.data !== null && r.data !== undefined ? r.data : fallback;
}

export function normalizeProfile(raw: unknown) {
  const r = raw as Record<string, unknown>;
  return {
    user_id: safeNumber(r.user_id),
    username: safeString(r.username),
    role: safeString(r.role),
    real_name: safeString(r.real_name),
    phone: safeString(r.phone),
    education: safeString(r.education),
    school: safeString(r.school),
    experience: safeString(r.experience),
    skills: safeString(r.skills),
  };
}

export function normalizePosition(raw: unknown) {
  const r = raw as Record<string, unknown>;
  return {
    id: safeNumber(r.id),
    hr_id: safeNumber(r.hr_id),
    title: safeString(r.title),
    description: safeString(r.description),
    requirements: safeString(r.requirements),
    salary_min: safeNumber(r.salary_min),
    salary_max: safeNumber(r.salary_max),
    location: safeString(r.location),
    status: safeString(r.status) as 'published' | 'offline',
    created_at: safeString(r.created_at),
    updated_at: safeString(r.updated_at),
  };
}

export function normalizeResume(raw: unknown) {
  const r = raw as Record<string, unknown>;
  return {
    id: safeNumber(r.id),
    file_name: safeString(r.file_name),
    file_type: safeString(r.file_type),
    file_size: safeNumber(r.file_size),
    uploaded_at: safeString(r.uploaded_at),
  };
}

export function normalizeApplication(raw: unknown) {
  const r = raw as Record<string, unknown>;
  return {
    application_id: safeNumber(r.application_id),
    position_id: safeNumber(r.position_id),
    position_title: safeString(r.position_title),
    status: safeString(r.status),
    applied_at: safeString(r.applied_at),
  };
}