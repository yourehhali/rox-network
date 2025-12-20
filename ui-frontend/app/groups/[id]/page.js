"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { apiFetch } from "../../../lib/api";
import styles from "./page.module.css";
import LeftBar from "../../components/sidebar/left/leftBar.js";
import RightBar from "../../components/sidebar/right/rightBar.js";

export default function GroupPage() {
  const params = useParams();
  const id = parseInt(params.id, 10);
  const [group, setGroup] = useState(null);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [joinState, setJoinState] = useState(null);

  useEffect(() => {
    const load = async () => {
      try {
        const gRes = await apiFetch(`/Get_Group_By_ID?id=${id}`);
        if (gRes.ok) {
          const g = await gRes.json();
          setGroup(g);
        }
        const pRes = await apiFetch(`/GetGroupPosts?id=${id}`);
        if (pRes.ok) {
          const p = await pRes.json();
          setPosts(Array.isArray(p) ? p : []);
        }
      } catch (_) {}
      setLoading(false);
    };
    if (id) load();
  }, [id]);

  const join = async () => {
    try {
      const res = await apiFetch("/join", {
        method: "POST",
        body: JSON.stringify({ group_id: id }),
      });
      if (res.ok) {
        const data = await res.json();
        setJoinState(data.status || "requested");
      }
    } catch (_) {}
  };

  return (
    <div className={styles.container}>
      <LeftBar />
      <main className={styles.main}>
        {loading ? (
          <h2>Loading...</h2>
        ) : (
          <>
            <div className={styles.header}>
              <h1 className={styles.title}>{group?.title || group?.name || `Group ${id}`}</h1>
              <p className={styles.desc}>{group?.description || 'No description available.'}</p>
              <div className={styles.actions}>
                <button onClick={join} className={styles.joinBtn}>
                  {joinState ? `Status: ${joinState}` : "Join Group"}
                </button>
              </div>
            </div>

            <h3 className={styles.sectionTitle}>Posts</h3>
            <div className={styles.posts}>
              {posts.map((p) => (
                <div key={p.id} className={styles.post}>
                  <div className={styles.postTitle}>{p.title || `Post ${p.id}`}</div>
                  <div className={styles.postContent}>{p.content}</div>
                </div>
              ))}
              {posts.length === 0 && <div className={styles.empty}>No posts yet in this group.</div>}
            </div>
          </>
        )}
      </main>
      <RightBar />
    </div>
  );
}
