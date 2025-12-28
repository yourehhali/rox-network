"use client";
import { useState, useEffect } from "react";
import { apiFetch } from "../../../../lib/api";
import { usePosts } from "../../../../context/PostsContext";
import EventCard from "../../events/EventCard";
import styles from "./EventsFeed.module.css";

export default function EventsFeed({ groupId }) {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [hasMore, setHasMore] = useState(true);
  const { refreshTrigger } = usePosts();

  const fetchEvents = async (isRefresh = false) => {
    try {
      // If no group ID, do nothing (or handle global events later)
      if (!groupId) {
        setLoading(false);
        return;
      }

      let cursor = 0;
      if (!isRefresh && events.length > 0) {
        cursor = events[events.length - 1].id;
      }

      const res = await apiFetch(`/group-events?group_id=${groupId}&limit=10&cursor=${cursor}`);

      if (res.ok) {
        const data = await res.json();
        const newEvents = Array.isArray(data) ? data : [];
        
        if (newEvents.length < 10) {
            setHasMore(false);
        } else {
            setHasMore(true);
        }

        if (isRefresh) {
            setEvents(newEvents);
        } else {
            setEvents(prev => {
                const existingIds = new Set(prev.map(e => e.id));
                const uniqueNew = newEvents.filter(e => !existingIds.has(e.id));
                return [...prev, ...uniqueNew];
            });
        }
      }
    } catch (e) {
      console.error("Failed to fetch events", e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents(true);
  }, [groupId, refreshTrigger]);

  const handleEventResponse = async (eventId, status) => {
    try {
      await apiFetch("/event-invite-response", {
        method: "POST",
        body: JSON.stringify({ event_id: eventId, status }),
      });
      // Refresh events to show updated status (simple reload)
      fetchEvents(true);
    } catch (e) {
      // Failed to respond
    }
  };

  if (loading && events.length === 0) {
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
        {hasMore && (
            <button onClick={() => fetchEvents(false)} className={styles.loadMoreBtn}>
                Load More
            </button>
        )}
      </div>
    </div>
  );
}
