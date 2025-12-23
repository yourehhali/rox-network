"use client";
import { useState, useEffect } from "react";
import { apiFetch, toMediaUrl } from "../../../../lib/api";
import styles from "./CreatePostModal.module.css";
import { X } from "lucide-react";
import { useAuth } from "../../../../context/AuthContext";

export default function CreatePostModal({ onClose, onSuccess, groupId = null }) {
  const { user } = useAuth();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [image, setImage] = useState(null);
  const [privacy, setPrivacy] = useState("0"); // 0: Public, 1: Private, 2: Almost Private
  const [allowedUsers, setAllowedUsers] = useState(new Set());
  const [followers, setFollowers] = useState([]);
  const [loadingFollowers, setLoadingFollowers] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!groupId && privacy === "2" && user?.id) {
      fetchFollowers();
    }
  }, [privacy, user, groupId]);

  const fetchFollowers = async () => {
    setLoadingFollowers(true);
    try {
      const res = await apiFetch(`/followers?id=${user.id}`);
      if (res.ok) {
        const data = await res.json();
        setFollowers(data || []);
      }
    } catch (err) {
      console.error("Failed to fetch followers", err);
    } finally {
      setLoadingFollowers(false);
    }
  };

  const handleAllowedUserToggle = (userId) => {
    const newAllowed = new Set(allowedUsers);
    if (newAllowed.has(userId)) {
      newAllowed.delete(userId);
    } else {
      newAllowed.add(userId);
    }
    setAllowedUsers(newAllowed);
  };

  const handleImageChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      setImage(e.target.files[0]);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    if (!title || !content) {
      setError("Please fill in title and content");
      setLoading(false);
      return;
    }

    const formData = new FormData();
    formData.append("title", title);
    formData.append("content", content);
    formData.append("privacy", privacy);
    
    if (privacy === "2" && !groupId) {
        const allowedList = Array.from(allowedUsers).join(",");
        formData.append("allowed_users", allowedList);
    }

    if (groupId) {
      formData.append("group_id", groupId);
    }

    if (image) {
      formData.append("image", image);
    }

    try {
      const endpoint = groupId ? "/CreateGroupPost" : "/Createpost";
      const res = await apiFetch(endpoint, {
        method: "POST",
        body: formData,
      });

      if (res.ok) {
        if (onSuccess) onSuccess();
        if (onClose) onClose();
      } else {
        // Try to parse JSON error first, then text
        try {
            const errJson = await res.json();
            setError(errJson.message || "Failed to create post");
        } catch {
            const errText = await res.text();
            setError(errText || "Failed to create post");
        }
      }
    } catch (err) {
      setError("An error occurred");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.overlay}>
      <div className={styles.modal}>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "1rem" }}>
          <h3>{groupId ? "Create Group Post" : "Create New Post"}</h3>
          <button onClick={onClose} style={{ background: "none", border: "none", color: "white", cursor: "pointer" }}>
            <X />
          </button>
        </div>
        
        {error && <p className={styles.error}>{error}</p>}
        
        <form onSubmit={handleSubmit} className={styles.form}>
          <input
            type="text"
            placeholder="Post Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            className={styles.input}
          />
          <textarea
            placeholder="What's on your mind?"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            required
            className={styles.textarea}
          />
          
          <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}>
            <label style={{ color: "#aaa", fontSize: "0.9rem" }}>Upload Image (optional)</label>
            <div className={styles.fileInputWrapper}>
              <input
                type="file"
                id="post-image-upload"
                accept="image/*"
                onChange={handleImageChange}
                className={styles.fileInput}
              />
              <label htmlFor="post-image-upload" className={styles.fileLabel}>
                {image ? (
                  <div className={styles.filePreview}>
                    <span className={styles.fileName}>{image.name}</span>
                    <span className={styles.changeFile}>Change</span>
                  </div>
                ) : (
                  <span className={styles.placeholder}>Click to upload image</span>
                )}
              </label>
            </div>
          </div>

            {!groupId && (
              <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}>
                <label style={{ color: "#aaa", fontSize: "0.9rem" }}>Privacy</label>
                <select
                  value={privacy}
                  onChange={(e) => setPrivacy(e.target.value)}
                  className={styles.select}
                >
                  <option value="0">Public</option>
                  <option value="1">Private (Followers only)</option>
                  <option value="2">Almost Private (Select users)</option>
                </select>
              </div>
            )}

            {!groupId && privacy === "2" && (
              <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}>
              <label style={{ color: "#aaa", fontSize: "0.9rem" }}>Select Allowed Users</label>
              {loadingFollowers ? (
                <p style={{ color: "#aaa" }}>Loading followers...</p>
              ) : (
                <div className={styles.followersList}>
                  {followers.length === 0 ? (
                    <p style={{ color: "#aaa", padding: "0.5rem" }}>No followers found.</p>
                  ) : (
                    followers.map((follower) => (
                      <label key={follower.id} className={styles.followerItem}>
                        <input
                          type="checkbox"
                          className={styles.followerCheckbox}
                          checked={allowedUsers.has(String(follower.id))}
                          onChange={() => handleAllowedUserToggle(String(follower.id))}
                        />
                        <span style={{ color: "white" }}>
                            {follower.first_name} {follower.last_name} (@{follower.username})
                        </span>
                      </label>
                    ))
                  )}
                </div>
              )}
            </div>
          )}

          <div className={styles.buttons}>
            <button
              type="button"
              className={styles.cancelBtn}
              onClick={onClose}
              disabled={loading}
            >
              Cancel
            </button>
            <button type="submit" className={styles.submitBtn} disabled={loading}>
              {loading ? "Posting..." : "Post"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
