import { useState, useEffect } from 'react';
import type { PositionInfo, CreatePositionRequest } from '../types';
import { listMyPositions, createPosition, updatePosition, offlinePosition, onlinePosition } from '../api/position';
import { listCandidates, getCandidateDetail, updateApplicationStatus } from '../api/application';
import { formatDateTime, formatSalary } from '../utils/format';
import { Loading } from '../components/Loading';
import { EmptyState } from '../components/EmptyState';
import { Pagination } from '../components/Pagination';
import { PositionForm } from '../components/PositionForm';
import { CandidateDetail } from '../components/CandidateDetail';
import { ConfirmDialog } from '../components/ConfirmDialog';
import { ToastContainer } from '../components/Toast';

interface Toast {
  id: number;
  message: string;
  type: 'success' | 'error' | 'warning';
}

interface CandidateBrief {
  user_id: number;
  application_id: number;
  username: string;
  real_name: string;
  phone: string;
  education: string;
  school: string;
  skills: string;
  resume_id: number;
  resume_name: string;
  applied_at: string;
  status: string;
}

interface CandidateDetailType {
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

export default function Positions() {
  const [positions, setPositions] = useState<PositionInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [showForm, setShowForm] = useState(false);
  const [editingPosition, setEditingPosition] = useState<PositionInfo | null>(null);
  const [showConfirm, setShowConfirm] = useState(false);
  const [offlineTargetId, setOfflineTargetId] = useState<number>(0);
  const [showOnlineConfirm, setShowOnlineConfirm] = useState(false);
  const [onlineTargetId, setOnlineTargetId] = useState<number>(0);
  const [toasts, setToasts] = useState<Toast[]>([]);
  
  const [showCandidateList, setShowCandidateList] = useState(false);
  const [currentPositionId, setCurrentPositionId] = useState(0);
  const [candidateDetail, setCandidateDetail] = useState<CandidateDetailType | null>(null);

  const showToast = (message: string, type: 'success' | 'error' | 'warning') => {
    const id = Date.now();
    setToasts(prev => [...prev, { id, message, type }]);
  };

  const removeToast = (id: number) => {
    setToasts(prev => prev.filter(t => t.id !== id));
  };

  const loadPositions = async () => {
    setLoading(true);
    try {
      const resp = await listMyPositions({ page, page_size: pageSize });
      setPositions(resp.data.positions);
      setTotal(resp.data.total);
    } catch (error) {
      showToast(error instanceof Error ? error.message : '加载失败', 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadPositions();
  }, [page, pageSize]);

  const handlePageSizeChange = (newPageSize: number) => {
    setPageSize(newPageSize);
    setPage(1);
  };

  const handleCreate = async (data: CreatePositionRequest) => {
    try {
      await createPosition(data);
      showToast('创建成功', 'success');
      setShowForm(false);
      loadPositions();
    } catch (error) {
      showToast(error instanceof Error ? error.message : '创建失败', 'error');
    }
  };

  const handleEdit = async (data: CreatePositionRequest) => {
    if (!editingPosition) return;
    try {
      await updatePosition(editingPosition.id, data);
      showToast('更新成功', 'success');
      setShowForm(false);
      setEditingPosition(null);
      loadPositions();
    } catch (error) {
      showToast(error instanceof Error ? error.message : '更新失败', 'error');
    }
  };

  const handleOffline = async () => {
    try {
      await offlinePosition(offlineTargetId);
      showToast('已下架', 'success');
      setShowConfirm(false);
      loadPositions();
    } catch (error) {
      showToast(error instanceof Error ? error.message : '操作失败', 'error');
    }
  };

  const handleOnline = async () => {
    try {
      await onlinePosition(onlineTargetId);
      showToast('已上架', 'success');
      setShowOnlineConfirm(false);
      loadPositions();
    } catch (error) {
      showToast(error instanceof Error ? error.message : '操作失败', 'error');
    }
  };

  const handleViewCandidates = async (positionId: number) => {
    setCurrentPositionId(positionId);
    setShowCandidateList(true);
  };

  return (
    <div className="max-w-6xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold text-gray-900">岗位管理</h2>
        <button
          onClick={() => { setEditingPosition(null); setShowForm(true); }}
          className="px-4 py-2 bg-primary-500 text-white rounded-btn hover:bg-primary-600 transition-colors"
        >
          + 新建岗位
        </button>
      </div>

      {loading ? (
        <Loading />
      ) : !positions || positions.length === 0 ? (
        <EmptyState title="暂无岗位" description="点击上方按钮创建新岗位" />
      ) : (
        <div className="bg-white rounded-card shadow-card">
          <table className="w-full">
            <thead>
              <tr className="bg-gray-50">
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">岗位名称</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">薪资</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">地点</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">状态</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">创建时间</th>
                <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">操作</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {(positions || []).map((position) => (
                <tr key={position.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">
                    <div>
                      <p className="font-medium text-gray-900">{position.title}</p>
                      <p className="text-sm text-gray-500 line-clamp-2 max-w-xs">{position.description}</p>
                    </div>
                  </td>
                  <td className="px-6 py-4 text-gray-700">{formatSalary(position.salary_min, position.salary_max)}</td>
                  <td className="px-6 py-4 text-gray-700">{position.location || '-'}</td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 text-xs font-medium rounded-full ${
                      position.status === 'published' 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-gray-100 text-gray-600'
                    }`}>
                      {position.status === 'published' ? '发布中' : '已下架'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-gray-500 text-sm">{formatDateTime(position.created_at)}</td>
                  <td className="px-6 py-4">
                    <div className="flex items-center justify-center gap-2">
                      <button
                        onClick={() => handleViewCandidates(position.id)}
                        className="px-3 py-1 text-sm text-blue-600 hover:bg-blue-50 rounded-btn transition-colors"
                      >
                        查看候选人
                      </button>
                      <button
                        onClick={() => { setEditingPosition(position); setShowForm(true); }}
                        className="px-3 py-1 text-sm text-gray-600 hover:bg-gray-100 rounded-btn transition-colors"
                      >
                        编辑
                      </button>
                      {position.status === 'published' ? (
                        <button
                          onClick={() => { setOfflineTargetId(position.id); setShowConfirm(true); }}
                          className="px-3 py-1 text-sm text-red-600 hover:bg-red-50 rounded-btn transition-colors"
                        >
                          下架
                        </button>
                      ) : (
                        <button
                          onClick={() => { setOnlineTargetId(position.id); setShowOnlineConfirm(true); }}
                          className="px-3 py-1 text-sm text-green-600 hover:bg-green-50 rounded-btn transition-colors"
                        >
                          上架
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <Pagination 
            total={total} 
            page={page} 
            pageSize={pageSize} 
            onPageChange={setPage}
            onPageSizeChange={handlePageSizeChange}
          />
        </div>
      )}

      {showForm && (
        <PositionForm
          title={editingPosition ? '编辑岗位' : '新建岗位'}
          position={editingPosition}
          onSubmit={editingPosition ? handleEdit : handleCreate}
          onClose={() => { setShowForm(false); setEditingPosition(null); }}
        />
      )}

      {showConfirm && (
        <ConfirmDialog
          title="确认下架"
          message="确认要下架此岗位吗？下架后候选人将无法投递。"
          onConfirm={handleOffline}
          onCancel={() => setShowConfirm(false)}
        />
      )}

      {showOnlineConfirm && (
        <ConfirmDialog
          title="确认上架"
          message="确认要上架此岗位吗？上架后候选人可以投递。"
          onConfirm={handleOnline}
          onCancel={() => setShowOnlineConfirm(false)}
        />
      )}

      {showCandidateList && (
        <CandidateList
          positionId={currentPositionId}
          onClose={() => setShowCandidateList(false)}
          onShowDetail={setCandidateDetail}
          onToast={showToast}
        />
      )}

      {candidateDetail && (
        <CandidateDetail
          candidate={candidateDetail}
          onClose={() => setCandidateDetail(null)}
          onToast={showToast}
        />
      )}

      <ToastContainer toasts={toasts} onRemove={removeToast} />
    </div>
  );
}

function CandidateList({ positionId, onClose, onShowDetail, onToast }: {
  positionId: number;
  onClose: () => void;
  onShowDetail: (detail: CandidateDetailType) => void;
  onToast: (message: string, type: 'success' | 'error' | 'warning') => void;
}) {
  const [candidates, setCandidates] = useState<CandidateBrief[]>([]);
  const [loading, setLoading] = useState(true);
  const [updatingId, setUpdatingId] = useState<number | null>(null);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);

  const statusOptions = [
    { value: 'pending', label: '待审核' },
    { value: 'reviewed', label: '已查看' },
    { value: 'rejected', label: '已拒绝' },
    { value: 'accepted', label: '已通过' },
  ];

  useEffect(() => {
    const loadCandidates = async () => {
      setLoading(true);
      try {
        const resp = await listCandidates({ positionId, page, page_size: pageSize });
        setCandidates(resp.data.candidates);
        setTotal(resp.data.total);
      } catch (error) {
        onToast(error instanceof Error ? error.message : '加载失败', 'error');
      } finally {
        setLoading(false);
      }
    };
    loadCandidates();
  }, [positionId, page, pageSize]);

  const handlePageSizeChange = (newPageSize: number) => {
    setPageSize(newPageSize);
    setPage(1); // 重置到第一页
  };

  const handleViewDetail = async (candidate: CandidateBrief) => {
    try {
      const resp = await getCandidateDetail({ candidateId: candidate.user_id, applicationId: candidate.application_id });
      onShowDetail(resp.data);
    } catch (error) {
      onToast(error instanceof Error ? error.message : '加载详情失败', 'error');
    }
  };

  const handleStatusChange = async (candidate: CandidateBrief, newStatus: string) => {
    if (candidate.status === newStatus) return;
    
    setUpdatingId(candidate.application_id);
    try {
      await updateApplicationStatus(candidate.application_id, newStatus);
      onToast('状态更新成功', 'success');
      // 刷新列表
      const resp = await listCandidates({ positionId, page, page_size: 10 });
      setCandidates(resp.data.candidates || []);
    } catch (error) {
      onToast(error instanceof Error ? error.message : '状态更新失败', 'error');
    } finally {
      setUpdatingId(null);
    }
  };

  const getStatusColor = (status: string) => {
    const map: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800',
      reviewed: 'bg-blue-100 text-blue-800',
      rejected: 'bg-red-100 text-red-800',
      accepted: 'bg-green-100 text-green-800',
    };
    return map[status] || 'bg-gray-100 text-gray-800';
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white rounded-card shadow-modal p-6 w-full max-w-4xl mx-4 max-h-[80vh] flex flex-col">
        <div className="flex items-center justify-between mb-6">
          <h3 className="text-lg font-semibold text-gray-900">候选人列表</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div className="flex-1 overflow-y-auto">
          {loading ? (
            <Loading />
          ) : !candidates || candidates.length === 0 ? (
            <EmptyState title="暂无候选人" description="该岗位还没有候选人投递" />
          ) : (
            <table className="w-full">
              <thead>
                <tr className="bg-gray-50">
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">姓名</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">学历</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">院校</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">技能</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">状态</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500">投递时间</th>
                  <th className="px-4 py-2 text-center text-xs font-medium text-gray-500">操作</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {(candidates || []).map((candidate) => (
                  <tr key={candidate.user_id}>
                    <td className="px-4 py-3">
                      <p className="font-medium text-gray-900">{candidate.real_name}</p>
                      <p className="text-xs text-gray-500">{candidate.phone}</p>
                    </td>
                    <td className="px-4 py-3 text-gray-700">{candidate.education}</td>
                    <td className="px-4 py-3 text-gray-700">{candidate.school}</td>
                    <td className="px-4 py-3">
                      <div className="flex flex-wrap gap-1">
                        {((candidate.skills || '').split(',')).slice(0, 3).map((skill, i) => (
                          <span key={i} className="px-2 py-0.5 bg-gray-100 text-gray-600 rounded text-xs">
                            {skill.trim()}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-2">
                        {updatingId === candidate.application_id ? (
                          <span className="text-xs text-gray-500">更新中...</span>
                        ) : (
                          <select
                            value={candidate.status}
                            onChange={(e) => handleStatusChange(candidate, e.target.value)}
                            className={`px-2 py-1 text-xs font-medium rounded-full border-0 appearance-none cursor-pointer ${getStatusColor(candidate.status)}`}
                          >
                            {statusOptions.map((option) => (
                              <option key={option.value} value={option.value}>
                                {option.label}
                              </option>
                            ))}
                          </select>
                        )}
                      </div>
                    </td>
                    <td className="px-4 py-3 text-gray-500 text-sm">{formatDateTime(candidate.applied_at)}</td>
                    <td className="px-4 py-3">
                      <button
                        onClick={() => handleViewDetail(candidate)}
                        className="px-3 py-1 text-sm text-blue-600 hover:bg-blue-50 rounded-btn transition-colors"
                      >
                        查看详情
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        
        {!loading && (candidates || []).length > 0 && (
          <Pagination 
            total={total} 
            page={page} 
            pageSize={pageSize} 
            onPageChange={setPage}
            onPageSizeChange={handlePageSizeChange}
          />
        )}
      </div>
    </div>
  );
}
