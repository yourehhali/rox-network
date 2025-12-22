import style from "./page.module.css";
import { Suspense } from "react";
import LeftBar from "./components/sidebar/left/leftBar";
import RightBar from "./components/sidebar/right/rightBar";
import PostsFeed from "./components/feeds/posts/feed";

export default function Home() {
  return (
    <main className={style.main}>
      <LeftBar/>
      <Suspense>
        <PostsFeed className={style.inner}/>
      </Suspense>
      <RightBar/>
    </main>
  );
}
