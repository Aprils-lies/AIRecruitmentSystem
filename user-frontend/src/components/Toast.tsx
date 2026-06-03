interface ToastItem {
  id: string;
  message: string;
  type: 'success' | 'error' | 'warning';
}

interface ToastProps {
  toasts: ToastItem[];
}

export default function Toast({ toasts }: ToastProps) {
  const getStyles = (type: ToastItem['type']) => {
    switch (type) {
      case 'success':
        return 'bg-green-500 text-white';
      case 'error':
        return 'bg-red-500 text-white';
      case 'warning':
        return 'bg-yellow-500 text-white';
      default:
        return 'bg-gray-500 text-white';
    }
  };

  return (
    <div className="fixed top-4 right-4 z-50 flex flex-col gap-2">
      {toasts.map(toast => (
        <div
          key={toast.id}
          className={`px-4 py-3 rounded-lg shadow-lg ${getStyles(toast.type)} animate-fade-in`}
        >
          {toast.message}
        </div>
      ))}
    </div>
  );
}