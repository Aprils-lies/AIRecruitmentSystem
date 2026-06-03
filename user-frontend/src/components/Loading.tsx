interface LoadingProps {
  text?: string;
}

export default function Loading({ text = '加载中...' }: LoadingProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12">
      <div className="relative">
        <div className="w-12 h-12 border-4 border-primary-200 rounded-full animate-spin border-t-primary-500"></div>
      </div>
      <p className="mt-4 text-gray-500">{text}</p>
    </div>
  );
}