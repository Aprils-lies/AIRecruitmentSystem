import { useLocation, useNavigate } from 'react-router-dom';
import { removeToken } from '../utils/token';

const navItems = [
  { path: '/positions', label: '岗位管理', icon: 'briefcase' },
  { path: '/ai-chat', label: 'AI 智能助手', icon: 'robot' },
];

export function Sidebar() {
  const location = useLocation();
  const navigate = useNavigate();

  const handleLogout = () => {
    removeToken();
    navigate('/login');
  };

  return (
    <aside className="w-64 bg-gray-800 text-white min-h-screen">
      <div className="p-6 border-b border-gray-700">
        <h1 className="text-xl font-bold">智能招聘系统</h1>
        <p className="text-gray-400 text-sm mt-1">HR 管理端</p>
      </div>
      <nav className="p-4">
        <ul className="space-y-1">
          {navItems.map((item) => {
            const isActive = location.pathname === item.path;
            return (
              <li key={item.path}>
                <button
                  onClick={() => navigate(item.path)}
                  className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
                    isActive ? 'bg-primary-600 text-white' : 'text-gray-300 hover:bg-gray-700'
                  }`}
                >
                  <Icon name={item.icon} />
                  <span>{item.label}</span>
                </button>
              </li>
            );
          })}
        </ul>
      </nav>
      <div className="absolute bottom-0 left-0 w-64 p-4 border-t border-gray-700">
        <button
          onClick={handleLogout}
          className="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-300 hover:bg-gray-700 transition-colors"
        >
          <Icon name="logout" />
          <span>退出登录</span>
        </button>
      </div>
    </aside>
  );
}

interface IconProps {
  name: string;
}

function Icon({ name }: IconProps) {
  const icons: Record<string, string> = {
    briefcase: 'M20 7h-9m0 0V5a2 2 0 012-2h4a2 2 0 012 2v2m-9 0V13a2 2 0 002 2h4a2 2 0 002-2V7m-9 0h9',
    robot: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10',
    logout: 'M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1',
  };

  return (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={icons[name]} />
    </svg>
  );
}
