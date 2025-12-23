"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { apiFetch } from "@/lib/api";
import styles from "./page.module.css";
import CreateEventModal from "@/app/components/events/CreateEventModal";

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
      
      const res = await apiFetch("/CreateGroupPost", {
        method: "POST",
        body: formData,
      });

      if (res.ok) {
          setPostTitle("");
          setPostContent("");
          loadGroupData();
      }
    } catch (err) {
       console.error(err);
    }
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
                       <small>By {p.author_id} on {new Date(p.created_at).toLocaleDateString()}</small>
                     </div>
                   ))}
                   {posts.length === 0 && <p>No posts yet.</p>}
                 </div>
              </div>
            )}

            {activeTab === "events" && (
              <div className={styles.eventsList}>
                {events.map(ev => (
                  <div key={ev.id} className={styles.eventCard}>
                    <h3>{ev.title}</h3>
                    <p>{ev.description}</p>
                    <p>Date: {new Date(ev.date).toLocaleDateString()}</p>
                    {/* Event actions like Going/Not Going could be here */}
                    <div className={styles.eventActions}>
                      <button onClick={() => handleEventResponse(ev.id, "going")}>Going</button>
                      <button onClick={() => handleEventResponse(ev.id, "not_going")}>Not Going</button>
                    </div>
                  </div>
                ))}
                {events.length === 0 && <p>No events found.</p>}
              </div>
            )}

            {activeTab === "members" && (
              <div className={styles.membersList}>
                {members.map(m => (
                  <div key={m.id} className={styles.memberCard}>
                     {m.first_name} {m.last_name} ({m.nickname})
                  </div>
                ))}
                {members.length === 0 && <p>No members found.</p>}
              </div>
            )}
          </div>
        </>
      )}

      {showEventModal && (
        <CreateEventModal 
          groupId={id} 
          onClose={() => setShowEventModal(false)} 
          onSuccess={() => {
            setShowEventModal(false);
            loadGroupData();
          }}
        />
      )}
    </div>
  );
}
