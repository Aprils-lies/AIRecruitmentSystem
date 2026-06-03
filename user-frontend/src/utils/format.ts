export function formatDateTime(isoStr: string): string {
  if (!isoStr) return '';
  try {
    const d = new Date(isoStr);
    if (isNaN(d.getTime())) return '';
    const pad = (n: number) => n.toString().padStart(2, '0');
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`;
  } catch {
    return '';
  }
}

export function formatSalary(min: number, max: number): string {
  if (min <= 0 && max <= 0) return '薪资面议';
  if (min <= 0) return `最高 ${max}K`;
  if (max <= 0) return `最低 ${min}K`;
  return `${min}K - ${max}K`;
}

export function formatFileSize(bytes: number): string {
  if (bytes <= 0) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB'];
  let i = 0;
  let size = bytes;
  while (size >= 1024 && i < units.length - 1) { size /= 1024; i++; }
  return `${size.toFixed(i === 0 ? 0 : 1)} ${units[i]}`;
}

export function formatApplicationStatus(status: string): string {
  const map: Record<string, string> = {
    pending: '待审核',
    reviewed: '已查看',
    rejected: '已拒绝',
    accepted: '已通过',
  };
  return map[status] || status;
}

export function getStatusColor(status: string): string {
  const map: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800',
    reviewed: 'bg-blue-100 text-blue-800',
    rejected: 'bg-red-100 text-red-800',
    accepted: 'bg-green-100 text-green-800',
  };
  return map[status] || 'bg-gray-100 text-gray-800';
}