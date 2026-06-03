import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { listPositions } from '../api/position';
import { normalizePosition } from '../utils/nullsafe';
import PositionCard from '../components/PositionCard';
import Pagination from '../components/Pagination';
import Loading from '../components/Loading';
import EmptyState from '../components/EmptyState';

export default function PositionList() {
  const [positions, setPositions] = useState<ReturnType<typeof normalizePosition>[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(9);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [keyword, setKeyword] = useState('');
  const [location, setLocation] = useState('');
  const navigate = useNavigate();

  const fetchPositions = useCallback(async (currentPage: number) => {
    setLoading(true);
    setError(null);
    try {
      const res = await listPositions({
        page: currentPage,
        page_size: pageSize,
        keyword: keyword || undefined,
        location: location || undefined,
      });
      const data = res.data;
      setPositions(data.positions?.map((p: unknown) => normalizePosition(p)) || []);
      setTotal(data.total || 0);
      setPage(currentPage);
    } catch (e) {
      setError(e instanceof Error ? e.message : '获取岗位列表失败');
      setPositions([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [keyword, location, pageSize]);

  const handlePageSizeChange = useCallback((newPageSize: number) => {
    setPageSize(newPageSize);
    setPage(1);
  }, []);

  useEffect(() => {
    fetchPositions(1);
  }, [fetchPositions]);

  useEffect(() => {
    if (pageSize !== 9 && page !== 1) {
      fetchPositions(1);
    }
  }, [pageSize]);

  const handleSearch = () => {
    fetchPositions(1);
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">招聘岗位</h1>
        <p className="text-gray-500">浏览最新的招聘机会，找到适合您的职位</p>
      </div>

      <div className="flex flex-col md:flex-row gap-4 mb-6">
        <div className="flex-1 relative">
          <input
            type="text"
            placeholder="搜索岗位关键词..."
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onKeyPress={handleKeyPress}
            className="w-full px-4 py-2 pl-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
          <svg className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        
        <div className="flex gap-3">
          <input
            type="text"
            placeholder="工作地点..."
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
          <button
            onClick={handleSearch}
            className="px-6 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600"
          >
            搜索
          </button>
        </div>
      </div>

      {loading ? (
        <Loading text="加载岗位列表..." />
      ) : error ? (
        <EmptyState title="加载失败" description={error} />
      ) : positions.length === 0 ? (
        <EmptyState title="暂无岗位" description="目前没有符合条件的招聘岗位" icon="briefcase" />
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {positions.map((position) => (
              <PositionCard
                key={position.id}
                position={position}
                onClick={() => navigate(`/positions/${position.id}`)}
              />
            ))}
          </div>
          
          <Pagination
            currentPage={page}
            total={total}
            pageSize={pageSize}
            onPageChange={fetchPositions}
            onPageSizeChange={handlePageSizeChange}
          />
        </>
      )}
    </div>
  );
}