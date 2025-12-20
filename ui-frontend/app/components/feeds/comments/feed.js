import style from "./feed.module.css";
import CommentCard from "../../post/comment/card/card";

export default function PostsFeed() {
  let comments = [
    {
      id: 1,
    },
    {
      id: 2,
    },
    {
      id: 3,
    },
    {
      id: 4,
    },
  ];
  let buildComments = comments.map((comment) => (
    <CommentCard key={comment.id} />
  ));
  return (
    <div className={style.commentsFeed}>
      <div className={style.commentsFeedContainer}>{buildComments}</div>
    </div>
  );
}
