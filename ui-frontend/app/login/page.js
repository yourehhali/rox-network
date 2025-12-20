"use client";

import { useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { apiFetch } from '../../lib/api';
import style from './login.module.css';
import Link from 'next/link';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const res = await apiFetch('/login', {
        method: 'POST',
        body: JSON.stringify({ username: email, password }),
      });

      if (!res.ok) {
        throw new Error('Login failed');
      }
      
      const loginData = await res.json();
      
      if (loginData.user_id) {
        const profileRes = await apiFetch(`/profile?id=${loginData.user_id}`);
        if (profileRes.ok) {
          const profileData = await profileRes.json();
          const user = profileData.user || profileData;
          login(user);
        } else {
          setError("Login successful but failed to fetch profile");
        }
      } else {
        setError("Login successful but no user ID returned");
      }
    } catch (err) {
      setError('Invalid credentials');
    }
  };

  return (
    <div className={style.container}>
      <div className={style.formCard}>
        <h2 className={style.title}>Sign in</h2>
        <p className={style.subtitle}>
          Or <Link href="/register" className={style.link}>create a new account</Link>
        </p>
        
        {error && (
          <div className={style.error}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className={style.inputGroup}>
            <input
              type="email"
              required
              className={style.input}
              placeholder="Email address"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div className={style.inputGroup}>
            <input
              type="password"
              required
              className={style.input}
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          <button type="submit" className={style.button}>
            Sign in
          </button>
        </form>
      </div>
    </div>
  );
}
