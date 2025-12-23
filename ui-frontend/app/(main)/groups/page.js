"use client";
import styles from "./page.module.css";
import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import Link from "next/link";
import CreateGroupModal from "@/app/components/groups/CreateGroupModal";

export default function GroupsPage() {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);

  const loadGroups = async () => {
    try {
      const res = await apiFetch("/Get_Groups");
      if (res.ok) {
        const data = await res.json();
        setGroups(Array.isArray(data) ? data : []);
      }
    } catch (_) {}
    setLoading(false);
  };

  useEffect(() => {
    loadGroups();
  }, []);

  return (
    <div className={styles.middle}>
      <div className={styles.header}>
        <h1>Groups</h1>
        <CreateGroupModal onSuccess={loadGroups} />
      </div>
      {loading ? (
        <div>Loading...</div>
      ) : (
        <div className={styles.grid}>
          {groups.map((g) => (
            <Link key={g.id} href={`/groups/${g.id}`} className={styles.cardLink}>
              <div className={styles.card}>
                <div className={styles.cardTitle}>
                  {g.title || g.name || `Group ${g.id}`}
                </div>
                <div className={styles.cardDesc}>{g.description || ""}</div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
