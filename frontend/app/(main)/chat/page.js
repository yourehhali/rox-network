"use client";
import { useEffect, useRef, useState, useLayoutEffect } from "react";
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
  const messagesEndRef = useRef(null);
  const chatContainerRef = useRef(null);
  const [hasMore, setHasMore] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const prevScrollHeightRef = useRef(0);

  useEffect(() => {
    activeChatRef.current = activeChat;
  }, [activeChat]);

  useEffect(() => {
    const loadData = async () => {
      
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

      
      try {
        const res = await apiFetch("/Get_Groups");
        if (res.ok) {
          const data = await res.json();
          
          const memberGroups = (data || []).filter(g => g.user_status === 'member' || g.user_status === 'creator').map(g => ({
            id: g.id,
            name: g.title,
            avatar: "/default-group.png", 
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

        let shouldScroll = false;

        if (msg.type === "private_message" && currentActive.type === "private") {
           
           if (msg.sender_id === currentActive.id || msg.receiver_id === currentActive.id) {
              setMessages((prev) => [...prev, msg]);
              shouldScroll = true;
           }
        } else if (msg.type === "group_message" && currentActive.type === "group") {
           if (msg.group_id === currentActive.id) {
              setMessages((prev) => [...prev, msg]);
              shouldScroll = true;
           }
        }
        
        if (shouldScroll) {
            setTimeout(scrollToBottom, 100);
        }
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const openChat = async (target) => {
    setActiveChat(target);
    setMessages([]);
    setHasMore(true);
    setLoadingMore(false);

    try {
      let url = "";
      if (target.type === "private") {
        url = `/chat/private/history?user_id=${target.id}&limit=10&offset=0`;
      } else if (target.type === "group") {
        url = `/chat/group/history?group_id=${target.id}&limit=10&offset=0`;
      }
      
      const res = await apiFetch(url);
      if (res.ok) {
        const data = await res.json();
        const initialMessages = Array.isArray(data) ? data : [];
        setMessages(initialMessages);
        if (initialMessages.length < 10) setHasMore(false);
        setTimeout(scrollToBottom, 100);
      }
    } catch (_) {}
  };

  const loadMoreMessages = async () => {
    if (loadingMore || !hasMore || !activeChat) return;

    setLoadingMore(true);
    if (chatContainerRef.current) {
        prevScrollHeightRef.current = chatContainerRef.current.scrollHeight;
    }

    try {
        const offset = messages.length;
        let url = "";
        if (activeChat.type === "private") {
            url = `/chat/private/history?user_id=${activeChat.id}&limit=10&offset=${offset}`;
        } else if (activeChat.type === "group") {
            url = `/chat/group/history?group_id=${activeChat.id}&limit=10&offset=${offset}`;
        }

        const res = await apiFetch(url);
        if (res.ok) {
            const data = await res.json();
            const newMessages = Array.isArray(data) ? data : [];
            
            if (newMessages.length < 10) setHasMore(false);
            
            if (newMessages.length > 0) {
                setMessages((prev) => [...newMessages, ...prev]);
            } else {
                setLoadingMore(false);
            }
        } else {
            setLoadingMore(false);
        }
    } catch (_) {
        setLoadingMore(false);
    }
  };

  useLayoutEffect(() => {
    if (loadingMore && chatContainerRef.current) {
        const newScrollHeight = chatContainerRef.current.scrollHeight;
        const diff = newScrollHeight - prevScrollHeightRef.current;
        if (diff > 0) {
            chatContainerRef.current.scrollTop = diff;
        }
        setLoadingMore(false);
    }
  }, [messages, loadingMore]);

  const onScroll = (e) => {
    if (e.target.scrollTop === 0 && hasMore) {
        loadMoreMessages();
    }
  };

  
  const isOnlyEmojis = (str) => {
    if (!str) return false;
    
    const noSpace = str.replace(/\s/g, '');
    if (noSpace.length === 0) return false;
    
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
    
    
    setTimeout(() => {
        inputEl.focus();
        const newCursorPos = start + emojiData.emoji.length;
        inputEl.setSelectionRange(newCursorPos, newCursorPos);
    }, 0);
  };

  return (
    <div className={style.middle}>
      <aside className={`${style.sidebar} ${activeChat ? style.hiddenMobile : ''}`}>
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
      <div className={`${style.chatArea} ${!activeChat ? style.hiddenMobile : ''}`}>
        {activeChat ? (
          <>
            <header className={style.header}>
              <button className={style.backButton} onClick={() => setActiveChat(null)}>
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="m15 18-6-6 6-6"/>
                </svg>
              </button>
              {activeChat.type === 'private' ? (
                 <img src={activeChat.avatar} alt={activeChat.name} />
              ) : (
                 <div style={{ width: '48px', height: '48px', borderRadius: '50%', background: '#444', display: 'flex', alignItems: 'center', justifyContent: 'center', border: '2px solid #F3FFAB', fontSize: '1.5rem' }}>👥</div>
              )}
              <span>{activeChat.name}</span>
            </header>
            <div 
                className={style.messages} 
                ref={chatContainerRef} 
                onScroll={onScroll}
            >
              {loadingMore && <div style={{ textAlign: 'center', padding: '10px', color: '#888' }}>Loading...</div>}
              {messages.map((m, i) => (
                <div key={m.id || i} className={`${style.message} ${m.sender_id === user?.id ? style.own : ''} ${m.is_emoji ? style.bigEmoji : ''}`}>
                  <span className={style.sender}>{m.sender?.username || (m.sender_id === user?.id ? 'You' : m.sender_id)}</span>
                  <p>{m.content}</p>
                </div>
              ))}
              <div ref={messagesEndRef} />
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
