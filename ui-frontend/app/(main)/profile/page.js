"use client";
import style from "./page.module.css";
import { Suspense } from "react";
import PostsFeed from "@/app/components/feeds/posts/feed";
import ProfileCard from "@/app/components/profile/profileCard";

export default function Profile() {
  return (
    <div className={style.middle}>
      <Suspense>
        <ProfileCard />
      </Suspense>
      <Suspense>
        <PostsFeed />
      </Suspense>
    </div>
  );
}
