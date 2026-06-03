import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { getPosition } from '../api/position';
import { listMyResumes } from '../api/resume';
import { apply, listMyApplications } from '../api/application';
import { normalizePosition, normalizeResume } from '../utils/nullsafe';
import { formatSalary, formatDateTime } from '../utils/format';
import Loading from '../components/Loading';
import EmptyState from '../components/EmptyState';
import ConfirmDialog from '../components/ConfirmDialog';

interface PositionDetailProps {
  isLoggedIn: boolean;
  isProfileComplete: boolean;
}

export default function PositionDetail({ isLoggedIn, isProfileComplete }: PositionDetailProps) {
  const { id } = useParams<{ id: string }>();
  const [position, setPosition] = useState<ReturnType<typeof normalizePosition> | null>(null);
  const [resumes, setResumes] = useState<ReturnType<typeof normalizeResume>[]>([]);
  const [selectedResume, setSelectedResume] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | undefined>(undefined);
  const [applying, setApplying] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [applyError, setApplyError] = useState<string | null>(null);
  const [hasApplied, setHasApplied] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    
    const fetchData = async () => {
      try {
        const [posRes, resumeRes, appRes] = await Promise.all([
          getPosition(parseInt(id)),
          isLoggedIn ? listMyResumes() : Promise.resolve({ data: { resumes: [] } }),
          isLoggedIn ? listMyApplications({ page: 1, page_size: 100 }) : Promise.resolve({ data: { applications: [] } }),
        ]);
        
        setPosition(normalizePosition(posRes.data));
        setResumes(resumeRes.data.resumes?.map((r: unknown) => normalizeResume(r)) || []);
        
        if (resumeRes.data.resumes && resumeRes.data.resumes.length > 0) {
          setSelectedResume(resumeRes.data.resumes[0].id);
        }
        
        if (appRes.data.applications) {
          const applied = appRes.data.applications.some(app => app.position_id === parseInt(id));
          setHasApplied(applied);
        }
      } catch (e) {
        setError(e instanceof Error ? e.message : '加载岗位详情失败');
      } finally {
        setLoading(false);
      }
    };
    
    fetchData();
  }, [id, isLoggedIn]);

  const handleApply = async () => {
    if (!selectedResume || !id) return;
    setApplying(true);
    try {
      const res = await apply(parseInt(id), selectedResume);
      if (res.data.application_id > 0) {
        setShowConfirm(true);
      }
    } catch (e) {
      setApplyError(e instanceof Error ? e.message : '投递失败');
    } finally {
      setApplying(false);
    }
  };

  const handleConfirmClose = () => {
    setShowConfirm(false);
    navigate('/applications');
  };

  if (loading) {
    return <Loading text="加载岗位详情..." />;
  }

  if (error || !position) {
    return <EmptyState title="岗位不存在" description={error} />;
  }

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <button onClick={() => navigate('/')} className="p-2 hover:bg-gray-100 rounded-lg">
          <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7"/>
          </svg>
        </button>
        <h1 className="text-2xl font-bold text-gray-900">{position.title || '岗位名称未填写'}</h1>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        <div className="flex flex-wrap items-center gap-4 mb-4">
          {position.salary_min > 0 || position.salary_max > 0 ? (
            <span className="text-xl font-bold text-primary-600">
              {formatSalary(position.salary_min, position.salary_max)}
            </span>
          ) : (
            <span className="text-lg text-gray-500">薪资面议</span>
          )}

          {position.location && (
            <span className="flex items-center gap-2 text-gray-600">
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
              </svg>
              {position.location}
            </span>
          )}

          <span className={`px-3 py-1 rounded-full text-sm ${
            position.status === 'published' ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-500'
          }`}>
            {position.status === 'published' ? '发布中' : '已下架'}
          </span>
        </div>

        <p className="text-sm text-gray-500 mb-6">发布于 {formatDateTime(position.created_at)}</p>

        <div className="mb-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-2">岗位职责</h3>
          <p className="text-gray-600 leading-relaxed">
            {position.description || '暂无描述'}
          </p>
        </div>

        <div>
          <h3 className="text-lg font-semibold text-gray-900 mb-2">任职要求</h3>
          <p className="text-gray-600 leading-relaxed">
            {position.requirements || '暂无要求'}
          </p>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">投递申请</h3>

        {!isLoggedIn ? (
          <div className="text-center py-8">
            <p className="text-gray-500 mb-4">请先登录后投递岗位</p>
            <Link to="/login" className="text-primary-600 hover:underline">
              立即登录
            </Link>
          </div>
        ) : !isProfileComplete ? (
          <div className="text-center py-8">
            <p className="text-gray-500 mb-4">请先完善个人资料再投递</p>
            <Link to="/profile" className="text-primary-600 hover:underline">
              完善资料
            </Link>
          </div>
        ) : resumes.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-500 mb-4">请先上传简历再投递</p>
            <Link to="/resume" className="text-primary-600 hover:underline">
              上传简历
            </Link>
          </div>
        ) : hasApplied ? (
          <div className="text-center py-8">
            <svg className="w-16 h-16 text-green-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
            <p className="text-gray-600 font-medium">已投递</p>
            <p className="text-sm text-gray-500 mt-2">您已投递此岗位，等待HR审核</p>
          </div>
        ) : (
          <>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">选择简历</label>
              <select
                value={selectedResume || ''}
                onChange={(e) => setSelectedResume(parseInt(e.target.value))}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
              >
                {resumes.map((resume) => (
                  <option key={resume.id} value={resume.id}>
                    {resume.file_name}
                  </option>
                ))}
              </select>
            </div>

            {applyError && <p className="text-red-500 text-sm mb-4">{applyError}</p>}

            <button
              onClick={() => { setSelectedResume(resumes[0]?.id || null); handleApply(); }}
              disabled={!selectedResume || applying}
              className="w-full py-3 bg-primary-500 text-white rounded-lg hover:bg-primary-600 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {applying ? '投递中...' : '立即投递'}
            </button>
          </>
        )}
      </div>

      <ConfirmDialog
        isOpen={showConfirm}
        title="投递成功"
        message="您的投递已成功提交，等待HR审核"
        onConfirm={handleConfirmClose}
        onCancel={handleConfirmClose}
        confirmText="查看投递记录"
      />
    </div>
  );
}