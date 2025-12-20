"use client";

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiFetch } from '../../lib/api';
import style from './register.module.css';
import Link from 'next/link';

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    username: '',
    first_name: '',
    last_name: '',
    email: '',
    password: '',
    confirm_password: '',
    date_of_birth: '',
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

  const calculateAge = (dob) => {
    const birthDate = new Date(dob);
    const today = new Date();
    let age = today.getFullYear() - birthDate.getFullYear();
    const m = today.getMonth() - birthDate.getMonth();
    if (m < 0 || (m === 0 && today.getDate() < birthDate.getDate())) {
      age--;
    }
    return age;
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
    
    if (formData.date_of_birth) {
        const age = calculateAge(formData.date_of_birth);
        data.append('age', age.toString());
    } else {
        data.append('age', '0'); 
    }

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
      <div className={style.formCard}>
        <h2 className={style.title}>Create Account</h2>
        <p className={style.subtitle}>
          Already have an account? <Link href="/login" className={style.link}>Sign in</Link>
        </p>
        
        {error && (
          <div className={style.error}>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className={style.inputGroup}>
            <input name="first_name" required className={style.input} placeholder="First Name" onChange={handleChange} />
          </div>
          <div className={style.inputGroup}>
            <input name="last_name" required className={style.input} placeholder="Last Name" onChange={handleChange} />
          </div>
          <div className={style.inputGroup}>
            <input name="nickname" className={style.input} placeholder="Nickname" onChange={handleChange} />
          </div>
          <div className={style.inputGroup}>
            <input name="username" required className={style.input} placeholder="Username" onChange={handleChange} />
          </div>
          <div className={style.inputGroup}>
            <input name="email" type="email" required className={style.input} placeholder="Email" onChange={handleChange} />
          </div>
          <div className={style.inputGroup}>
            <input name="password" type="password" required className={style.input} placeholder="Password" onChange={handleChange} />
          </div>
          <div className={style.inputGroup}>
            <input name="confirm_password" type="password" required className={style.input} placeholder="Confirm Password" onChange={handleChange} />
          </div>
          <div className={style.inputGroup}>
            <label className={style.label}>Date of Birth</label>
            <input name="date_of_birth" type="date" required className={style.input} onChange={handleChange} />
          </div>
          <div className={style.inputGroup}>
            <select name="gender" required className={style.select} onChange={handleChange} defaultValue="">
              <option value="" disabled>Select Gender</option>
              <option value="Male">Male</option>
              <option value="Female">Female</option>
            </select>
          </div>
          <div className={style.inputGroup}>
            <textarea name="about_me" className={style.textarea} placeholder="About Me" onChange={handleChange}></textarea>
          </div>
          <div className={style.inputGroup}>
            <label className={style.label}>Avatar</label>
            <input type="file" className={style.input} onChange={handleFileChange} accept="image/*" />
          </div>

          <button type="submit" className={style.button}>
            Register
          </button>
        </form>
      </div>
    </div>
  );
}
