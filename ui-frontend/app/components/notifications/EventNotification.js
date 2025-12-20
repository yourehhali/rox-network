"use client";
import { useState } from "react";
import style from "./followRequest.module.css";
import { apiFetch } from "../../../lib/api";

export default function EventNotification({ notification, onHandled }) {
  const [status, setStatus] = useState("pending");

  const handleAction = async (action) => {
    try {
      const res = await apiFetch("/event-invite-response", {
        method: "POST",
        body: JSON.stringify({
          event_id: notification.event_id,
          status: action, // 'going' or 'not_going'
        }),
      });
      if (res.ok) {
        setStatus(action);
        if (onHandled) onHandled(action);
      }
    } catch (_) {}
  };

  const sender = notification.sender || {};
  
  return (
    <div className={style.item}>
      <div className={style.info}>
        <div className={style.name}>
          <p className={style.message}>{notification.message}</p>
          <p className={style.username}>By {sender.username}</p>
        </div>
      </div>
      {status === "pending" && (
        <div className={style.actions}>
          <button
            className={`${style.btn} ${style.accept}`}
            onClick={() => handleAction("going")}
          >
            Going
          </button>
          <button
            className={`${style.btn} ${style.reject}`}
            onClick={() => handleAction("not_going")}
          >
            Ignore
          </button>
        </div>
      )}
    </div>
  );
}
