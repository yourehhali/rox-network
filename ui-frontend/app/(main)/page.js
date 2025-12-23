import style from "./page.module.css";
import { Suspense } from "react";
import PostsFeed from "@/app/components/feeds/posts/feed";

export default function Home() {
  return (
    <Suspense>
      <PostsFeed className={style.inner}/>
    </Suspense>
  );
}
