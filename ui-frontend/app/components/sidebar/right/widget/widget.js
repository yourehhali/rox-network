"use client";

import style from "./widget.module.css";
import { useRouter } from "next/navigation";
import { Bell } from "../../../icons";
import { useEffect, useRef, useState } from "react";
import { apiFetch, toMediaUrl } from "../../../../../lib/api";
import FollowRequest from "../../../notifications/followRequest";
import GroupInvite from "../../../notifications/GroupInvite";
import EventNotification from "../../../notifications/EventNotification";
import JoinRequest from "../../../notifications/JoinRequest";
import { useAuth } from "../../../../../context/AuthContext";

export default function Widget() {
  const router = useRouter();
  const { user } = useAuth();
  const [open, setOpen] = useState(false);
  const [notifications, setNotifications] = useState([]);
  const wsRef = useRef(null);

  const navigateTo = (endpoint) => {
    router.push(endpoint);
  };

  useEffect(() => {
    if (!open) return;
    const load = async () => {
      try {
        const res = await apiFetch("/notifications");
        if (res.ok) {
          const data = await res.json();
          if (Array.isArray(data)) {
            // Deduplicate by ID
            const unique = [];
            const ids = new Set();
            data.forEach(n => {
              if (!ids.has(n.id)) {
                ids.add(n.id);
                unique.push(n);
              }
            });
            setNotifications(unique);
          } else {
            setNotifications([]);
          }
        }
      } catch (_) {}
    };
    load();
    const ws = new WebSocket((toMediaUrl("/ws") || "").replace(/^http/, "ws"));
    wsRef.current = ws;
    
    ws.onmessage = (ev) => {
      try {
        const notif = JSON.parse(ev.data);
        const type = (notif.type || "").trim();
        
        // Handle removal instruction from backend
        if (type === "notification_removed") {
             setNotifications((prev) => prev.filter(n => {
                 if (notif.message === "group_invite") {
                     return !(n.type === "group_invite" && n.group_id === notif.group_id);
                 }
                 if (notif.message === "request_join_group") {
                     // Filter by sender_id (requester) and group_id
                     const sId = n.sender_id || n.sender?.id;
                     return !(n.type === "request_join_group" && n.group_id === notif.group_id && sId === notif.sender_id);
                 }
                 if (notif.message === "event_created") {
                     return !(n.type === "event_created" && n.event_id === notif.event_id);
                 }
                 if (notif.message === "started following you") {
                     const sId = n.sender_id || n.sender?.id;
                     return !(n.type === "new_follower" && sId === notif.sender_id);
                 }
                 if (notif.message === "Invitation_friendships") {
                     // Filter by sender_id
                     const senderId = n.sender_id || n.sender?.id || n.sender?.ID;
                     return !(n.type === "Invitation_friendships" && senderId === notif.sender_id);
                 }
                 return true;
             }));
             return;
        }

        // Add new notification, preventing duplicates
        setNotifications((prev) => {
            // Check if already exists by ID
            if (prev.some(n => n.id === notif.id)) {
                return prev;
            }
            // Also check if logically same (same type, sender, group/event) to be safe
            // But ID check should suffice if backend sends ID
            return [notif, ...prev];
        });
      } catch (_) {}
    };
    
    ws.onclose = () => { wsRef.current = null; };
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [open]);

  const removeLocally = (notifToRemove) => {
    setNotifications((prev) => prev.filter((n) => n !== notifToRemove));
  };

  return (
    <div className={style.widget}>
      <div className={style.notifications} onClick={() => setOpen((v)=>!v)}>
        <Bell/>
      </div>
      <div
        className={style.userAvatar}
        onClick={() =>
          user?.id ? navigateTo(`/profile?id=${user.id}`) : navigateTo("/profile")
        }
      >
        <img
          src={toMediaUrl(user?.avatar || user?.image_path) || "/default-avatar.png"}
          alt="User Avatar"
          className={style.avatarImage}
        />
        <span className={style.onlineIndicator}></span>
      </div>
      {open && (
        <div className={style.notificationsPanel}>
          {notifications.length === 0 ? (
            <div className={style.notificationItem}>No notifications</div>
          ) : notifications.map((n, idx) => {
            const key = n.id || `notif-${idx}`;
            
            if (n.type === "Invitation_friendships" || n.type === "new_follower") {
              return (
                <FollowRequest
                  key={key}
                  notification={n}
                  onHandled={() => removeLocally(n)}
                />
              );
            }
            if (n.type === "group_invite") {
              return (
                <GroupInvite
                  key={key}
                  notification={n}
                  onHandled={() => removeLocally(n)}
                />
              );
            }
            if (n.type === "event_created") {
              return (
                <EventNotification
                  key={key}
                  notification={n}
                  onHandled={() => removeLocally(n)}
                />
              );
            }
            if (n.type === "request_join_group") {
              return (
                <JoinRequest
                  key={key}
                  notification={n}
                  onHandled={() => removeLocally(n)}
                />
              );
            }

            // Generic fallback
            return (
              <div key={key} className={style.notificationItem}>
                <span className={style.notificationTitle}>{n.type || "Notification"}</span>
                <span className={style.notificationMessage}>{n.message || n.content || ""}</span>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
