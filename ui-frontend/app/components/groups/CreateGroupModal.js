"use client";
import { useState } from "react";
import { apiFetch } from "../../../lib/api";
import styles from "./CreateGroupModal.module.css";

export default function CreateGroupModal({ onSuccess }) {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    try {
      const res = await apiFetch("/Create_Group", {
        method: "POST",
        body: JSON.stringify({ title, description: desc }),
      });

      if (res.ok) {
        setTitle("");
        setDesc("");
        setOpen(false);
        if (onSuccess) onSuccess();
      } else {
        const errData = await res.text();
        setError(errData || "Failed to create group");
      }
    } catch (err) {
      setError("An error occurred");
    }
  };

  return (
    <>
      <button className={styles.createBtn} onClick={() => setOpen(true)}>
        + Create Group
      </button>

      {open && (
        <div className={styles.overlay}>
          <div className={styles.modal}>
            <h3>Create Group</h3>
            {error && <p className={styles.error}>{error}</p>}
            <form onSubmit={handleSubmit} className={styles.form}>
              <input
                type="text"
                placeholder="Title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
                className={styles.input}
              />
              <textarea
                placeholder="Description"
                value={desc}
                onChange={(e) => setDesc(e.target.value)}
                required
                className={styles.textarea}
              />
              <div className={styles.buttons}>
                <button
                  type="button"
                  className={styles.cancelBtn}
                  onClick={() => setOpen(false)}
                >
                  Cancel
                </button>
                <button type="submit" className={styles.submitBtn}>
                  Create
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}
