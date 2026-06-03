import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { login } from '../api/auth';
import { setToken } from '../utils/token';
import { LoadingButton } from '../components/Loading';

interface Toast {
  id: number;
  message: string;
  type: 'success' | 'error' | 'warning';
}

export default function Login() {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
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
    if (!username.trim() || !password.trim()) {
      showToast('请填写用户名和密码', 'warning');
      return;
    }

    setLoading(true);
    try {
      const resp = await login({ username, password });
      const { token, role } = resp.data;
      
      if (role !== 'hr') {
        showToast('您不是HR用户，无法登录HR管理端', 'error');
        return;
      }
      
      setToken(token);
      showToast('登录成功', 'success');
      setTimeout(() => navigate('/positions'), 1000);
    } catch (error) {
      showToast(error instanceof Error ? error.message : '登录失败', 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white rounded-card shadow-modal p-8 w-full max-w-md mx-4">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900">智能招聘系统</h1>
          <p className="text-gray-500 mt-2">HR 管理端</p>
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
          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-primary-500 text-white rounded-btn hover:bg-primary-600 transition-colors disabled:opacity-50"
          >
            {loading ? <LoadingButton /> : '登录'}
          </button>
        </form>
        
        <div className="mt-6 text-center">
          <p className="text-gray-500 text-sm">
            还没有账号？ <Link to="/register" className="text-primary-500 hover:text-primary-600">立即注册</Link>
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
