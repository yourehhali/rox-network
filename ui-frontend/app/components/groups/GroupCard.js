import Link from "next/link";
import styles from "./GroupCard.module.css";

export default function GroupCard({ group }) {
  const grads = ["bg-1", "bg-2", "bg-3", "bg-4", "bg-5"];
  // Use group.id to select a stable background, or fallback to random if no id (unlikely)
  const bgIndex = group.id ? group.id % grads.length : Math.floor(Math.random() * grads.length);
  const bgClass = grads[bgIndex];

  return (
    <Link href={`/groups/${group.id}`} className={styles.cardLink}>
      <div className={styles.groupCard}>
        <div className={styles.groupCardBg}>
            <div className={styles[bgClass]}></div>
        </div>
        
        <div className={styles.groupCardInfo}>
          <div className={styles.cardTitle}>
            {group.title || group.name || `Group ${group.id}`}
          </div>
          <div className={styles.cardDesc}>
            {group.description || "No description available"}
          </div>
        </div>
      </div>
    </Link>
  );
}
