"use client";

import style from "./friendsSuggestions.module.css";
import { useRouter } from "next/navigation";
import { useEffect, useState, useRef } from "react";
import { apiFetch, toMediaUrl } from "../../../../../lib/api";
import { Search } from "../../../icons";

export default function FriendsSuggestions() {
  const router = useRouter();
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  
  const [isExpanded, setIsExpanded] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const [debouncedSearchTerm, setDebouncedSearchTerm] = useState("");

  const containerRef = useRef(null);

  useEffect(() => {
    const handler = setTimeout(() => {
        setDebouncedSearchTerm(searchTerm);
    }, 500); 
    return () => clearTimeout(handler);
  }, [searchTerm]);

  
  useEffect(() => {
      function handleClickOutside(event) {
          if (containerRef.current && !containerRef.current.contains(event.target)) {
              setIsExpanded(false);
          }
      }
      document.addEventListener("mousedown", handleClickOutside);
      return () => {
          document.removeEventListener("mousedown", handleClickOutside);
      };
  }, [containerRef]);

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

  const filteredSuggestions = suggestions.filter((user) =>
    user.name.toLowerCase().includes(debouncedSearchTerm.toLowerCase())
  );

  const buildSuggestions = filteredSuggestions.map((user) => (
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
        </div>
      </div>
    </div>
  ));

  return (
    <div className={style.containerWrapper} ref={containerRef}>
        <div 
            className={style.searchIconWrapper} 
            onClick={() => setIsExpanded(!isExpanded)}
            title="Search Friends"
        >
            <Search />
        </div>

        {isExpanded && (
            <div className={style.expandedPanel}>
                <h2 style={{ fontSize: '1.2rem', margin: 0 }}>Friends Suggestions</h2>
                <input
                    type="text"
                    placeholder="Search users..."
                    className={style.searchInput}
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    autoFocus
                />
                <div className={style.data}>
                    {loading && <div>Loading...</div>}
                    {error && <div>{error}</div>}
                    {!loading && !error && (
                        buildSuggestions.length > 0 ? buildSuggestions : <div className={style.noResults}>No users found</div>
                    )}
                </div>
            </div>
        )}
    </div>
  );
}
