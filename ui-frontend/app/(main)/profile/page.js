"use client";
import style from "./page.module.css";
import { Suspense } from "react";
import PostsFeed from "@/app/components/feeds/posts/feed";
import ProfileCard from "@/app/components/profile/profileCard";
import { useAuth } from "@/context/AuthContext";
import { useSearchParams } from "next/navigation"; 
import NotFound from "@/app/not-found"

function ProfileContent() {
    const { user } = useAuth();
    const params = useSearchParams();
    const idParam = params.get("id");
    const userId = idParam ? parseInt(idParam, 10) : user?.id; 

    const isNum = /^\d+$/.test(idParam)  
    if (!userId || !isNum) {
        return <NotFound/>
    }
    
    return (
        <div className={style.middle}>
          <ProfileCard />
          <PostsFeed userId={userId} />
        </div>
    );
}

export default function Profile() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
        <ProfileContent />
    </Suspense>
  );
}
