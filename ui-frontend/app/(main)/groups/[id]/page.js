"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { apiFetch } from "@/lib/api";
import styles from "./page.module.css";
import PostsFeed from "@/app/components/feeds/posts/feed";
import EventsFeed from "@/app/components/feeds/events/EventsFeed";
import MemberCard from "@/app/components/groups/MemberCard";
import { usePosts } from "@/context/PostsContext";

export default function GroupPage() {
  const params = useParams();
  const id = parseInt(params.id, 10);
  const { refreshTrigger } = usePosts();
  
  const [group, setGroup] = useState(null);
  const [members, setMembers] = useState([]);
  
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("posts");
  const [joinState, setJoinState] = useState(null);

  const loadGroupData = async () => {
    try {
      // 1. Group Details
      const gRes = await apiFetch(`/Get_Group_By_ID?id=${id}`, { cache: "no-store" });
      if (gRes.ok) {
        const data = await gRes.json();
        // Backend returns { group: {...}, members: [...] }
        setGroup(data.group || data); 
        setMembers(data.members || []);
        
        if (data.group?.user_status) {
          setJoinState(data.group.user_status);
        } else if (data.group?.is_member) {
          setJoinState("member");
        } else {
          setJoinState(null);
        }
      }
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  useEffect(() => {
    if (id) loadGroupData();
  }, [id, refreshTrigger]);

  const join = async () => {
    try {
      const res = await apiFetch("/join", {
        method: "POST",
        body: JSON.stringify({ group_id: id }),
      });
      if (res.ok) {
        const data = await res.json();
        setJoinState(data.status || "requested");
        loadGroupData(); // Refresh to update member status
      }
    } catch (_) {}
  };

  return (
    <div className={styles.middle}>
      {loading ? (
        <h2>Loading...</h2>
      ) : (
        <>
          <div className={styles.header}>
            <h1 className={styles.title}>{group?.title || `Group ${id}`}</h1>
            <p className={styles.desc}>{group?.description}</p>
            
            <div className={styles.actions}>
              {joinState === "member" ? (
                 <span className={styles.memberBadge}>✓ Member</span>
              ) : (
                <button onClick={join} className={styles.joinBtn}>
                  {joinState ? `Status: ${joinState}` : "Join Group"}
                </button>
              )}
            </div>
          </div>

          {/* Tabs */}
          <div className={styles.tabs}>
            <button 
              className={`${styles.tab} ${activeTab === "posts" ? styles.activeTab : ""}`}
              onClick={() => setActiveTab("posts")}
            >
              Posts
            </button>
            <button 
              className={`${styles.tab} ${activeTab === "events" ? styles.activeTab : ""}`}
              onClick={() => setActiveTab("events")}
            >
              Events
            </button>
            <button 
              className={`${styles.tab} ${activeTab === "members" ? styles.activeTab : ""}`}
              onClick={() => setActiveTab("members")}
            >
              Members
            </button>
          </div>

          {/* Content */}
          <div className={styles.content}>
            {activeTab === "posts" && (
              <div>
                 <div className={styles.postsList}>
                    <PostsFeed 
                        groupId={id}
                        endpoints={{
                            comments: null,
                            submitComment: null,
                            submitReaction: null
                        }}
                    />
                 </div>
              </div>
            )}

            {activeTab === "events" && (
              <div className={styles.eventsList}>
                <EventsFeed groupId={id} />
              </div>
            )}

            {activeTab === "members" && (
              <div className={styles.membersList}>
                {members.map(m => (
                  <MemberCard key={m.id} member={m} />
                ))}
                {members.length === 0 && <p>No members found.</p>}
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}
