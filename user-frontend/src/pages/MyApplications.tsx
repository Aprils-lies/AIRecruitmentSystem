import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { listMyApplications, withdrawApplication } from '../api/application';
import { normalizeApplication } from '../utils/nullsafe';
import { formatDateTime, formatApplicationStatus, getStatusColor } from '../utils/format';
import Loading from '../components/Loading';
import EmptyState from '../components/EmptyState';
import Pagination from '../components/Pagination';
import ConfirmDialog from '../components/ConfirmDialog';
import { useToast } from '../hooks/useToast';

export default function MyApplications() {
  const [applications, setApplications] = useState<ReturnType<typeof normalizeApplication>[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showConfirm, setShowConfirm] = useState(false);
  const [applicationToWithdraw, setApplicationToWithdraw] = useState<number | null>(null);
  const [withdrawing, setWithdrawing] = useState(false);
  const navigate = useNavigate();
  const { success, error: showError } = useToast();

  const pageSize = 10;

  const fetchApplications = useCallback(async (currentPage: number) => {
    setLoading(true);
    setError(null);
    try {
      const res = await listMyApplications({
        page: currentPage,
        page_size: pageSize,
      });
      const data = res.data;
      setApplications(data.applications?.map((a: unknown) => normalizeApplication(a)) || []);
      setTotal(data.total || 0);
      setPage(currentPage);
    } catch (e) {
      setError(e instanceof Error ? e.message : '获取投递记录失败');
      setApplications([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchApplications(1);
  }, [fetchApplications]);

  const handleWithdrawClick = (applicationId: number) => {
    setApplicationToWithdraw(applicationId);
    setShowConfirm(true);
  };

  const handleConfirmWithdraw = async () => {
    if (!applicationToWithdraw) return;
    setWithdrawing(true);
    try {
      await withdrawApplication(applicationToWithdraw);
      success('撤回成功');
      setShowConfirm(false);
      fetchApplications(page);
    } catch (e) {
      showError(e instanceof Error ? e.message : '撤回失败');
    } finally {
      setWithdrawing(false);
    }
  };

  const handleCancelWithdraw = () => {
    setShowConfirm(false);
    setApplicationToWithdraw(null);
  };

  const canWithdraw = (status: string) => status === 'pending';

  if (loading) {
    return <Loading text="加载投递记录..." />;
  }

  if (error) {
    return <EmptyState title="加载失败" description={error} />;
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">我的投递</h1>

      {applications.length === 0 ? (
        <EmptyState title="暂无投递记录" description="快去浏览岗位并投递吧" icon="briefcase" />
      ) : (
        <>
          <div className="bg-white rounded-lg shadow-md overflow-hidden">
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">岗位名称</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">投递状态</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">投递时间</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {applications.map((app) => (
                  <tr key={app.application_id} className="hover:bg-gray-50">
                    <td className="px-6 py-4">
                      <div className="text-sm font-medium text-gray-900">{app.position_title}</div>
                    </td>
                    <td className="px-6 py-4">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(app.status)}`}>
                        {formatApplicationStatus(app.status)}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-500">
                      {formatDateTime(app.applied_at)}
                    </td>
                    <td className="px-6 py-4 space-x-2">
                      <button
                        onClick={() => navigate(`/positions/${app.position_id}`)}
                        className="text-primary-600 hover:text-primary-900 text-sm font-medium"
                      >
                        查看岗位
                      </button>
                      {canWithdraw(app.status) && (
                        <button
                          onClick={() => handleWithdrawClick(app.application_id)}
                          disabled={withdrawing}
                          className="text-red-600 hover:text-red-900 text-sm font-medium disabled:opacity-50"
                        >
                          撤回
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <Pagination
            currentPage={page}
            total={total}
            pageSize={pageSize}
            onPageChange={fetchApplications}
          />
        </>
      )}

      <ConfirmDialog
        isOpen={showConfirm}
        title="确认撤回"
        message="确定要撤回这次投递吗？撤回后该记录将不再显示。"
        onConfirm={handleConfirmWithdraw}
        onCancel={handleCancelWithdraw}
        confirmText="确认撤回"
        cancelText="取消"
      />
    </div>
  );
}