"use client";
import { useEffect, useRef, useState } from "react";
import { apiFetch, toMediaUrl } from "../../lib/api";
import LeftBar from "../components/sidebar/left/leftBar.js";
import style from "./chat.module.css";

export default function ChatPage() {
  const [friends, setFriends] = useState([]);
  const [activeUser, setActiveUser] = useState(null);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const wsRef = useRef(null);

  useEffect(() => {
    const loadFriends = async () => {
      try {
        const res = await apiFetch("/Friends");
        if (res.ok) {
          const data = await res.json();
          const mapped = (data || []).map((f) => ({
            id: f.id,
            name: f.full_name || f.username,
            avatar: toMediaUrl(f.image_path) || "/default-avatar.png",
          }));
          setFriends(mapped);
        }
      } catch (_) {}
    };
    loadFriends();
    const ws = new WebSocket((toMediaUrl("/chat/ws") || "").replace(/^http/, "ws"));
    wsRef.current = ws;
    ws.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data);
        if (msg.type === "private_message" && (msg.sender_id === activeUser?.id || msg.receiver_id === activeUser?.id)) {
          setMessages((prev) => [...prev, msg]);
        }
      } catch (_) {}
    };
    ws.onclose = () => { wsRef.current = null; };
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [activeUser?.id]);

  const openChat = async (user) => {
    setActiveUser(user);
    try {
      const res = await apiFetch(`/chat/private/history?user_id=${user.id}`);
      if (res.ok) {
        const data = await res.json();
        setMessages(Array.isArray(data) ? data : []);
      }
    } catch (_) {}
  };

  const sendMessage = () => {
    if (!wsRef.current || !activeUser || !input.trim()) return;
    wsRef.current.send(JSON.stringify({
      type: "private_message",
      receiver_id: activeUser.id,
      content: input,
      is_emoji: false,
    }));
    setInput("");
  };

  return (
    <div className={style.container}>
      <LeftBar />
      <aside className={style.sidebar}>
        <h3>Chats</h3>
        <div className={style.list}>
          {friends.map((f) => (
            <div key={f.id} className={`${style.user} ${activeUser?.id === f.id ? style.active : ''}`} onClick={() => openChat(f)}>
              <img src={f.avatar} alt={f.name} />
              <span>{f.name}</span>
            </div>
          ))}
        </div>
      </aside>
      <main className={style.chatArea}>
        {activeUser ? (
          <>
            <header className={style.header}>
              <img src={activeUser.avatar} alt={activeUser.name} />
              <span>{activeUser.name}</span>
            </header>
            <div className={style.messages}>
              {messages.map((m, i) => (
                <div key={m.id || i} className={`${style.message} ${m.sender_id === user?.id ? style.own : ''}`}>
                  <span className={style.sender}>{m.sender?.username || (m.sender_id === user?.id ? 'You' : m.sender_id)}</span>
                  <p>{m.content}</p>
                </div>
              ))}
            </div>
            <div className={style.inputRow}>
              <input value={input} onChange={(e)=>setInput(e.target.value)} placeholder="Type a message..." />
              <button onClick={sendMessage}>Send</button>
            </div>
          </>
        ) : (
          <div className={style.placeholder}>Select a chat</div>
        )}
      </main>
    </div>
  );
}
