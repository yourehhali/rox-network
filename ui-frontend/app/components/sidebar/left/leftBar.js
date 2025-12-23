"use client";

import styles from "./leftBar.module.css";
import { useRouter } from "next/navigation";
import { useAuth } from "../../../../context/AuthContext";
import { useState } from "react";
import CreatePostModal from "../../post/create/CreatePostModal";

import {
  Plus,
  Home,
  UserRound,
  UsersRound,
  MessageCircle,
  LogOut,
} from "lucide-react";

export default function LeftBar() {
  const router = useRouter();
  const { user } = useAuth();
  const [showCreateModal, setShowCreateModal] = useState(false);

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
        alert("Logout failed: " + err);
      }
    } catch (err) {
      console.error("Logout error:", err);
    }
    
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
          <a className={styles.menuItem}
            onClick={() => handleLogout()}
          >
            
            <LogOut />
          </a>

          <div className={styles.createButton} onClick={() => setShowCreateModal(true)}>
            <span className={styles.plusIcon}>
              <Plus />
            </span>
          </div>
        </div>
      </div>
      {showCreateModal && (
        <CreatePostModal
          onClose={() => setShowCreateModal(false)}
          onSuccess={() => {
            setShowCreateModal(false);
            if (window.location.pathname === "/") {
              window.location.reload();
            } else {
              navigateTo("/");
            }
          }}
        />
      )}
    </aside>
  );
}
