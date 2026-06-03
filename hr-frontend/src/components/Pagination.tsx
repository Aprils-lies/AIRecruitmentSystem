interface PaginationProps {
  total: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  onPageSizeChange?: (pageSize: number) => void;
}

export function Pagination({ total, page, pageSize, onPageChange, onPageSizeChange }: PaginationProps) {
  const totalPages = Math.ceil(total / pageSize);
  
  const pageSizeOptions = [
    { value: 5, label: '5条/页' },
    { value: 10, label: '10条/页' },
    { value: 15, label: '15条/页' },
  ];

  const getPageNumbers = () => {
    const pages: number[] = [];
    if (totalPages <= 5) {
      for (let i = 1; i <= totalPages; i++) pages.push(i);
    } else {
      if (page <= 3) {
        pages.push(1, 2, 3, 4, totalPages);
      } else if (page >= totalPages - 2) {
        pages.push(1, totalPages - 3, totalPages - 2, totalPages - 1, totalPages);
      } else {
        pages.push(1, page - 1, page, page + 1, totalPages);
      }
    }
    return pages;
  };

  const pageNumbers = getPageNumbers();

  return (
    <div className="flex items-center justify-between mt-4 px-4 py-3 bg-white border-t border-gray-200">
      {onPageSizeChange && (
        <div className="flex items-center gap-2">
          <span className="text-sm text-gray-600">每页显示：</span>
          <select
            value={pageSize}
            onChange={(e) => onPageSizeChange(Number(e.target.value))}
            className="px-3 py-1 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            {pageSizeOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </div>
      )}

      <div className="flex items-center justify-center gap-1 flex-1">
        <button
          onClick={() => onPageChange(page - 1)}
          disabled={page === 1}
          className="px-3 py-1 rounded-btn text-sm text-gray-600 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          上一页
        </button>
        {pageNumbers.map((p, index) => {
          if (pageNumbers[index - 1] && p > pageNumbers[index - 1] + 1) {
            return <span key={`ellipsis-${p}`} className="px-2 text-gray-400">...</span>;
          }
          return (
            <button
              key={p}
              onClick={() => onPageChange(p)}
              className={`px-3 py-1 rounded-btn text-sm transition-colors ${
                page === p
                  ? 'bg-primary-500 text-white'
                  : 'text-gray-600 hover:bg-gray-100'
              }`}
            >
              {p}
            </button>
          );
        })}
        <button
          onClick={() => onPageChange(page + 1)}
          disabled={page === totalPages}
          className="px-3 py-1 rounded-btn text-sm text-gray-600 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          下一页
        </button>
      </div>

      <div className="text-sm text-gray-500 w-32 text-right">
        共 {total} 条
      </div>
    </div>
  );
}
