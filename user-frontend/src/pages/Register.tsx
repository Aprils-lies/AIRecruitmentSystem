import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import type { RegisterRequest } from '../types';

interface RegisterProps {
  onRegister: (data: RegisterRequest) => Promise<boolean>;
  loading: boolean;
  error: string | null;
}

export default function Register({ onRegister, loading, error }: RegisterProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username.trim() || !password || !confirmPassword) {
      return;
    }
    if (password !== confirmPassword) {
      return;
    }
    const success = await onRegister({ username: username.trim(), password, role: 'candidate' });
    if (success) {
      navigate('/login');
    }
  };

  const passwordsMatch = password === confirmPassword;

  return (
    <div className="max-w-md mx-auto">
      <div className="text-center mb-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">注册账号</h1>
        <p className="text-gray-500">创建您的候选人账号</p>
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">用户名</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="请输入用户名"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">密码</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="请输入密码"
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
          </div>

          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">确认密码</label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="请再次输入密码"
              className={`w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 ${
                confirmPassword && !passwordsMatch ? 'border-red-500' : 'border-gray-300'
              }`}
            />
            {confirmPassword && !passwordsMatch && (
              <p className="text-red-500 text-sm mt-1">两次输入的密码不一致</p>
            )}
          </div>

          {error && (
            <p className="text-red-500 text-sm mb-4">{error}</p>
          )}

          <button
            type="submit"
            disabled={loading || !username.trim() || !password || !confirmPassword || !passwordsMatch}
            className="w-full py-3 bg-primary-500 text-white rounded-lg hover:bg-primary-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? '注册中...' : '注册'}
          </button>
        </form>

        <p className="mt-4 text-center text-gray-500">
          已有账号？{' '}
          <Link to="/login" className="text-primary-600 hover:underline">
            立即登录
          </Link>
        </p>
      </div>
    </div>
  );
}