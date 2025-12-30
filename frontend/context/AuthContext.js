"use client";

import { createContext, useContext, useState, useEffect } from 'react';
import { apiFetch } from '../lib/api';
import { useRouter } from 'next/navigation';

const AuthContext = createContext(undefined);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  const checkSession = async () => {
    try {
      const res = await apiFetch('/checkstate');
      if (res.ok) {
        const data = await res.json();
        if (data.authenticated) {
          const profileRes = await apiFetch('/profile?id=' + data.userID);
          if (profileRes.ok) {
             const profileData = await profileRes.json();
             setUser(profileData.user || profileData); 
          }
        }
      }
    } catch (err) {
      console.error("Session check failed", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkSession();
  }, []);

  useEffect(() => {
    const handleUnauthorized = () => {
      setUser(null);
      router.push('/login');
    };

    window.addEventListener('auth:unauthorized', handleUnauthorized);
    return () => window.removeEventListener('auth:unauthorized', handleUnauthorized);
  }, [router]);

  const login = (userData) => {
    setUser(userData);
    router.push('/');
  };

  const logout = async () => {
    await apiFetch('/logout', { method: 'POST' });
    setUser(null);
    router.push('/login');
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, checkSession }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
