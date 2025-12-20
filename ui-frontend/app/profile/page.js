import style from "./page.module.css";
import { Suspense } from "react";
import LeftBar from "../components/sidebar/left/leftBar";
import RightBar from "../components/sidebar/right/rightBar";
import PostsFeed from "../components/feeds/posts/feed";
import ProfileCard from "../components/profile/profileCard";
import { useAuth } from "../../context/AuthContext";

export default function Profile() {
  return (
    <main className={style.main}>
      <LeftBar />
      <div className={style.middle}>
        <Suspense>
          <ProfileCard />
        </Suspense>
        <Suspense>
          <PostsFeed />
        </Suspense>
      </div>
      <RightBar />
    </main>
  );
}
