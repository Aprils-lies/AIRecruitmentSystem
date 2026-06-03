import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { register } from '../api/auth';
import { LoadingButton } from '../components/Loading';

interface Toast {
  id: number;
  message: string;
  type: 'success' | 'error' | 'warning';
}

export default function Register() {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [toasts, setToasts] = useState<Toast[]>([]);

  const showToast = (message: string, type: 'success' | 'error' | 'warning') => {
    const id = Date.now();
    setToasts(prev => [...prev, { id, message, type }]);
    setTimeout(() => {
      setToasts(prev => prev.filter(t => t.id !== id));
    }, 3000);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim()) {
      showToast('请输入用户名', 'warning');
      return;
    }
    if (!password) {
      showToast('请输入密码', 'warning');
      return;
    }
    if (password !== confirmPassword) {
      showToast('两次输入的密码不一致', 'warning');
      return;
    }

    setLoading(true);
    try {
      await register({ username, password, role: 'hr' });
      showToast('注册成功，请登录', 'success');
      setTimeout(() => navigate('/login'), 1000);
    } catch (error) {
      showToast(error instanceof Error ? error.message : '注册失败', 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white rounded-card shadow-modal p-8 w-full max-w-md mx-4">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900">智能招聘系统</h1>
          <p className="text-gray-500 mt-2">HR 注册</p>
        </div>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">用户名</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="请输入用户名"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">密码</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="请输入密码"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">确认密码</label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-btn focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="请再次输入密码"
            />
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-primary-500 text-white rounded-btn hover:bg-primary-600 transition-colors disabled:opacity-50"
          >
            {loading ? <LoadingButton /> : '注册'}
          </button>
        </form>
        
        <div className="mt-6 text-center">
          <p className="text-gray-500 text-sm">
            已有账号？ <Link to="/login" className="text-primary-500 hover:text-primary-600">立即登录</Link>
          </p>
        </div>
      </div>
      
      {toasts.map(toast => (
        <div
          key={toast.id}
          className={`fixed top-4 right-4 px-4 py-3 rounded-lg text-white shadow-lg ${
            toast.type === 'success' ? 'bg-success' : toast.type === 'error' ? 'bg-danger' : 'bg-warning'
          }`}
        >
          {toast.message}
        </div>
      ))}
    </div>
  );
}
