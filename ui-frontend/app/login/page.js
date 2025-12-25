"use client";

import { useState } from 'react';
import { useAuth } from '../../context/AuthContext';
import { apiFetch } from '../../lib/api';
import style from './login.module.css';
import Link from 'next/link';
import { LogIn } from '../components/icons';

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
      <div className={style.formModalContainer}>
        <div className={style.icon}>
          <LogIn size={32} />
        </div>
        <div className={style.texts}>
          <div className={style.heading}>
            <h2>Welcome Back</h2>
          </div>
          <div className={style.subHeading}>
            <p>Sign in to your 01Social account</p>
          </div>
        </div>

        {error && (
          <div className={style.error}>
            {error}
          </div>
        )}

        <div className={style.formContainer}>
          <form onSubmit={handleSubmit}>
            <div className={style.inputContainerLarge}>
              <label htmlFor="email" className={style.label}>Email</label>
              <input
                id="email"
                type="text" 
                name="email"
                placeholder="email"
                className={style.input}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className={style.inputContainerLarge}>
              <label htmlFor="password" className={style.label}>Password</label>
              <input
                id="password"
                type="password"
                name="password"
                placeholder="password"
                className={style.input}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <div className={style.buttons}>
              <Link href="/register" className={style.secondAction}>
                Register
              </Link>
              <input className={style.mainAction} type="submit" value="Login" />
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
