"use client";
import { useState } from "react";
import style from "./followRequest.module.css";
import { apiFetch, toMediaUrl } from "../../../lib/api";

export default function FollowRequest({ notification, onHandled }) {
  const request = notification;
  const [status, setStatus] = useState(
    request.receiver_is_private === true || request.message === "request"
      ? "pending"
      : request.message === "started following you" 
        ? "accepted" 
        : request.message
  );

  const handleFollowAction = async (senderId, action) => {
    try {
      const res = await apiFetch("/request_follow", {
        method: "POST",
        body: JSON.stringify({
          sender_id: senderId,
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

  const sender = request.sender || {};
  const avatar = toMediaUrl(sender.avatar) || "/default-avatar.png";

  return (
    <div className={style.item}>
      <div className={style.info}>
        <img
          className={style.avatar}
          src={avatar}
          alt={`${sender.first_name || ""} ${sender.last_name || ""}`}
        />
        <div className={style.name}>
          <p>
            {(sender.first_name || "") + " " + (sender.last_name || "")}
          </p>
          <p className={style.username}>@{sender.username}</p>
          {status === "accepted" && (
            <p className={style.message}>Started following you</p>
          )}
          {status === "rejected" && (
            <p className={style.message}>You rejected this request</p>
          )}
        </div>
      </div>
      {status === "pending" && (
        <div className={style.actions}>
          <button
            className={`${style.btn} ${style.accept}`}
            onClick={() => handleFollowAction(sender.id, "accept")}
          >
            Accept
          </button>
          <button
            className={`${style.btn} ${style.reject}`}
            onClick={() => handleFollowAction(sender.id, "reject")}
          >
            Reject
          </button>
        </div>
      )}
    </div>
  );
}
