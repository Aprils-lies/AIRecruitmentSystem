import { RouteObject } from 'react-router-dom';
import { isLoggedIn } from '../utils/token';
import { Layout } from '../components/Layout';
import Login from '../pages/Login';
import Register from '../pages/Register';
import Positions from '../pages/Positions';
import AIChat from '../pages/AIChat';
import { Navigate } from 'react-router-dom';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  if (!isLoggedIn()) {
    return <Navigate to="/login" />;
  }
  return children;
}

function PublicRoute({ children }: { children: React.ReactNode }) {
  if (isLoggedIn()) {
    return <Navigate to="/positions" />;
  }
  return children;
}

export const routes: RouteObject[] = [
  {
    path: '/login',
    element: <PublicRoute><Login /></PublicRoute>,
  },
  {
    path: '/register',
    element: <PublicRoute><Register /></PublicRoute>,
  },
  {
    element: <ProtectedRoute><Layout /></ProtectedRoute>,
    children: [
      {
        path: '/positions',
        element: <Positions />,
      },
      {
        path: '/ai-chat',
        element: <AIChat />,
      },
    ],
  },
  {
    path: '/',
    element: <Navigate to="/positions" />,
  },
];
