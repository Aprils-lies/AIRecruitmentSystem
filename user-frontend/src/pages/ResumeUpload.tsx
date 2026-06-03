import { useState, useEffect, useCallback, useRef } from 'react';
import { listMyResumes, getUploadSignURL, confirmUpload, getDownloadSignURL, deleteResume } from '../api/resume';
import { normalizeResume } from '../utils/nullsafe';
import { formatFileSize, formatDateTime } from '../utils/format';
import Loading from '../components/Loading';
import EmptyState from '../components/EmptyState';

interface ResumeUploadProps {
  onSuccess: () => void;
}

export default function ResumeUpload({ onSuccess }: ResumeUploadProps) {
  const [resumes, setResumes] = useState<ReturnType<typeof normalizeResume>[]>([]);
  const [loading, setLoading] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState('');
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const fetchResumes = useCallback(async () => {
    setLoading(true);
    try {
      const res = await listMyResumes();
      setResumes(res.data.resumes?.map((r: unknown) => normalizeResume(r)) || []);
    } catch (e) {
      setError(e instanceof Error ? e.message : '获取简历列表失败');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchResumes();
  }, [fetchResumes]);

  const handleFile = async (file: File) => {
    const allowedTypes = ['application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'];
    const fileExtension = file.name.split('.').pop()?.toLowerCase();
    
    if (!allowedTypes.includes(file.type) && !['pdf', 'doc', 'docx'].includes(fileExtension || '')) {
      setError('只支持 PDF、DOC、DOCX 格式的文件');
      return;
    }

    if (file.size > 10 * 1024 * 1024) {
      setError('文件大小不能超过 10MB');
      return;
    }

    setError(null);
    setUploading(true);

    try {
      const signRes = await getUploadSignURL(file.name, file.type);
      const { upload_url, oss_key } = signRes.data;

      const uploadRes = await fetch(upload_url, {
        method: 'PUT',
        headers: {
          'Content-Type': file.type,
        },
        body: file,
      });

      if (!uploadRes.ok) {
        const errText = await uploadRes.text();
        console.error('OSS Upload Error:', errText);
        throw new Error(`OSS上传失败: ${uploadRes.status}`);
      }

      await confirmUpload({
        oss_key,
        file_name: file.name,
        file_type: fileExtension || 'pdf',
        file_size: file.size,
      });

      setSuccessMessage('上传成功');
      setTimeout(() => setSuccessMessage(''), 3000);
      await fetchResumes();
      onSuccess();
    } catch (e) {
      setError(e instanceof Error ? e.message : '上传失败');
    } finally {
      setUploading(false);
      setIsDragging(false);
    }
  };

  const handleFileSelect = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    await handleFile(file);
    e.target.value = '';
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);

    const file = e.dataTransfer.files?.[0];
    if (file) {
      handleFile(file);
    }
  };

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handlePreview = async (resumeId: number) => {
    try {
      const res = await getDownloadSignURL(resumeId);
      const { download_url } = res.data;
      if (download_url) {
        window.open(download_url, '_blank');
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : '获取预览链接失败');
    }
  };

  const handleDelete = async (resumeId: number) => {
    if (!window.confirm('确定要删除这份简历吗？此操作不可撤销。')) {
      return;
    }
    
    try {
      await deleteResume(resumeId);
      setSuccessMessage('删除成功');
      setTimeout(() => setSuccessMessage(''), 3000);
      await fetchResumes();
    } catch (e) {
      setError(e instanceof Error ? e.message : '删除失败');
    }
  };

  if (loading) {
    return <Loading text="加载简历列表..." />;
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">简历管理</h1>

      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">上传简历</h3>
        
        {error && (
          <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg">
            {error}
          </div>
        )}

        {successMessage && (
          <div className="mb-4 p-3 bg-green-100 text-green-700 rounded-lg">
            {successMessage}
          </div>
        )}

        <div
          className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${
            uploading 
              ? 'border-primary-300 bg-primary-50' 
              : isDragging 
                ? 'border-primary-500 bg-primary-50' 
                : 'border-gray-300 hover:border-primary-500 hover:bg-gray-50'
          }`}
          onClick={() => fileInputRef.current?.click()}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onDragEnter={handleDragEnter}
        >
          <input
            ref={fileInputRef}
            type="file"
            accept=".pdf,.doc,.docx"
            onChange={handleFileSelect}
            disabled={uploading}
            className="hidden"
          />
          
          {uploading ? (
            <>
              <svg className="w-12 h-12 text-primary-500 mx-auto mb-3 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <p className="text-gray-600">上传中...</p>
            </>
          ) : (
            <>
              <svg className="w-12 h-12 text-gray-400 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              <p className="text-gray-600 mb-2">
                {isDragging ? '松开以上传文件' : '点击或拖拽文件到此处上传'}
              </p>
              <p className="text-sm text-gray-400">支持 PDF、DOC、DOCX 格式，最大 10MB</p>
            </>
          )}
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">我的简历</h3>
        
        {resumes.length === 0 ? (
          <EmptyState title="暂无简历" description="请上传您的简历以便投递岗位" icon="file" />
        ) : (
          <div className="space-y-4">
            {resumes.map((resume) => (
              <div
                key={resume.id}
                className="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
              >
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center">
                    <svg className="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">{resume.file_name}</p>
                    <p className="text-sm text-gray-500">
                      {formatFileSize(resume.file_size)} · {formatDateTime(resume.uploaded_at)}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => handlePreview(resume.id)}
                    className="px-3 py-1 text-sm text-primary-600 hover:bg-primary-50 rounded-lg transition-colors"
                  >
                    预览
                  </button>
                  <button
                    onClick={() => handleDelete(resume.id)}
                    className="px-3 py-1 text-sm text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                  >
                    删除
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}