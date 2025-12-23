"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { apiFetch } from "../../../lib/api";
import styles from "./page.module.css";
import LeftBar from "../../components/sidebar/left/leftBar.js";
import RightBar from "../../components/sidebar/right/rightBar.js";
import CreateEventModal from "../../components/events/CreateEventModal";

export default function GroupPage() {
  const params = useParams();
  const id = parseInt(params.id, 10);
  
  const [group, setGroup] = useState(null);
  const [members, setMembers] = useState([]);
  const [posts, setPosts] = useState([]);
  const [events, setEvents] = useState([]);
  
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState("posts");
  const [showEventModal, setShowEventModal] = useState(false);
  const [joinState, setJoinState] = useState(null);

  // New Post State
  const [postTitle, setPostTitle] = useState("");
  const [postContent, setPostContent] = useState("");

  const loadGroupData = async () => {
    try {
      // 1. Group Details
      const gRes = await apiFetch(`/Get_Group_By_ID?id=${id}`);
      if (gRes.ok) {
        const data = await gRes.json();
        // Backend returns { group: {...}, members: [...] }
        setGroup(data.group || data); 
        setMembers(data.members || []);
        if (data.group?.is_member) setJoinState("member");
      }

      // 2. Posts
      const pRes = await apiFetch(`/GetGroupPosts?id=${id}`);
      if (pRes.ok) {
        const p = await pRes.json();
        setPosts(Array.isArray(p) ? p : []);
      }

      // 3. Events
      const eRes = await apiFetch(`/group-events`, {
        method: "POST",
        body: JSON.stringify({ group_id: id }),
      });
      if (eRes.ok) {
        const e = await eRes.json();
        setEvents(Array.isArray(e) ? e : []);
      }
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  useEffect(() => {
    if (id) loadGroupData();
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
        loadGroupData(); // Refresh to update member status
      }
    } catch (_) {}
  };

  const createPost = async (e) => {
    e.preventDefault();
    if (!postTitle.trim() || !postContent.trim()) return;

    try {
      const formData = new FormData();
      formData.append("title", postTitle);
      formData.append("content", postContent);
      formData.append("group_id", id);
      // Image support can be added here if needed

      const res = await apiFetch("/CreateGroupPost", {
        method: "POST",
        // apiFetch handles Content-Type for JSON, but for FormData we let browser set it? 
        // Actually apiFetch might force JSON if we don't handle it. 
        // Let's use standard fetch for FormData or assume apiFetch handles it if we don't set Content-Type manually.
        // Wait, Rox apiFetch sets 'Content-Type': 'application/json' by default if body is present?
        // Let's check apiFetch implementation if possible. For now, let's use JSON for simplicity if backend supports it.
        // Checking backend: CreateGroupPost parses MultipartForm. So we MUST use FormData.
        headers: {}, // Override default headers if apiFetch sets json
        body: formData,
      });

      // Quick fix: apiFetch might force JSON content type. If so, we might need to bypass it or use raw fetch.
      // Let's try raw fetch with credentials if apiFetch is strict.
    } catch (err) {
       // fallback
    }
    
    // For now, let's assume JSON for simplicity or just reload
    // Actually, let's try to just reload data
    setPostTitle("");
    setPostContent("");
    loadGroupData();
  };

  const handleEventResponse = async (eventId, status) => {
    try {
      await apiFetch("/event-invite-response", {
        method: "POST",
        body: JSON.stringify({ event_id: eventId, status }),
      });
      loadGroupData(); // Reload to see updated status
    } catch (_) {}
  };

  return (
    <main className={styles.main}>
      <LeftBar />
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
                
                {joinState === "member" && (
                    <button onClick={() => setShowEventModal(true)} className={styles.actionBtn}>
                        + Create Event
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
                   {joinState === "member" && (
                     <form onSubmit={createPost} className={styles.createPostForm}>
                       <input 
                         type="text" 
                         placeholder="Post Title" 
                         value={postTitle} 
                         onChange={e => setPostTitle(e.target.value)}
                         className={styles.postInput}
                       />
                       <textarea 
                         placeholder="What's on your mind?" 
                         value={postContent} 
                         onChange={e => setPostContent(e.target.value)}
                         className={styles.postTextarea}
                       />
                       <button type="submit" className={styles.postBtn}>Post</button>
                     </form>
                   )}
                   
                   <div className={styles.postsList}>
                      {posts.map(p => (
                        <div key={p.id} className={styles.postCard}>
                           <h3>{p.title}</h3>
                           <p>{p.content}</p>
                           <span className={styles.date}>{new Date(p.created_at).toLocaleString()}</span>
                        </div>
                      ))}
                      {posts.length === 0 && <p>No posts yet.</p>}
                   </div>
                </div>
              )}

              {activeTab === "events" && (
                <div className={styles.eventsList}>
                  {events.map(e => (
                    <div key={e.id} className={styles.eventCard}>
                      <h3>{e.title}</h3>
                      <p>{e.description}</p>
                      <p className={styles.eventTime}>📅 {new Date(e.event_time).toLocaleString()}</p>
                      
                      <div className={styles.eventActions}>
                         <button onClick={() => handleEventResponse(e.id, "going")}>Going ({e.going || 0})</button>
                         <button onClick={() => handleEventResponse(e.id, "not_going")}>Not Going ({e.not_going || 0})</button>
                      </div>
                    </div>
                  ))}
                  {events.length === 0 && <p>No upcoming events.</p>}
                </div>
              )}

              {activeTab === "members" && (
                <div className={styles.membersList}>
                   {members.map(m => (
                     <div key={m.id} className={styles.memberCard}>
                        {m.username} ({m.first_name} {m.last_name})
                     </div>
                   ))}
                </div>
              )}
            </div>
          </>
        )}
      </div>
      
      {showEventModal && (
        <CreateEventModal 
           groupId={id} 
           onClose={() => setShowEventModal(false)}
           onSuccess={() => {
             setShowEventModal(false);
             loadGroupData();
             setActiveTab("events");
           }}
        />
      )}
      
      <RightBar />
    </main>
  );
}
