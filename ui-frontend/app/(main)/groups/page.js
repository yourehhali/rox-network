"use client";
import styles from "./page.module.css";
import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import CreateGroupModal from "@/app/components/groups/CreateGroupModal";
import GroupCard from "@/app/components/groups/GroupCard";
import { usePosts } from "@/context/PostsContext"; 
import {NotFound} from "@/app/not-found" 

export default function GroupsPage() {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const { groupRefreshTrigger } = usePosts();

  const loadGroups = async () => {
    try {
      const res = await apiFetch("/Get_Groups", { cache: "no-store" });
      if (res.ok) {
        const data = await res.json();
        setGroups(Array.isArray(data) ? data : []);
      }else {
        setLoading(false)
        throw new Error("group not found")
        
      }
    } catch (_) {}
    setLoading(false);
  };

  useEffect(() => {
    loadGroups();
  }, [groupRefreshTrigger]);

  return (
    <div className={styles.middle}>
      <div className={styles.header}>
        <h1>Groups</h1>
      </div>
      {loading ? (
        <div>Loading...</div>
      ) : (
        <div className={styles.grid}>
          {groups.map((g) => (
            <GroupCard key={g.id} group={g} />
          ))}
        </div>
      )}
    </div>
  );
}
