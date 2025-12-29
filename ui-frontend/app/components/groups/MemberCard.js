"use client"

import styles from "./MemberCard.module.css";
import { useRouter } from "next/navigation";




export default function MemberCard({ member }) {

  const router = useRouter();
  const handleClick = () => {
  if (router) {
    router.push(`/profile?id=${member.id}`);
  }
};

  return (
    <div 
    className={styles.memberCard}
    onClick={()=>handleClick()}
    
    >
      {console.log(member)}
      {member.first_name} {member.last_name}
    </div>
  );
}
