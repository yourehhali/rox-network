"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useAuth } from "./AuthContext";
import { toMediaUrl, apiFetch } from "../lib/api";

const ChatContext = createContext(undefined);

export function ChatProvider({ children }) {
  const { user } = useAuth();
  const wsRef = useRef(null);
  const [unreadCount, setUnreadCount] = useState(0);

  const subscribersRef = useRef(new Set());

  const subscribe = (callback) => {
    subscribersRef.current.add(callback);
    return () => subscribersRef.current.delete(callback);
  };

  const fetchUnreadCount = async () => {
    if (!user) return;
    try {
      const res = await apiFetch("/chat/unread/count");
      if (res.ok) {
        const data = await res.json();
        setUnreadCount(data.count || 0);
      }
    } catch (e) {
      console.error("Failed to fetch unread count", e);
    }
  };

  useEffect(() => {
    if (!user) {
      setUnreadCount(0);
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
      return;
    }

    fetchUnreadCount();

    const connect = () => {
      if (wsRef.current?.readyState === WebSocket.OPEN) return;

      const ws = new WebSocket((toMediaUrl("/chat/ws") || "").replace(/^http/, "ws"));
      wsRef.current = ws;

      ws.onopen = () => {
        
      };

      ws.onmessage = (ev) => {
        try {
          const msg = JSON.parse(ev.data);
          if (msg.type === "private_message") {
            fetchUnreadCount();
          }
          subscribersRef.current.forEach(cb => cb(msg));
        } catch (e) {
          
        }
      };

      ws.onclose = () => {
        wsRef.current = null;
        
      };

      ws.onerror = (err) => {
        
      };
    };

    connect();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [user]);

  const sendMessage = (msg) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    }
  };

  return (
    <ChatContext.Provider value={{ sendMessage, subscribe, isConnected: !!wsRef.current, unreadCount, refreshUnreadCount: fetchUnreadCount }}>
      {children}
    </ChatContext.Provider>
  );
}

export function useChat() {
  const context = useContext(ChatContext);
  if (context === undefined) {
    throw new Error("useChat must be used within a ChatProvider");
  }
  return context;
}
