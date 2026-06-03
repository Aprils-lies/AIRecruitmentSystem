import { useState, useEffect } from 'react';
import type { PositionInfo, CreatePositionRequest } from '../types';

interface PositionFormProps {
  title: string;
  position?: PositionInfo | null;
  onSubmit: (data: CreatePositionRequest) => void;
  onClose: () => void;
}

export function PositionForm({ title, position, onSubmit, onClose }: PositionFormProps) {
  const [formData, setFormData] = useState<CreatePositionRequest>({
    title: '',
    description: '',
    requirements: '',
    salary_min: 0,
    salary_max: 0,
    location: '',
  });

  useEffect(() => {
    if (position) {
      setFormData({
        title: position.title,
        description: position.description,
        requirements: position.requirements,
        salary_min: position.salary_min,
        salary_max: position.salary_max,
        location: position.location,
      });
    }
  }, [position]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white rounded-card shadow-modal p-6 w-full max-w-lg mx-4 max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">岗位名称 *</label>
            <input
              type="text"
              required
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="请输入岗位名称"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">岗位职责描述 *</label>
            <textarea
              required
              rows={3}
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="请输入岗位职责描述"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">任职要求</label>
            <textarea
              rows={3}
              value={formData.requirements}
              onChange={(e) => setFormData({ ...formData, requirements: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="请输入任职要求"
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">最低薪资（K）</label>
              <input
                type="number"
                min="0"
                value={formData.salary_min}
                onChange={(e) => setFormData({ ...formData, salary_min: parseInt(e.target.value) || 0 })}
                className="w-full px-4 py-2 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">最高薪资（K）</label>
              <input
                type="number"
                min="0"
                value={formData.salary_max}
                onChange={(e) => setFormData({ ...formData, salary_max: parseInt(e.target.value) || 0 })}
                className="w-full px-4 py-2 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500"
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">工作地点</label>
            <input
              type="text"
              value={formData.location}
              onChange={(e) => setFormData({ ...formData, location: e.target.value })}
              className="w-full px-4 py-2 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="请输入工作地点"
            />
          </div>
          <div className="flex justify-end gap-2 mt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-btn transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-primary-500 text-white rounded-btn hover:bg-primary-600 transition-colors"
            >
              保存
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
