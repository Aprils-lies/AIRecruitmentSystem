export function Loading() {
  return (
    <div className="flex items-center justify-center py-12">
      <div className="flex flex-col items-center">
        <div className="w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full animate-spin" />
        <span className="mt-2 text-gray-500 text-sm">加载中...</span>
      </div>
    </div>
  );
}

export function LoadingButton() {
  return (
    <span className="flex items-center gap-2">
      <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
      处理中...
    </span>
  );
}
