"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useAuth } from "./AuthContext";
import { toMediaUrl } from "../lib/api";

const ChatContext = createContext(undefined);

export function ChatProvider({ children }) {
  const { user } = useAuth();
  const wsRef = useRef(null);
  
  // We can track the active chat user globally if we want to filter messages here,
  // or just store all messages and let the UI filter.
  // For simplicity and performance, let's keep the socket management here 
  // and expose the socket and a message listener mechanism.
  
  // However, the current ChatPage logic filters messages based on `activeUser`.
  // To keep it compatible, we can broadcast all incoming messages to subscribers.
  
  const subscribersRef = useRef(new Set());

  const subscribe = (callback) => {
    subscribersRef.current.add(callback);
    return () => subscribersRef.current.delete(callback);
  };

  useEffect(() => {
    if (!user) {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
      return;
    }

    // Connect
    const connect = () => {
      if (wsRef.current?.readyState === WebSocket.OPEN) return;

      const ws = new WebSocket((toMediaUrl("/chat/ws") || "").replace(/^http/, "ws"));
      wsRef.current = ws;

      ws.onopen = () => {
        // Connected
      };

      ws.onmessage = (ev) => {
        try {
          const msg = JSON.parse(ev.data);
          // Notify subscribers
          subscribersRef.current.forEach(cb => cb(msg));
        } catch (e) {
          // Parse Error
        }
      };

      ws.onclose = () => {
        wsRef.current = null;
        // Optional: Reconnect logic could go here
      };

      ws.onerror = (err) => {
        // WS Error
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
    <ChatContext.Provider value={{ sendMessage, subscribe, isConnected: !!wsRef.current }}>
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
