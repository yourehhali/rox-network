import styles from "./MemberCard.module.css";

export default function MemberCard({ member }) {
  return (
    <div className={styles.memberCard}>
      {member.first_name} {member.last_name} ({member.nickname})
    </div>
  );
}
