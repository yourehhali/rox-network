"use client";
import { useEffect, useState } from "react";
import style from "./userListModal.module.css";
import { apiFetch, toMediaUrl } from "../../../lib/api";
import { X } from "lucide-react";
import Link from "next/link";

export default function UserListModal({ type, userId, onClose }) {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const endpoint = type === "followers" ? "/followers" : "/following";
        const res = await apiFetch(`${endpoint}?id=${userId}`);
        if (res.ok) {
          const data = await res.json();
          setUsers(data || []);
        }
      } catch (e) {
        console.error("Failed to fetch users", e);
      } finally {
        setLoading(false);
      }
    };
    if (userId) fetchUsers();
  }, [type, userId]);

  return (
    <div className={style.modalOverlay} onClick={onClose}>
      <div className={style.modalContent} onClick={(e) => e.stopPropagation()}>
        <div className={style.header}>
          <h2>{type === "followers" ? "Followers" : "Following"}</h2>
          <button className={style.closeButton} onClick={onClose}>
            <X size={24} />
          </button>
        </div>
        <div className={style.list}>
          {loading ? (
            <div className={style.loading}>Loading...</div>
          ) : users.length === 0 ? (
            <div className={style.empty}>No users found.</div>
          ) : (
            users.map((u) => (
              <Link href={`/profile?id=${u.id}`} key={u.id} className={style.userItem} onClick={onClose}>
                <img
                  src={toMediaUrl(u.image_path) || "/default-avatar.png"}
                  alt={u.username}
                  className={style.avatar}
                />
                <div className={style.userInfo}>
                  <div className={style.fullName}>{u.full_name}</div>
                  <div className={style.username}>@{u.username}</div>
                </div>
              </Link>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
