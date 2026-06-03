import { Routes, Route, Navigate } from 'react-router-dom';
import Layout from '../components/Layout';
import PositionList from '../pages/PositionList';
import PositionDetail from '../pages/PositionDetail';
import Login from '../pages/Login';
import Register from '../pages/Register';
import Profile from '../pages/Profile';
import ResumeUpload from '../pages/ResumeUpload';
import MyApplications from '../pages/MyApplications';
import Loading from '../components/Loading';
import { useAuth } from '../hooks/useAuth';
import { useToast } from '../hooks/useToast';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isLoggedIn, loading } = useAuth();
  
  if (loading) {
    return <Loading />;
  }
  
  return isLoggedIn ? <>{children}</> : <Navigate to="/login" />;
}

function PublicRoute({ children }: { children: React.ReactNode }) {
  const { isLoggedIn, loading } = useAuth();
  
  if (loading) {
    return <Loading />;
  }
  
  return !isLoggedIn ? <>{children}</> : <Navigate to="/" />;
}

export function AppRoutes() {
  const { login, register, logout, isLoggedIn, loading, error, isProfileComplete, refreshProfile } = useAuth();
  const { toasts, success } = useToast();

  if (loading) {
    return <Loading />;
  }

  const handleProfileUpdateSuccess = async () => {
    success('个人资料更新成功');
    await refreshProfile();
  };

  const handleResumeUploadSuccess = () => {
    success('简历上传成功');
  };

  return (
    <Routes>
      <Route path="/" element={
        <Layout isLoggedIn={isLoggedIn} onLogout={logout} toasts={toasts}>
          <PositionList />
        </Layout>
      } />
      
      <Route path="/positions/:id" element={
        <Layout isLoggedIn={isLoggedIn} onLogout={logout} toasts={toasts}>
          <PositionDetail isLoggedIn={isLoggedIn} isProfileComplete={isProfileComplete} />
        </Layout>
      } />
      
      <Route path="/login" element={
        <PublicRoute>
          <Login onLogin={login} loading={loading} error={error} />
        </PublicRoute>
      } />
      
      <Route path="/register" element={
        <PublicRoute>
          <Register onRegister={register} loading={loading} error={error} />
        </PublicRoute>
      } />
      
      <Route path="/profile" element={
        <ProtectedRoute>
          <Layout isLoggedIn={isLoggedIn} onLogout={logout} toasts={toasts}>
            <Profile onUpdateSuccess={handleProfileUpdateSuccess} />
          </Layout>
        </ProtectedRoute>
      } />
      
      <Route path="/resume" element={
        <ProtectedRoute>
          <Layout isLoggedIn={isLoggedIn} onLogout={logout} toasts={toasts}>
            <ResumeUpload onSuccess={handleResumeUploadSuccess} />
          </Layout>
        </ProtectedRoute>
      } />
      
      <Route path="/applications" element={
        <ProtectedRoute>
          <Layout isLoggedIn={isLoggedIn} onLogout={logout} toasts={toasts}>
            <MyApplications />
          </Layout>
        </ProtectedRoute>
      } />
    </Routes>
  );
}