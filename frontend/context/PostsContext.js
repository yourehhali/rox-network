"use client";
import { createContext, useContext, useState, useCallback } from "react";

const PostsContext = createContext();

export function PostsProvider({ children }) {
  const [refreshTrigger, setRefreshTrigger] = useState(0);
  const [groupRefreshTrigger, setGroupRefreshTrigger] = useState(0);

  const triggerFeedRefresh = useCallback(() => {
    setRefreshTrigger((prev) => prev + 1);
  }, []);

  const triggerGroupRefresh = useCallback(() => {
    setGroupRefreshTrigger((prev) => prev + 1);
  }, []);

  return (
    <PostsContext.Provider value={{ refreshTrigger, triggerFeedRefresh, groupRefreshTrigger, triggerGroupRefresh }}>
      {children}
    </PostsContext.Provider>
  );
}

export function usePosts() {
  const context = useContext(PostsContext);
  if (!context) {
    throw new Error("usePosts must be used within a PostsProvider");
  }
  return context;
}
