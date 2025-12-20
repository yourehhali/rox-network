"use client";

import style from "./friendsSuggestions.module.css";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { apiFetch, toMediaUrl } from "../../../../../lib/api";

export default function FriendsSuggestions() {
  const router = useRouter();
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const loadSuggestions = async () => {
      try {
        const res = await apiFetch("/users-sug");
        if (!res.ok) {
          throw new Error("Failed to load suggestions");
        }
        const data = await res.json();
        const mapped = (data || []).map((u) => {
          const avatar = toMediaUrl(u.image_path) || "/default-avatar.png";
          return {
            id: u.user_id,
            name: u.full_name || u.username,
            avatar,
          };
        });
        setSuggestions(mapped);
      } catch (e) {
        setError(e.message || "Unexpected error");
      } finally {
        setLoading(false);
      }
    };
    loadSuggestions();
  }, []);

  const navigateTo = (endpoint) => {
    router.push(endpoint);
  };

  const handleFollow = async (userId) => {
    try {
      const res = await apiFetch("/toggle-follow", {
        method: "POST",
        body: JSON.stringify({ followed_id: userId }),
      });
      if (!res.ok) {
        return;
      }
      // Optionally refresh suggestions to remove followed user
      setSuggestions((prev) => prev.filter((u) => u.id !== userId));
    } catch (e) {
      // swallow errors for now
    }
  };

  if (loading) {
    return (
      <div className={style.FriendsSuggestions}>
        <h2>Friends Suggestions</h2>
        <div className={style.data}>Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={style.FriendsSuggestions}>
        <h2>Friends Suggestions</h2>
        <div className={style.data}>{error}</div>
      </div>
    );
  }

  const buildSuggestions = suggestions.map((user) => (
    <div key={user.id} className={style.FriendsSuggestion}>
      <div className={style.container}>
        <div className={style.left}>
          <img src={user.avatar} alt={user.name} />
          <span>{user.name}</span>
        </div>
        <div className={style.right}>
          <button
            className={style.profileButton}
            onClick={() => navigateTo(`/profile?id=${user.id}`)}
          >
            Profile
          </button>
          <button
            className={style.followButton}
            onClick={() => handleFollow(user.id)}
          >
            Follow
          </button>
        </div>
      </div>
    </div>
  ));

  return (
    <div className={style.FriendsSuggestions}>
      <h2>Friends Suggestions</h2>
      <div className={style.data}>{buildSuggestions}</div>
    </div>
  );
}
