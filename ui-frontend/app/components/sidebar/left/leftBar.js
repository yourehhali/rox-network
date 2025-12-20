"use client";

import styles from "./leftBar.module.css";
import { useRouter } from "next/navigation";
import { useAuth } from "../../../../context/AuthContext";

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

  const navigateTo = (endpoint) => {
    router.push(endpoint);
  };

  return (
    <aside className={styles.leftSidebar}>
      <div className={styles.container}>
        <div className={styles.topSection}>
          <div className={styles.logo}>
            <img src="/logo.png" alt="Logo" className={styles.logoImage} />
          </div>

          <div className={styles.menu}>
            <a className={styles.menuItem} onClick={() => navigateTo('/')}>
              <Home />
            </a>

            <a
              className={styles.menuItem}
              onClick={() =>
                user?.id ? navigateTo(`/profile?id=${user.id}`) : navigateTo("/profile")
              }
            >
              <UserRound />
            </a>

            <a className={styles.menuItem} onClick={() => navigateTo('/chat')}>
              <MessageCircle />
            </a>

            <a className={styles.menuItem} onClick={() => navigateTo('/groups')}>
              <UsersRound />
            </a>
          </div>
        </div>

        <div className={styles.bottomSection}>
          <a className={styles.menuItem}>
            <LogOut />
          </a>

          <div className={styles.createButton}>
            <span className={styles.plusIcon}>
              <Plus />
            </span>
          </div>
        </div>
      </div>
    </aside>
  );
}
