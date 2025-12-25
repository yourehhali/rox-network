"use client";
import { useEffect, useState } from "react";
import style from "./InviteGroupModal.module.css";
import { apiFetch } from "../../../lib/api";
import { X } from "../icons";

export default function InviteGroupModal({ userId, onClose }) {
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [invited, setInvited] = useState(new Set());

  useEffect(() => {
    const fetchGroups = async () => {
      try {
        const res = await apiFetch("/groups/not-member", {
            method: "POST",
            body: JSON.stringify({ user_id: userId })
        });
        if (res.ok) {
          const data = await res.json();
          setGroups(data || []);
          
          const initialInvited = new Set();
          if (data) {
            data.forEach(g => {
                if (g.invite_status === 'invited' || g.invite_status === 'requested') {
                    initialInvited.add(g.id);
                }
            });
          }
          setInvited(initialInvited);
        }
      } catch (e) {
        console.error("Failed to fetch groups", e);
      } finally {
        setLoading(false);
      }
    };
    if (userId) fetchGroups();
  }, [userId]);

  const handleInvite = async (groupId) => {
    try {
        const res = await apiFetch("/groups/invite", {
            method: "POST",
            body: JSON.stringify({ group_id: groupId, user_id: userId })
        });
        if (res.ok) {
            setInvited(prev => new Set(prev).add(groupId));
        }
    } catch (e) {
        console.error("Failed to invite user to group", e);
    }
  };

  return (
    <div className={style.modalOverlay} onClick={onClose}>
      <div className={style.modalContent} onClick={(e) => e.stopPropagation()}>
        <div className={style.header}>
          <h2>Invite to Group</h2>
          <button className={style.closeButton} onClick={onClose}>
            <X size={24} />
          </button>
        </div>
        <div className={style.list}>
          {loading ? (
            <div className={style.loading}>Loading...</div>
          ) : groups.length === 0 ? (
            <div className={style.empty}>No groups available to invite.</div>
          ) : (
            groups.map((g) => (
              <div key={g.id} className={style.groupItem}>
                <div className={style.groupInfo}>
                  <div className={style.groupTitle}>{g.title}</div>
                  <div className={style.groupDesc}>{g.description}</div>
                </div>
                <button 
                    className={style.inviteBtn}
                    onClick={() => handleInvite(g.id)}
                    disabled={invited.has(g.id)}
                >
                    {invited.has(g.id) ? "Invited" : "Invite"}
                </button>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
