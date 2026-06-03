import type { PositionInfo } from '../types';
import { formatSalary, formatDateTime } from '../utils/format';

interface PositionCardProps {
  position: PositionInfo;
  onClick: () => void;
}

export default function PositionCard({ position, onClick }: PositionCardProps) {
  return (
    <div
      onClick={onClick}
      className="bg-white rounded-lg shadow-md p-5 cursor-pointer hover:shadow-lg transition-shadow border border-gray-100"
    >
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-lg font-semibold text-gray-900 line-clamp-1">
          {position.title || '岗位名称未填写'}
        </h3>
      </div>
      
      <div className="flex items-center gap-4 text-sm text-gray-500 mb-3">
        {position.salary_min > 0 || position.salary_max > 0 ? (
          <span className="text-primary-600 font-medium">
            {formatSalary(position.salary_min, position.salary_max)}
          </span>
        ) : (
          <span className="text-gray-400">薪资面议</span>
        )}
        {position.location && (
          <span className="flex items-center gap-1">
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            {position.location}
          </span>
        )}
      </div>
      
      <p className="text-sm text-gray-600 line-clamp-2 mb-3">
        {position.description || '暂无描述'}
      </p>
      
      <div className="flex items-center justify-between text-xs text-gray-400">
        <span>发布于 {formatDateTime(position.created_at)}</span>
        <span className={`px-2 py-0.5 rounded-full text-xs ${
          position.status === 'published' ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-500'
        }`}>
          {position.status === 'published' ? '发布中' : '已下架'}
        </span>
      </div>
    </div>
  );
}