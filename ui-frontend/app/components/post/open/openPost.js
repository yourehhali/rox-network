export default function OpenPost({ postId }) {
  let fetchPost = (postId) => {
    return {
      id: postId,
      login: "yorehhali",
      title: "Younes Test Post",
      content:
        "testing this new open post format, it should be cool and look good on modal open",
      image: "/post1.png",
      avatar: "/avatar1.png",
      likesCount: 5,
      dislikeCount: 2,
    };
  };
  let fetchComments = (postId) => {
    return [
      {
        id: 1,
        parentId: postId,
        avatar: "/avatar2.png",
        login: "yorehhali",
        body: "testing the comment",
      },
      {
        id: 2,
        parentId: postId,
        avatar: "/avatar2.png",
        login: "yorehhali",
        body: "testing the comment 2",
      },
      {
        id: 3,
        parentId: postId,
        avatar: "/avatar2.png",
        login: "yorehhali",
        body: "testing the comment 3",
      },
    ];
  };
  fetchPost(postId);
  fetchComments(postId);
  return (
    <div className={style.openPost}>
      <div className={style.postContainer}>{buildPosts}</div>
    </div>
  );
}
