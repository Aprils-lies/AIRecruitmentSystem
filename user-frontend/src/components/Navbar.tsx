import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

interface NavbarProps {
  isLoggedIn: boolean;
  onLogout: () => void;
}

export default function Navbar({ isLoggedIn, onLogout }: NavbarProps) {
  const [showDropdown, setShowDropdown] = useState(false);
  const navigate = useNavigate();

  const handleLogout = () => {
    onLogout();
    setShowDropdown(false);
    navigate('/');
  };

  return (
    <nav className="bg-white shadow-sm sticky top-0 z-40">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center gap-8">
            <Link to="/" className="text-xl font-bold text-primary-600">
              智能招聘
            </Link>
            <div className="hidden md:flex items-center gap-6">
              <Link to="/" className="text-gray-600 hover:text-primary-600">
                岗位列表
              </Link>
              {isLoggedIn && (
                <>
                  <Link to="/profile" className="text-gray-600 hover:text-primary-600">
                    个人资料
                  </Link>
                  <Link to="/resume" className="text-gray-600 hover:text-primary-600">
                    简历管理
                  </Link>
                  <Link to="/applications" className="text-gray-600 hover:text-primary-600">
                    我的投递
                  </Link>
                </>
              )}
            </div>
          </div>

          <div className="flex items-center gap-4">
            {isLoggedIn ? (
              <div className="relative">
                <button
                  onClick={() => setShowDropdown(!showDropdown)}
                  className="flex items-center gap-2 px-3 py-2 text-gray-600 hover:text-primary-600"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  <span className="hidden sm:inline">我的</span>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
                
                {showDropdown && (
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-100 py-2">
                    <Link to="/profile" className="block px-4 py-2 text-gray-700 hover:bg-gray-50">
                      个人资料
                    </Link>
                    <Link to="/resume" className="block px-4 py-2 text-gray-700 hover:bg-gray-50">
                      简历管理
                    </Link>
                    <Link to="/applications" className="block px-4 py-2 text-gray-700 hover:bg-gray-50">
                      我的投递
                    </Link>
                    <hr className="my-2 border-gray-100" />
                    <button
                      onClick={handleLogout}
                      className="w-full text-left px-4 py-2 text-red-600 hover:bg-gray-50"
                    >
                      退出登录
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <Link to="/login" className="px-4 py-2 text-gray-600 hover:text-primary-600">
                  登录
                </Link>
                <Link
                  to="/register"
                  className="px-4 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600"
                >
                  注册
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}