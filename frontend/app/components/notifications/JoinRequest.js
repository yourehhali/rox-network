"use client";
import { useState } from "react";
import style from "./followRequest.module.css";
import { apiFetch, toMediaUrl } from "../../../lib/api";

export default function JoinRequest({ notification, onHandled }) {
  const [status, setStatus] = useState("pending");

  const handleAction = async (action) => {
    try {
      const res = await apiFetch("/accept-reject-join", {
        method: "POST",
        body: JSON.stringify({
          notif_id: notification.id,
          group_id: notification.group_id,
          user_id: notification.sender?.id, 
          action: action, 
        }),
      });
      if (res.ok) {
        setStatus(action);
        if (onHandled) onHandled(action);
      }
    } catch (_) {}
  };

  const sender = notification.sender || {};
  const avatar = toMediaUrl(sender.avatar) || "/default-avatar.png";

  return (
    <div className={style.item}>
       <img
          className={style.avatar}
          src={avatar}
          alt={sender.username}
        />
      <div className={style.info}>
        <div className={style.name}>
          <p className={style.message}>{notification.message}</p>
        </div>
      </div>
      {status === "pending" && (
        <div className={style.actions}>
          <button
            className={`${style.btn} ${style.accept}`}
            onClick={() => handleAction("accept")}
          >
            Accept
          </button>
          <button
            className={`${style.btn} ${style.reject}`}
            onClick={() => handleAction("reject")}
          >
            Reject
          </button>
        </div>
      )}
    </div>
  );
}
