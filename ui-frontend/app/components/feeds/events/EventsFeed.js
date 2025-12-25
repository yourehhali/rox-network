"use client";
import { useState, useEffect } from "react";
import { apiFetch } from "../../../../lib/api";
import { usePosts } from "../../../../context/PostsContext";
import EventCard from "../../events/EventCard";
import styles from "./EventsFeed.module.css";

export default function EventsFeed({ groupId }) {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const { refreshTrigger } = usePosts();

  const fetchEvents = async () => {
    try {
      let res;
      if (groupId) {
        res = await apiFetch(`/group-events?group_id=${groupId}`);
      } else {
        // Fallback or future implementation for global events
        // For now, if no groupId, maybe empty or different endpoint
        setLoading(false);
        return;
      }

      if (res.ok) {
        const data = await res.json();
        setEvents(Array.isArray(data) ? data : []);
      }
    } catch (e) {
      console.error("Failed to fetch events", e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents();
  }, [groupId, refreshTrigger]);

  const handleEventResponse = async (eventId, status) => {
    try {
      await apiFetch("/event-invite-response", {
        method: "POST",
        body: JSON.stringify({ event_id: eventId, status }),
      });
      // Refresh events to show updated status
      fetchEvents();
    } catch (e) {
      // Failed to respond
    }
  };

  if (loading) {
    return <div>Loading events...</div>;
  }

  return (
    <div className={styles.eventsFeed}>
      <div className={styles.eventsFeedContainer}>
        {events.length === 0 ? (
          <p>No events found.</p>
        ) : (
          events.map((ev) => (
            <EventCard
              key={ev.id}
              event={ev}
              onResponse={handleEventResponse}
            />
          ))
        )}
      </div>
    </div>
  );
}
