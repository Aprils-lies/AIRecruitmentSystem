import { useState } from 'react';
import type { CandidateDetail as CandidateDetailType } from '../types';
import { formatDateTime } from '../utils/format';
import { getDownloadSignURL } from '../api/resume';
import { LoadingButton } from './Loading';

interface CandidateDetailProps {
  candidate: CandidateDetailType;
  onClose: () => void;
  onToast: (message: string, type: 'success' | 'error' | 'warning') => void;
}

export function CandidateDetail({ candidate, onClose, onToast }: CandidateDetailProps) {
  const [downloading, setDownloading] = useState(false);

  const handleDownloadResume = async () => {
    if (candidate.resume_id <= 0) return;
    
    setDownloading(true);
    try {
      const resp = await getDownloadSignURL(candidate.resume_id);
      const { download_url, file_name } = resp.data;
      const link = document.createElement('a');
      link.href = download_url;
      link.download = file_name;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      onToast('简历下载成功', 'success');
    } catch (error) {
      onToast('下载失败，请重试', 'error');
    } finally {
      setDownloading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white rounded-card shadow-modal p-6 w-full max-w-2xl mx-4 max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-semibold text-gray-900">候选人详情</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm text-gray-500 mb-1">姓名</label>
              <p className="font-medium">{candidate.real_name || '-'}</p>
            </div>
            <div>
              <label className="block text-sm text-gray-500 mb-1">联系电话</label>
              <p>{candidate.phone || '-'}</p>
            </div>
            <div>
              <label className="block text-sm text-gray-500 mb-1">学历</label>
              <p>{candidate.education || '-'}</p>
            </div>
            <div>
              <label className="block text-sm text-gray-500 mb-1">毕业院校</label>
              <p>{candidate.school || '-'}</p>
            </div>
          </div>
          
          <div>
            <label className="block text-sm text-gray-500 mb-1">工作/项目经历</label>
            <p className="text-gray-700 whitespace-pre-wrap">{candidate.experience || '-'}</p>
          </div>
          
          <div>
            <label className="block text-sm text-gray-500 mb-1">核心技能</label>
            <div className="flex flex-wrap gap-2">
              {candidate.skills?.split(',').map((skill, index) => (
                <span key={index} className="px-2 py-1 bg-gray-100 text-gray-700 rounded text-sm">
                  {skill.trim()}
                </span>
              ))}
            </div>
          </div>
          
          <div>
            <label className="block text-sm text-gray-500 mb-1">投递时间</label>
            <p>{formatDateTime(candidate.applied_at)}</p>
          </div>
          
          {candidate.resume_id > 0 && (
            <div className="pt-4 border-t border-gray-200">
              <label className="block text-sm text-gray-500 mb-2">简历</label>
              <button
                onClick={handleDownloadResume}
                disabled={downloading}
                className="flex items-center gap-2 px-4 py-2 bg-primary-500 text-white rounded-btn hover:bg-primary-600 transition-colors disabled:opacity-50"
              >
                {downloading ? <LoadingButton /> : (
                  <>
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                    </svg>
                    下载简历 ({candidate.resume_name})
                  </>
                )}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
