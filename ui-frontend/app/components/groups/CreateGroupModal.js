"use client";
import { useState } from "react";
import { apiFetch } from "../../../lib/api";
import styles from "./CreateGroupModal.module.css";

export default function CreateGroupModal({ onSuccess, isOpen, onClose }) {
  const [internalOpen, setInternalOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [error, setError] = useState("");

  const isControlled = typeof isOpen !== "undefined";
  const open = isControlled ? isOpen : internalOpen;

  const handleClose = () => {
    if (isControlled) {
      if (onClose) onClose();
    } else {
      setInternalOpen(false);
    }
  };

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
        if (isControlled) {
             if (onClose) onClose();
        } else {
             setInternalOpen(false);
        }
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
      {!isControlled && (
        <button className={styles.createBtn} onClick={() => setInternalOpen(true)}>
          + Create Group
        </button>
      )}

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
                minLength={3}
                maxLength={30}
              />
              <textarea
                placeholder="Description"
                value={desc}
                onChange={(e) => setDesc(e.target.value)}
                required
                className={styles.textarea}
                maxLength={100}
              />
              <div className={styles.buttons}>
                <button
                  type="button"
                  className={styles.cancelBtn}
                  onClick={handleClose}
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
