import { useState, useEffect, useCallback } from 'react';
import { login as apiLogin, register as apiRegister } from '../api/auth';
import { getProfile } from '../api/user';
import { setToken, getToken, removeToken } from '../utils/token';
import type { LoginRequest, RegisterRequest, LoginData, UserProfile } from '../types';

interface AuthState {
  user: LoginData | null;
  profile: UserProfile | null;
  loading: boolean;
  error: string | null;
}

export function useAuth() {
  const [state, setState] = useState<AuthState>({
    user: null,
    profile: null,
    loading: true,
    error: null,
  });

  const fetchProfile = useCallback(async () => {
    if (!getToken()) {
      setState(prev => ({ ...prev, loading: false }));
      return;
    }
    try {
      const res = await getProfile();
      setState(prev => ({ ...prev, profile: res.data }));
    } catch {
      // profile fetch failed, but user may still be logged in
    } finally {
      setState(prev => ({ ...prev, loading: false }));
    }
  }, []);

  useEffect(() => {
    fetchProfile();
  }, [fetchProfile]);

  const login = useCallback(async (data: LoginRequest) => {
    setState(prev => ({ ...prev, loading: true, error: null }));
    try {
      const res = await apiLogin(data);
      if (res.data.token) {
        
        if (res.data.role !== 'candidate') {
          setState(prev => ({ ...prev, error: '您不是候选人用户，无法登录用户端', loading: false }));
          return false;
        }
        
        setToken(res.data.token);
        setState(prev => ({ ...prev, user: res.data, loading: false }));
        await fetchProfile();
        return true;
      }
      setState(prev => ({ ...prev, error: '用户名或密码错误', loading: false }));
      return false;
    } catch (err) {
      setState(prev => ({ ...prev, error: err instanceof Error ? err.message : '登录失败', loading: false }));
      return false;
    }
  }, [fetchProfile]);

  const register = useCallback(async (data: RegisterRequest) => {
    setState(prev => ({ ...prev, loading: true, error: null }));
    try {
      const res = await apiRegister(data);
      if (res.data.user_id > 0) {
        setState(prev => ({ ...prev, loading: false }));
        return true;
      }
      setState(prev => ({ ...prev, error: '注册失败', loading: false }));
      return false;
    } catch (err) {
      setState(prev => ({ ...prev, error: err instanceof Error ? err.message : '注册失败', loading: false }));
      return false;
    }
  }, []);

  const logout = useCallback(() => {
    removeToken();
    setState({ user: null, profile: null, loading: false, error: null });
  }, []);

  const isProfileComplete = useCallback((): boolean => {
    const p = state.profile;
    if (!p) return false;
    return !!(p.real_name && p.phone && p.education && p.school && p.experience && p.skills);
  }, [state.profile]);

  const refreshProfile = useCallback(async () => {
    if (!getToken()) {
      return;
    }
    try {
      const res = await getProfile();
      setState(prev => ({ ...prev, profile: res.data }));
    } catch {
      // profile fetch failed
    }
  }, []);

  return {
    ...state,
    login,
    register,
    logout,
    refreshProfile,
    isLoggedIn: !!getToken(),
    isProfileComplete: isProfileComplete(),
  };
}