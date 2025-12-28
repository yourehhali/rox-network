"use client";
import { useState, useEffect } from "react";
import { useAuth } from "@/context/AuthContext";
import { useToast } from "@/context/ToastContext";
import { apiFetch, toMediaUrl } from "@/lib/api";
import styles from "./page.module.css";

export default function SettingsPage() {
  const { user, setUser } = useAuth(); // Assuming setUser can update context if needed
  const { addToast } = useToast();
  const [loading, setLoading] = useState(true);
  const [formData, setFormData] = useState({
    first_name: "",
    last_name: "",
    nickname: "",
    about: "",
    is_private: false,
    avatar: null, // File object
    avatarUrl: "", // For preview
  });

  useEffect(() => {
    if (user) {
      fetchProfile();
    }
  }, [user]);

  const fetchProfile = async () => {
    try {
      const res = await apiFetch(`/profile?id=${user.id}`);
      if (res.ok) {
        const data = await res.json();
        const u = data.user;
        setFormData({
          first_name: u.first_name || "",
          last_name: u.last_name || "",
          nickname: u.nickname || "",
          about: u.about || "",
          is_private: u.is_private === 1,
          avatar: null,
          avatarUrl: u.avatar ? toMediaUrl(u.avatar) : "/avatar.png",
        });
      }
    } catch (e) {
      console.error("Failed to load profile", e);
      addToast("Failed to load profile settings", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value,
    }));
  };

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setFormData((prev) => ({
        ...prev,
        avatar: file,
        avatarUrl: URL.createObjectURL(file),
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const data = new FormData();
      data.append("first_name", formData.first_name);
      data.append("last_name", formData.last_name);
      data.append("nickname", formData.nickname);
      data.append("about", formData.about);
      data.append("is_private", formData.is_private ? "1" : "0");
      if (formData.avatar) {
        data.append("avatar", formData.avatar);
      }

      const res = await apiFetch("/update-profile", {
        method: "POST",
        body: data, // apiFetch handles headers for FormData usually, or we need to omit Content-Type
      });

      if (res.ok) {
        addToast("Profile updated successfully", "success");
        // Optionally update auth context if it holds these details
        // setUser({...user, ...updatedFields}) - requires parsing response or reusing form data
      } else {
        const err = await res.text();
        // Try to parse JSON error if possible
        try {
            const jsonErr = JSON.parse(err);
            addToast(jsonErr.message || "Update failed", "error");
        } catch {
            addToast("Update failed: " + err, "error");
        }
      }
    } catch (e) {
      console.error("Update error", e);
      addToast("An error occurred", "error");
    }
  };

  return (
    <div className={styles.middle}>
      <div className={styles.header}>
        <h1 className={styles.title}>Account Settings</h1>
      </div>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <form className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.inputGroup} style={{ alignItems: "center" }}>
            <img
              src={formData.avatarUrl || "/avatar.png"}
              alt="Avatar Preview"
              className={styles.avatarPreview}
            />
            <div className={styles.fileInputWrapper}>
              <input
                type="file"
                id="avatar-upload"
                onChange={handleFileChange}
                accept="image/*"
                className={styles.fileInput}
              />
              <label htmlFor="avatar-upload" className={styles.fileLabel}>
                <span className={styles.changeFile}>Change Avatar</span>
              </label>
            </div>
          </div>

          <div className={styles.inputGroup}>
            <label className={styles.label}>First Name</label>
            <input
              name="first_name"
              value={formData.first_name}
              onChange={handleChange}
              className={styles.input}
              placeholder="First Name"
              maxLength={50}
            />
          </div>

          <div className={styles.inputGroup}>
            <label className={styles.label}>Last Name</label>
            <input
              name="last_name"
              value={formData.last_name}
              onChange={handleChange}
              className={styles.input}
              placeholder="Last Name"
              maxLength={50}
            />
          </div>

          <div className={styles.inputGroup}>
            <label className={styles.label}>Nickname</label>
            <input
              name="nickname"
              value={formData.nickname}
              onChange={handleChange}
              className={styles.input}
              placeholder="Nickname"
              maxLength={50}
            />
          </div>

          <div className={styles.inputGroup}>
            <label className={styles.label}>About Me</label>
            <textarea
              name="about"
              value={formData.about}
              onChange={handleChange}
              className={styles.textarea}
              placeholder="Tell us about yourself..."
              maxLength={500}
            />
          </div>

          <div className={styles.checkboxGroup}>
            <input
              type="checkbox"
              name="is_private"
              checked={formData.is_private}
              onChange={handleChange}
              className={styles.checkbox}
              id="is_private"
            />
            <label htmlFor="is_private" className={styles.label} style={{ cursor: "pointer" }}>
              Private Account
            </label>
          </div>

          <button type="submit" className={styles.button}>
            Save Changes
          </button>
        </form>
      )}
    </div>
  );
}
