"use client";
import { useState } from "react";
import { apiFetch } from "../../../lib/api";
import styles from "../groups/CreateGroupModal.module.css"; // Reuse same styles

export default function CreateEventModal({ groupId, onSuccess, onClose }) {
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [eventTime, setEventTime] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    if (!title || !desc || !eventTime) {
      setError("Please fill all fields");
      return;
    }

    try {
      const res = await apiFetch("/create_event", {
        method: "POST",
        body: JSON.stringify({
          group_id: parseInt(groupId, 10),
          title,
          description: desc,
          event_time: eventTime,
        }),
      });

      if (res.ok) {
        if (onSuccess) onSuccess();
        if (onClose) onClose();
      } else {
        const errData = await res.text();
        setError(errData || "Failed to create event");
      }
    } catch (err) {
      setError("An error occurred");
    }
  };

  return (
    <div className={styles.overlay}>
      <div className={styles.modal}>
        <h3>Create Event</h3>
        {error && <p className={styles.error}>{error}</p>}
        <form onSubmit={handleSubmit} className={styles.form}>
          <input
            type="text"
            placeholder="Event Title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            className={styles.input}
            minLength={3}
            maxLength={50}
          />
          <textarea
            placeholder="Description"
            value={desc}
            onChange={(e) => setDesc(e.target.value)}
            required
            className={styles.textarea}
            maxLength={500}
          />
          <input
            type="datetime-local"
            value={eventTime}
            onChange={(e) => setEventTime(e.target.value)}
            required
            min={new Date(Date.now() + 86400000).toISOString().slice(0, 16)}
            className={styles.input}
          />
          <div className={styles.buttons}>
            <button
              type="button"
              className={styles.cancelBtn}
              onClick={onClose}
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
  );
}
