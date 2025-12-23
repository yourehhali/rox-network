"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiFetch } from '../../lib/api';
import style from './register.module.css';
import Link from 'next/link';
import { UserPlus } from 'lucide-react';

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    username: '',
    first_name: '',
    last_name: '',
    email: '',
    password: '',
    confirm_password: '',
    age: '',
    gender: '',
    nickname: '',
    about_me: '',
  });
  const [avatar, setAvatar] = useState(null);
  const [error, setError] = useState('');
  const router = useRouter();

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleFileChange = (e) => {
    if (e.target.files?.[0]) {
      setAvatar(e.target.files[0]);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (formData.password !== formData.confirm_password) {
      setError("Passwords do not match");
      return;
    }

    if (!formData.gender) {
        setError("Please select a gender");
        return;
    }

    const data = new FormData();
    data.append('username', formData.username);
    data.append('first_name', formData.first_name);
    data.append('last_name', formData.last_name);
    data.append('email', formData.email);
    data.append('password', formData.password);
    data.append('confirmpassword', formData.confirm_password);
    data.append('gender', formData.gender);
    data.append('nickname', formData.nickname);
    data.append('about', formData.about_me);
    data.append('age', formData.age || '0');

    if (avatar) {
      data.append('avatar', avatar);
    }

    try {
      const res = await apiFetch('/register', {
        method: 'POST',
        body: data,
      });

      if (!res.ok) {
        const errData = await res.json();
        throw new Error(errData.message || 'Registration failed');
      }

      router.push('/login');
    } catch (err) {
      setError(err.message || 'An unexpected error occurred');
    }
  };

  return (
    <div className={style.container}>
      <div className={style.formModalContainer}>
        <div className={style.icon}>
          <UserPlus size={32} />
        </div>
        <div className={style.texts}>
          <div className={style.heading}>
            <h2>Welcome!</h2>
          </div>
          <div className={style.subHeading}>
            <p>Create your 01Social account</p>
          </div>
        </div>

        {error && (
          <div className={style.error}>
            {error}
          </div>
        )}

        <div className={style.formContainer}>
          <form onSubmit={handleSubmit} className={style.form}>
            <div className={style.inputContainer}>
              <label htmlFor="first_name" className={style.label}>First Name</label>
              <input
                id="first_name"
                type="text"
                name="first_name"
                placeholder="first name"
                className={style.input}
                value={formData.first_name}
                onChange={handleChange}
                required
              />
            </div>
            <div className={style.inputContainer}>
              <label htmlFor="last_name" className={style.label}>Last Name</label>
              <input
                id="last_name"
                type="text"
                name="last_name"
                placeholder="last name"
                className={style.input}
                value={formData.last_name}
                onChange={handleChange}
                required
              />
            </div>
            <div className={style.inputContainer}>
              <label htmlFor="username" className={style.label}>Username</label>
              <input
                id="username"
                type="text"
                name="username"
                placeholder="username"
                className={style.input}
                value={formData.username}
                onChange={handleChange}
                required
              />
            </div>
            <div className={style.inputContainer}>
              <label htmlFor="gender" className={style.label}>Gender</label>
              <select
                id="gender"
                name="gender"
                className={style.select}
                value={formData.gender}
                onChange={handleChange}
                required
              >
                <option value="">Select Gender</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
              </select>
            </div>
            <div className={style.inputContainerLarge}>
              <label htmlFor="email" className={style.label}>Email</label>
              <input
                id="email"
                type="email"
                name="email"
                placeholder="email"
                className={style.input}
                value={formData.email}
                onChange={handleChange}
                required
              />
            </div>
            <div className={style.inputContainer}>
              <label htmlFor="password" className={style.label}>Password</label>
              <input
                id="password"
                type="password"
                name="password"
                placeholder="password"
                className={style.input}
                value={formData.password}
                onChange={handleChange}
                required
              />
            </div>
            <div className={style.inputContainer}>
              <label htmlFor="confirm_password" className={style.label}>Confirm Password</label>
              <input
                id="confirm_password"
                type="password"
                name="confirm_password"
                placeholder="confirm password"
                className={style.input}
                value={formData.confirm_password}
                onChange={handleChange}
                required
              />
            </div>
            <div className={style.inputContainer}>
              <label htmlFor="nickname" className={style.label}>Nickname (Optional)</label>
              <input
                id="nickname"
                type="text"
                name="nickname"
                placeholder="nickname"
                className={style.input}
                value={formData.nickname}
                onChange={handleChange}
              />
            </div>
            <div className={style.inputContainer}>
              <label htmlFor="age" className={style.label}>Age</label>
              <input
                id="age"
                type="number"
                name="age"
                placeholder="age"
                className={style.input}
                value={formData.age}
                onChange={handleChange}
                required
              />
            </div>
            <div className={style.inputContainerLarge}>
              <label htmlFor="about_me" className={style.label}>About me (Optional)</label>
              <textarea
                id="about_me"
                name="about_me"
                placeholder="about me"
                className={style.textarea}
                value={formData.about_me}
                onChange={handleChange}
              />
            </div>
            <div className={style.inputContainerLarge}>
              <div className={style.fileInputWrapper}>
                <input
                  type="file"
                  id="avatar-upload"
                  accept="image/*"
                  onChange={handleFileChange}
                  className={style.fileInput}
                />
                <label htmlFor="avatar-upload" className={style.fileLabel}>
                  {avatar ? (
                    <div className={style.filePreview}>
                      <span className={style.fileName}>{avatar.name}</span>
                      <span className={style.changeFile}>Change</span>
                    </div>
                  ) : (
                    <span className={style.placeholder}>Choose Avatar (Optional)</span>
                  )}
                </label>
              </div>
            </div>
            <div className={style.buttons}>
              <Link href="/login" className={style.secondAction}>
                Login
              </Link>
              <input className={style.mainAction} type="submit" value="Register" />
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
