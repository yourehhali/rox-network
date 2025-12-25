import Link from "next/link";
import styles from "./GroupCard.module.css";

export default function GroupCard({ group }) {
  return (
    <Link href={`/groups/${group.id}`} className={styles.cardLink}>
      <div className={styles.card}>
        <div className={styles.cardTitle}>
          {group.title || group.name || `Group ${group.id}`}
        </div>
        <div className={styles.cardDesc}>{group.description || ""}</div>
      </div>
    </Link>
  );
}
