"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useAuth } from "./AuthContext";
import { toMediaUrl } from "../lib/api";

const ChatContext = createContext(undefined);

export function ChatProvider({ children }) {
  const { user } = useAuth();
  const wsRef = useRef(null);
  
  
  
  
  
  
  
  
  
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

    
    const connect = () => {
      if (wsRef.current?.readyState === WebSocket.OPEN) return;

      const ws = new WebSocket((toMediaUrl("/chat/ws") || "").replace(/^http/, "ws"));
      wsRef.current = ws;

      ws.onopen = () => {
        
      };

      ws.onmessage = (ev) => {
        try {
          const msg = JSON.parse(ev.data);
          
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
