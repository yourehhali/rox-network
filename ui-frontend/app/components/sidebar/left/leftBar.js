"use client";

import styles from "./leftBar.module.css";
import { useRouter, usePathname } from "next/navigation";
import { useAuth } from "../../../../context/AuthContext";
import { useToast } from "../../../../context/ToastContext";
import { useState, useEffect } from "react";
import { apiFetch } from "../../../../lib/api";
import CreatePostModal from "../../post/create/CreatePostModal";
import CreateGroupModal from "../../groups/CreateGroupModal";
import CreateEventModal from "../../events/CreateEventModal";

import {
  Plus,
  Home,
  UserRound,
  UsersRound,
  MessageCircle,
  Settings,
  LogOut,
} from "lucide-react";

export default function LeftBar() {
  const router = useRouter();
  const pathname = usePathname();
  const { user } = useAuth();
  const { addToast } = useToast();
  
  const [showCreatePostModal, setShowCreatePostModal] = useState(false);
  const [showCreateGroupModal, setShowCreateGroupModal] = useState(false);
  const [showCreateEventModal, setShowCreateEventModal] = useState(false);
  const [showGroupChoice, setShowGroupChoice] = useState(false);
  const [targetGroupId, setTargetGroupId] = useState(null);
  const [isGroupMember, setIsGroupMember] = useState(false);

  useEffect(() => {
    // Reset state when path changes
    setShowGroupChoice(false);
    setTargetGroupId(null);
    setIsGroupMember(false);

    if (pathname && pathname.startsWith("/groups/")) {
        const parts = pathname.split("/");
        // Path is like /groups/123. parts[0]='', parts[1]='groups', parts[2]='123'
        const gid = parts[2]; 
        if (gid && !isNaN(gid)) {
            setTargetGroupId(gid);
            fetchGroupMembership(gid);
        }
    }
  }, [pathname, user]);

  const fetchGroupMembership = async (gid) => {
    if (!user) return;
    try {
       const res = await apiFetch(`/Get_Group_By_ID?id=${gid}`);
       if (res.ok) {
         const data = await res.json();
         if (data.group?.is_member) {
            setIsGroupMember(true);
         }
       }
    } catch (e) {
      console.error("Failed to check membership", e);
    }
  };

  const navigateTo = (endpoint) => {
    router.push(endpoint);
  };
  const handleLogout = async () => {
    try {
      const res = await fetch("http://localhost:8080/logout", {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        localStorage.clear();
        navigateTo("/");
      } else {
        const err = await res.text();
        addToast("Logout failed: " + err, "error");
      }
    } catch (err) {
      console.error("Logout error:", err);
      addToast("Logout error occurred", "error");
    }
    
  };

  const handleCreateButtonClick = () => {
    if (pathname === "/groups") {
        setShowGroupChoice(!showGroupChoice);
    } else if (targetGroupId) {
        // We are on a group page
        setShowGroupChoice(!showGroupChoice);
    } else {
        // Normal post
        setShowCreatePostModal(true);
    }
  };

  const handleGroupPostClick = (e) => {
    e.stopPropagation();
    if (isGroupMember) {
        setShowCreatePostModal(true);
    } else {
        addToast("You must be a member of this group to post.", "error");
    }
    setShowGroupChoice(false);
  };

  const handleGroupEventClick = (e) => {
    e.stopPropagation();
    if (isGroupMember) {
        setShowCreateEventModal(true);
    } else {
        addToast("You must be a member of this group to create events.", "error");
    }
    setShowGroupChoice(false);
  };

  return (
    <aside className={styles.leftSidebar}>
      <div className={styles.container}>
        <div className={styles.topSection}>
          <div className={styles.logo}>
            <img src="/logo.png" alt="Logo" className={styles.logoImage} />
          </div>

          <div className={styles.menu}>
            <a className={styles.menuItem} onClick={() => navigateTo("/")}>
              <Home />
            </a>

            <a
              className={styles.menuItem}
              onClick={() =>
                user?.id
                  ? navigateTo(`/profile?id=${user.id}`)
                  : navigateTo("/profile")
              }
            >
              <UserRound />
            </a>

            <a className={styles.menuItem} onClick={() => navigateTo("/chat")}>
              <MessageCircle />
            </a>

            <a
              className={styles.menuItem}
              onClick={() => navigateTo("/groups")}
            >
              <UsersRound />
            </a>
          </div>
        </div>

        <div className={styles.bottomSection}>
          <a className={styles.menuItem} onClick={() => navigateTo("/settings")}>
            <Settings />
          </a>
          <a className={styles.menuItem}
            onClick={() => handleLogout()}
          >
            
            <LogOut />
          </a>

          <div className={styles.createButton} onClick={handleCreateButtonClick}>
            <span className={styles.plusIcon}>
              <Plus />
            </span>
          </div>
          
          {/* Choice Pill for Group Page */}
          {showGroupChoice && (
            <div className={styles.choicePill}>
              {pathname === "/groups" ? (
                <button onClick={(e) => {
                    e.stopPropagation();
                    setShowCreateGroupModal(true);
                    setShowGroupChoice(false);
                }}>
                    Create Group
                </button>
              ) : (
                <>
                    <button onClick={handleGroupPostClick}>
                        Post
                    </button>
                    <div className={styles.divider}></div>
                    <button onClick={handleGroupEventClick}>
                        Event
                    </button>
                </>
              )}
            </div>
          )}
        </div>
      </div>

      {showCreatePostModal && (
        <CreatePostModal
          groupId={targetGroupId}
          onClose={() => {
            setShowCreatePostModal(false);
          }}
          onSuccess={() => {
            setShowCreatePostModal(false);
            window.location.reload();
          }}
        />
      )}
      
      {showCreateGroupModal && (
          <CreateGroupModal 
            isOpen={showCreateGroupModal}
            onClose={() => setShowCreateGroupModal(false)}
            onSuccess={() => {
                setShowCreateGroupModal(false);
                window.location.reload();
            }}
          />
      )}

      {showCreateEventModal && targetGroupId && (
          <CreateEventModal
            groupId={targetGroupId}
            isOpen={showCreateEventModal}
            onClose={() => {
                setShowCreateEventModal(false);
            }}
            onSuccess={() => {
                setShowCreateEventModal(false);
                window.location.reload();
            }}
          />
      )}
    </aside>
  );
}
