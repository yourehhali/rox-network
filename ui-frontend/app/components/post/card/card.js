import style from "./card.module.css";
import { Eye } from "lucide-react";


const grads = ["bg-1", "bg-2", "bg-3", "bg-4", "bg-5"];

export default function PostCard({postData, onOpen}) {

  const rndGrad = grads[Math.floor(Math.random() * grads.length)];
  return (
    <div className={style.postCard}>
      {postData.imageUrl ? (
        <div className={style.postCardImage}>
          <img
            src={postData.imageUrl}
            alt="Post Image"
            className={style.image}
          />
        </div>
      ) : (
        <div className={`${style.postCardImage} ${style[rndGrad]}`}></div>
      )}

      <div className={style.postCardInfo}>
        {!postData.imageUrl ? (
          <div className={style.postContentNoImage}>
            <p className={style.lines5}>{postData.content}</p>
          </div>
        ) : null}

        <div className={style.postCardFooter}>
          <div className={style.postCardProfile}>
            <div className={style.postAvatar}>
              <img
                src={postData.avatar}
                alt="Author Avatar"
                className={style.avatar}
              />
            </div>
            <div className={style.postAuthor}>
              <span className={style.authorName}>{postData.author}</span>
              <div className={style.postAgo}>{postData.ago}</div>
            </div>
          </div>
          <div className={style.OpenPost} onClick={onOpen}>
            <Eye className={style.eyeIcon} />
          </div>
        </div>

        {postData.imageUrl ? (
          <div className={style.postContent}>
            <p className={style.lines2}>{postData.content}</p>
          </div>
        ) : null}
      </div>
    </div>
  );
}
