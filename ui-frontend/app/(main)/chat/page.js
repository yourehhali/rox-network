"use client";
import { useEffect, useRef, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { useChat } from "@/context/ChatContext";
import { useToast } from "@/context/ToastContext";
import { apiFetch, toMediaUrl } from "@/lib/api";
import style from "./chat.module.css";
import EmojiPicker from "../../components/emoji/EmojiPicker";
import { Smile } from "../../components/icons";

export default function ChatPage() {
  const { user } = useAuth();
  const { sendMessage: sendWsMessage, subscribe } = useChat();
  const { addToast } = useToast();
  const [friends, setFriends] = useState([]);
  const [groups, setGroups] = useState([]);
  const [activeChat, setActiveChat] = useState(null);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);

  const activeChatRef = useRef(activeChat);
  const inputRef = useRef(null);

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
    
    const unsubscribe = subscribe((msg) => {
        if (msg.type === "error") {
            addToast(msg.message, "error");
            return;
        }

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

  // Helper to check if string contains only emojis
  const isOnlyEmojis = (str) => {
    if (!str) return false;
    // Remove whitespace
    const noSpace = str.replace(/\s/g, '');
    if (noSpace.length === 0) return false;
    // Regex for emojis
    const emojiRegex = /^(?:[\u2700-\u27bf]|(?:\ud83c[\udde6-\uddff]){2}|[\ud800-\udbff][\udc00-\udfff]|[\u0023-\u0039]\ufe0f?\u20e3|\u3299|\u3297|\u303d|\u3030|\u24c2|\ud83c[\udd70-\udd71]|\ud83c[\udd7e-\udd7f]|\ud83c\udd8e|\ud83c[\udd91-\udd9a]|\ud83c[\udde6-\uddff]|\ud83c[\ude01-\ude02]|\ud83c\ude1a|\ud83c\ude2f|\ud83c[\ude32-\ude3a]|\ud83c[\ude50-\ude51]|\u203c|\u2049|[\u25aa-\u25ab]|\u25b6|\u25c0|[\u25fb-\u25fe]|\u00a9|\u00ae|\u2122|\u2139|\ud83c\udc04|[\u2600-\u26FF]|\u2b05|\u2b06|\u2b07|\u2b1b|\u2b1c|\u2b50|\u2b55|\u231a|\u231b|\u2328|\u23cf|[\u23e9-\u23f3]|[\u23f8-\u23fa]|\ud83c\udccf|\u2934|\u2935|[\u2190-\u21ff])+$/;
    return emojiRegex.test(noSpace);
  };

  const sendMessage = () => {
    if (!activeChat || !input.trim()) return;
    
    const isEmoji = isOnlyEmojis(input);

    if (activeChat.type === "private") {
      sendWsMessage({
        type: "private_message",
        receiver_id: activeChat.id,
        content: input,
        is_emoji: isEmoji,
      });
    } else if (activeChat.type === "group") {
      sendWsMessage({
        type: "group_message",
        group_id: activeChat.id,
        content: input,
        is_emoji: isEmoji,
      });
    }
    setInput("");
    setShowEmojiPicker(false);
  };

  const onEmojiClick = (emojiData) => {
    const inputEl = inputRef.current;
    if (!inputEl) {
        setInput((prev) => prev + emojiData.emoji);
        return;
    }

    const start = inputEl.selectionStart;
    const end = inputEl.selectionEnd;
    const text = input;
    
    const before = text.substring(0, start);
    const after = text.substring(end);
    
    const newText = before + emojiData.emoji + after;
    setInput(newText);
    
    // Restore cursor position after the emoji
    setTimeout(() => {
        inputEl.focus();
        const newCursorPos = start + emojiData.emoji.length;
        inputEl.setSelectionRange(newCursorPos, newCursorPos);
    }, 0);
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
                <div key={m.id || i} className={`${style.message} ${m.sender_id === user?.id ? style.own : ''} ${m.is_emoji ? style.bigEmoji : ''}`}>
                  <span className={style.sender}>{m.sender?.username || (m.sender_id === user?.id ? 'You' : m.sender_id)}</span>
                  <p>{m.content}</p>
                </div>
              ))}
            </div>
            <div className={style.inputRow}>
              {showEmojiPicker && (
                <div className={style.emojiPickerContainer}>
                  <EmojiPicker onEmojiClick={onEmojiClick} />
                </div>
              )}
              <button 
                className={style.emojiButton} 
                onClick={() => setShowEmojiPicker(!showEmojiPicker)}
              >
                <Smile size={20} />
              </button>
              <input 
                ref={inputRef}
                value={input} 
                onChange={(e)=>setInput(e.target.value)} 
                placeholder="Type a message..." 
                maxLength={2000}
              />
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
