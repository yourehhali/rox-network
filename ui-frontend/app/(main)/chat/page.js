"use client";
import { useEffect, useRef, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { useChat } from "@/context/ChatContext";
import { apiFetch, toMediaUrl } from "@/lib/api";
import style from "./chat.module.css";

export default function ChatPage() {
  const { user } = useAuth();
  const { sendMessage: sendWsMessage, subscribe } = useChat();
  const [friends, setFriends] = useState([]);
  const [groups, setGroups] = useState([]);
  const [activeChat, setActiveChat] = useState(null);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");

  const activeChatRef = useRef(activeChat);

  useEffect(() => {
    activeChatRef.current = activeChat;
  }, [activeChat]);

  useEffect(() => {
    const loadData = async () => {
      // Load Friends
      try {
        const res = await apiFetch("/Friends");
        if (res.ok) {
          const data = await res.json();
          const mapped = (data || []).map((f) => ({
            id: f.id,
            name: f.full_name || f.username,
            avatar: toMediaUrl(f.image_path) || "/default-avatar.png",
            type: "private",
          }));
          setFriends(mapped);
        }
      } catch (_) {}

      // Load Groups
      try {
        const res = await apiFetch("/Get_Groups");
        if (res.ok) {
          const data = await res.json();
          // Filter only groups where user is a member
          const memberGroups = (data || []).filter(g => g.user_status === 'member' || g.user_status === 'creator').map(g => ({
            id: g.id,
            name: g.title,
            avatar: "/default-group.png", // Groups don't have avatars yet
            type: "group",
          }));
          setGroups(memberGroups);
        }
      } catch (_) {}
    };
    loadData();
    
    // Subscribe to global chat messages
    const unsubscribe = subscribe((msg) => {
        const currentActive = activeChatRef.current;
        if (!currentActive) return;

        if (msg.type === "private_message" && currentActive.type === "private") {
           // Only add message if it belongs to the current active chat
           if (msg.sender_id === currentActive.id || msg.receiver_id === currentActive.id) {
              setMessages((prev) => [...prev, msg]);
           }
        } else if (msg.type === "group_message" && currentActive.type === "group") {
           if (msg.group_id === currentActive.id) {
              setMessages((prev) => [...prev, msg]);
           }
        }
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const openChat = async (target) => {
    setActiveChat(target);
    try {
      let url = "";
      if (target.type === "private") {
        url = `/chat/private/history?user_id=${target.id}`;
      } else if (target.type === "group") {
        url = `/chat/group/history?group_id=${target.id}`;
      }
      
      const res = await apiFetch(url);
      if (res.ok) {
        const data = await res.json();
        setMessages(Array.isArray(data) ? data : []);
      }
    } catch (_) {}
  };

  const sendMessage = () => {
    if (!activeChat || !input.trim()) return;
    
    if (activeChat.type === "private") {
      sendWsMessage({
        type: "private_message",
        receiver_id: activeChat.id,
        content: input,
        is_emoji: false,
      });
    } else if (activeChat.type === "group") {
      sendWsMessage({
        type: "group_message",
        group_id: activeChat.id,
        content: input,
        is_emoji: false,
      });
    }
    setInput("");
  };

  return (
    <div className={style.middle}>
      <aside className={style.sidebar}>
        <h3>Friends</h3>
        <div className={style.list} style={{ flex: '0 0 auto', maxHeight: '40%' }}>
          {friends.map((f) => (
            <div key={f.id} className={`${style.user} ${activeChat?.id === f.id && activeChat?.type === 'private' ? style.active : ''}`} onClick={() => openChat(f)}>
              <img src={f.avatar} alt={f.name} />
              <span>{f.name}</span>
            </div>
          ))}
        </div>
        
        <h3 style={{ marginTop: '20px', borderTop: '1px solid #333', paddingTop: '10px' }}>Groups</h3>
        <div className={style.list}>
          {groups.map((g) => (
            <div key={g.id} className={`${style.user} ${activeChat?.id === g.id && activeChat?.type === 'group' ? style.active : ''}`} onClick={() => openChat(g)}>
               <div style={{ width: '40px', height: '40px', borderRadius: '50%', background: '#444', display: 'flex', alignItems: 'center', justifyContent: 'center', border: '2px solid #333' }}>
                  <span style={{ fontSize: '1.2rem' }}>👥</span>
               </div>
              <span>{g.name}</span>
            </div>
          ))}
        </div>
      </aside>
      <div className={style.chatArea}>
        {activeChat ? (
          <>
            <header className={style.header}>
              {activeChat.type === 'private' ? (
                 <img src={activeChat.avatar} alt={activeChat.name} />
              ) : (
                 <div style={{ width: '48px', height: '48px', borderRadius: '50%', background: '#444', display: 'flex', alignItems: 'center', justifyContent: 'center', border: '2px solid #F3FFAB', fontSize: '1.5rem' }}>👥</div>
              )}
              <span>{activeChat.name}</span>
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
          <div className={style.placeholder}>Select a chat to start messaging</div>
        )}
      </div>
    </div>
  );
}
