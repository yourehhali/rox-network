"use client";
import style from "./profileCard.module.css";
import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { apiFetch, toMediaUrl } from "../../../lib/api";
import { useAuth } from "../../../context/AuthContext";
import { Mail, Cake, User } from "../icons";
import UserListModal from "./userListModal";
import InviteGroupModal from "./InviteGroupModal";
export default function ProfileCard() {
  const params = useSearchParams();
  const { user } = useAuth();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [followState, setFollowState] = useState({
    following: false,
    pending: false,
  });
  const [modalType, setModalType] = useState(null);
  const [isPrivate, setIsPrivate] = useState(false);

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
          ismyfollower: !!data.ismyfollower,
        });
        setIsPrivate(!!u.is_private);
        setFollowState({ following: !!data.isfriend, pending: !!data.p });
      } catch (e) {
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
      setFollowState({
        following: !!data.following,
        pending: data.status === "pending",
      });
      setProfile((prev) =>
        prev
          ? { ...prev, followers: data.followers_Count || prev.followers }
          : prev
      );
    } catch (e) {
      // ignore
    }
  };

  const handlePrivacyToggle = async () => {
    const newVal = !isPrivate;
    try {
      const res = await apiFetch("/update-privacy", {
        method: "POST",
        body: JSON.stringify({ is_private: newVal ? 1 : 0 }),
      });
      if (res.ok) {
        setIsPrivate(newVal);
      }
    } catch (e) {
      console.error("Failed to update privacy", e);
    }
  };

  if (loading) {
    return (
      <div
        className={style.profileCardContainer}
        style={{ textAlign: "center", paddingTop: "50px" }}
      >
        Loading...
      </div>
    );
  }
  if (!profile) {
    return (
      <div
        className={style.profileCardContainer}
        style={{ textAlign: "center", paddingTop: "50px" }}
      >
        No profile found
      </div>
    );
  }

  const followLabel = followState.pending
    ? "Pending"
    : followState.following
    ? "Unfollow"
    : "Follow";

  const isMyProfile = user && user.id === profile.id;

  return (
    <div className={style.profileCardContainer}>
      <div className={style.profileCoverImage}>
        <img src={"/post1.png"} alt="Cover" className={style.coverImage} />
        <div className={style.overlay}>
          <div>
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
                <div className={style.description}>
                  {profile.description || "No description provided."}
                </div>
                <div className={style.actions}>
                  {!isMyProfile && (
                    <>
                      <button
                        className={`${style.followButton} ${
                          followState.following ? style.unfollowButton : ""
                        }`}
                        onClick={handleToggleFollow}
                      >
                        {followLabel}
                      </button>
                      {ismyfollower && (
                        <button
                          className={style.inviteButton}
                          onClick={() => setModalType("invite")}
                        >
                          Invite to Group
                        </button>
                      )}
                    </>
                  )}
                  {isMyProfile && (
                    <button
                      className={`${style.privacyButton} ${
                        !isPrivate ? style.public : ""
                      }`}
                      onClick={handlePrivacyToggle}
                    >
                      {isPrivate ? "Private" : "Public"}
                    </button>
                  )}
                </div>
              </div>
            </div>
            <div className={style.bottom}>
              <div className={style.contact}>
                {profile.email && followState.following &&(
                  <div className={style.contactItem}>
                    <Mail size={16} color="white" /> <div>{profile.email}</div>
                  </div>
                )}

                {profile.age && followState.following && (
                  <div className={style.contactItem}>
                    <Cake size={16} color="white" /> <div>{profile.age} yo</div>
                  </div>
                )}

                {profile.gender && followState.following && (
                  <div className={style.contactItem}>
                    <User size={16} color="white" /> <div>{profile.gender}</div>
                  </div>
                )}
              </div>
              <div className={style.stats}>
                <div
                  className={style.statItem}
                  onClick={() => followState.following && setModalType("followers")}
                >
                  <div className={style.statValue}>{profile.followers}</div>
                  <div className={style.statLabel}>Followers</div>
                </div>
                <div
                  className={style.statItem}
                  onClick={() => followState.following && setModalType("following")}
                >
                  <div className={style.statValue}>{profile.following}</div>
                  <div className={style.statLabel}>Following</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {modalType === "invite" && (
        <InviteGroupModal
          userId={profile.id}
          onClose={() => setModalType(null)}
        />
      )}

      {/* Reusing UserListModal for followers/following if implemented */}
      {modalType && modalType !== "invite" && (
        <UserListModal
          type={modalType}
          userId={profile.id}
          onClose={() => setModalType(null)}
        />
      )}
    </div>
  );
}
