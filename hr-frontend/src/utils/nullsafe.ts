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
