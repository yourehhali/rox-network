"use client";
import { useState } from "react";
import { apiFetch } from "../../../../lib/api";
import styles from "./CreatePostModal.module.css";
import { X } from "lucide-react";

export default function CreatePostModal({ onClose, onSuccess }) {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [image, setImage] = useState(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

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
    if (image) {
      formData.append("image", image);
    }

    try {
      const res = await apiFetch("/Createpost", {
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
          <h3>Create New Post</h3>
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
            <input
              type="file"
              accept="image/*"
              onChange={handleImageChange}
              className={styles.fileInput}
            />
          </div>

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
