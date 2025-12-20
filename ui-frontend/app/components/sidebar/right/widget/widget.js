"use client";

import style from "./widget.module.css";
import { useRouter } from "next/navigation";
import { Bell } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { apiFetch, toMediaUrl } from "../../../../../lib/api";
import FollowRequest from "../../../notifications/followRequest";
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
          setNotifications(Array.isArray(data) ? data : []);
        }
      } catch (_) {}
    };
    load();
    const ws = new WebSocket((toMediaUrl("/ws") || "").replace(/^http/, "ws"));
    wsRef.current = ws;
    ws.onmessage = (ev) => {
      try {
        const notif = JSON.parse(ev.data);
        setNotifications((prev) => [notif, ...prev]);
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
          ) : notifications.map((n) => {
            if (n.type === "Invitation_friendships") {
              return (
                <FollowRequest
                  key={n.id || Math.random()}
                  notification={n}
                  onHandled={(updated) => {
                     // remove locally if handled
                    setNotifications((prev)=> prev.filter((x)=> (x.id||x._k)!==(n.id||n._k)) )
                  }}
                />
              );
            }
            return (
              <div key={n.id || Math.random()} className={style.notificationItem}>
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
