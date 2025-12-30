import LeftBar from "@/app/components/sidebar/left/leftBar";
import RightBar from "@/app/components/sidebar/right/rightBar";
import styles from "./layout.module.css";

export default function MainLayout({ children }) {
  return (
    <main className={styles.main}>
      <LeftBar />
      {children}
      <RightBar />
    </main>
  );
}
