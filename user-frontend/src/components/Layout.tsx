import Navbar from './Navbar';
import ToastComponent from './Toast';
import type { Toast } from '../hooks/useToast';

interface LayoutProps {
  isLoggedIn: boolean;
  onLogout: () => void;
  toasts: Toast[];
  children: React.ReactNode;
}

export default function Layout({ isLoggedIn, onLogout, toasts, children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar isLoggedIn={isLoggedIn} onLogout={onLogout} />
      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>
      <ToastComponent toasts={toasts} />
    </div>
  );
}