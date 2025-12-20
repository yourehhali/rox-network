"use client";
import style from "./profileCard.module.css";
import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { apiFetch, toMediaUrl } from "../../../lib/api";
import { useAuth } from "../../../context/AuthContext";

export default function ProfileCard() {
  const params = useSearchParams();
  const { user } = useAuth();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [followState, setFollowState] = useState({ following: false, pending: false });

  useEffect(() => {
    const loadProfile = async () => {
      try {
        const idParam = params.get("id");
        const idToLoad = idParam ? parseInt(idParam, 10) : user?.id;
        if (!idToLoad) {
          setLoading(false);
          return;
        }
        const res = await apiFetch(`/profile?id=${idToLoad}`);
        if (!res.ok) {
          throw new Error("Failed to load profile");
        }
        const data = await res.json();
        const u = data.user || data;
        const avatarUrl = toMediaUrl(u.avatar) || "/default-avatar.png";
        setProfile({
          id: u.id,
          name: `${u.first_name} ${u.last_name}`.trim() || u.username,
          username: u.username,
          avatar: avatarUrl,
          description: u.about || "",
          email: u.email,
          gender: u.gender || "",
          age: u.age || "",
          followers: data.followers_count || 0,
          following: data.following_count || 0,
        });
        setFollowState({ following: !!data.isfriend, pending: !!data.p });
      } catch (e) {
        // ignore
      } finally {
        setLoading(false);
      }
    };
    loadProfile();
  }, [params, user]);

  const handleToggleFollow = async () => {
    if (!profile?.id) return;
    try {
      const res = await apiFetch("/toggle-follow", {
        method: "POST",
        body: JSON.stringify({ followed_id: profile.id }),
      });
      if (!res.ok) return;
      const data = await res.json();
      setFollowState({ following: !!data.following, pending: data.status === "pending" });
      setProfile((prev) => prev ? { ...prev, followers: data.followers_Count || prev.followers } : prev);
    } catch (e) {
      // ignore
    }
  };

  if (loading) {
    return <div className={style.profileCardContainer} style={{textAlign:'center', paddingTop:'50px'}}>Loading...</div>;
  }
  if (!profile) {
    return <div className={style.profileCardContainer} style={{textAlign:'center', paddingTop:'50px'}}>No profile found</div>;
  }

  const followLabel = followState.pending ? "Pending" : (followState.following ? "Unfollow" : "Follow");
  // Don't show follow button if it's my own profile
  const isMyProfile = user && user.id === profile.id;

  return (
    <div className={style.profileCardContainer}>
      <div className={style.profileCoverImage}>
        <img
          src={"/post1.png"}
          alt="Cover"
          className={style.coverImage}
        />
        <div className={style.overlay}></div>
        
        <div className={style.profileInfo}>
          <img
            src={profile.avatar}
            alt={profile.name}
            className={style.avatar}
          />
          <div className={style.userInfo}>
            <div>
              <h2 className={style.name}>{profile.name}</h2>
              <p className={style.username}>@{profile.username}</p>
            </div>
            
            <div className={style.actions}>
              {!isMyProfile && (
                <button 
                  className={`${style.followButton} ${followState.following ? style.unfollowButton : ''}`} 
                  onClick={handleToggleFollow}
                >
                  {followLabel}
                </button>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className={style.stats}>
        <div className={style.statItem}>
          <span className={style.statValue}>{profile.followers}</span>
          <span className={style.statLabel}>Followers</span>
        </div>
        <div className={style.statItem}>
          <span className={style.statValue}>{profile.following}</span>
          <span className={style.statLabel}>Following</span>
        </div>
      </div>

      <div className={style.description}>
        {profile.description || "No description provided."}
        <div style={{marginTop: '20px', display: 'flex', gap: '20px', fontSize: '0.9rem', color: '#888', flexWrap: 'wrap'}}>
          {profile.email && <span>📧 {profile.email}</span>}
          {profile.age && <span>🎂 {profile.age} yo</span>}
          {profile.gender && <span>⚧ {profile.gender}</span>}
        </div>
      </div>
    </div>
  );
}
