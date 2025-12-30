"use client";
import { useState } from "react";
import style from "./followRequest.module.css"; 
import { apiFetch } from "../../../lib/api";

export default function GroupInvite({ notification, onHandled }) {
  const [status, setStatus] = useState("pending");

  const handleAction = async (action) => {
    try {
      const res = await apiFetch("/group/invite/respond", {
        method: "POST",
        body: JSON.stringify({
          group_id: notification.group_id,
          notif_id: notification.id,
          status: action, 
        }),
      });
      if (res.ok) {
        setStatus(action);
        if (onHandled) onHandled(action);
      }
    } catch (_) {}
  };

  const sender = notification.sender || {};
  const groupName = notification.group_title || `Group ${notification.group_id}`;

  return (
    <div className={style.item}>
      <div className={style.info}>
        <div className={style.name}>
          <p className={style.message}>{sender.first_name} invited you to join</p>
          <p style={{ fontWeight: "bold", color: "#F3FFAB" }}>{groupName}</p>
        </div>
      </div>
      {status === "pending" && (
        <div className={style.actions}>
          <button
            className={`${style.btn} ${style.accept}`}
            onClick={() => handleAction("accept")}
          >
            Join
          </button>
          <button
            className={`${style.btn} ${style.reject}`}
            onClick={() => handleAction("reject")}
          >
            Decline
          </button>
        </div>
      )}
    </div>
  );
}
