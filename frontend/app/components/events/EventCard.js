"use client";
import styles from "./EventCard.module.css";

export default function EventCard({ event, onResponse }) {
  const eventDate = new Date(event.event_time);
  const isEnded = eventDate < new Date();

  return (
    <div className={styles.eventCard}>
      <h3>{event.title}</h3>
      <p>{event.description}</p>
      <p className={styles.date}>Date: {eventDate.toLocaleString()}</p>

      {isEnded ? (
        <div className={styles.eventActions}>
          <button
            disabled
            style={{ backgroundColor: "#ccc", cursor: "not-allowed" }}
          >
            Event Ended
          </button>
        </div>
      ) : (
        <div className={styles.eventActions}>
          <button
            onClick={() => onResponse(event.id, "going")}
            disabled={event.user_status === "going"}
          >
            {event.user_status === "going" ? "Going ✓" : "Going"}
          </button>
          <button
            onClick={() => onResponse(event.id, "not_going")}
            disabled={event.user_status === "not_going"}
          >
            {event.user_status === "not_going" ? "Not Going ✓" : "Not Going"}
          </button>
        </div>
      )}
    </div>
  );
}
