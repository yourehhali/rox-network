import style from "./openPost.module.css";
import { X, Heart, ThumbsDown, Image as ImageIcon } from "../../icons";
import { useState, useEffect, useRef } from "react";
import { apiFetch, toMediaUrl } from "../../../../lib/api";
import { useAuth } from "../../../../context/AuthContext";

export default function OpenPost({ postId: propPostId, onClose, postData, endpoints }) {
  const { user } = useAuth();
  const postId = propPostId || postData?.id;
  
  // Default endpoints
  const commentsUrl = endpoints?.comments || `/posts/${postId}/comments`;
  const submitCommentUrl = endpoints?.submitComment || `/posts/${postId}/comments`;
  const submitReactionUrl = endpoints?.submitReaction || "/reactions";

  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [bgClass, setBgClass] = useState("");
  const [newComment, setNewComment] = useState("");
  const [commentImage, setCommentImage] = useState(null);
  const [likesCount, setLikesCount] = useState(0);
  const [dislikeCount, setDislikeCount] = useState(0);
  const [userReaction, setUserReaction] = useState(null);
  const fileInputRef = useRef(null);

  useEffect(() => {
    if (postData) {
        setPost(postData);
        if (!postData.imageUrl) {
            const gradients = ["bg1", "bg2", "bg3", "bg4", "bg5"];
            setBgClass(gradients[Math.floor(Math.random() * gradients.length)]);
        }
        // Initialize counts from postData if available (assuming postData might have them)
        // But better to fetch fresh data if possible or use what we have
        setLikesCount(postData.likesCount || 0); // Need to ensure postData has these
        setDislikeCount(postData.dislikeCount || 0); // Need to ensure postData has these
        if (postData.userReaction) {
            setUserReaction(postData.userReaction);
        }
    } else {
        // Fallback fetch
        const fetchPost = async () => {
            try {
                const res = await apiFetch(`/posts/${postId}`); // Assuming this endpoint exists
                if (res.ok) {
                    const data = await res.json();
                    setPost(data);
                    // Set counts
                }
            } catch (e) {
                console.error(e);
            }
        };
        // fetchPost(); 
        // For now sticking to the logic that postData is passed mostly or simple fallback
         const p = {
          id: postId,
          login: "User",
          title: "Post",
          content: "Loading...",
          image: null,
          avatar: "/default-avatar.png",
          likesCount: 0,
          dislikeCount: 0,
          time: ""
        };
        setPost(p);
    }
  }, [postId, postData]);

  // Fetch comments and reaction status
  useEffect(() => {
    if (!postId) return;

    const loadData = async () => {
        try {
            // 1. Fetch Comments
            const cRes = await apiFetch(`/posts/${postId}/comments`);
            if (cRes.ok) {
                const cData = await cRes.json();
                setComments(Array.isArray(cData) ? cData : []);
            }

            // 2. Fetch Reaction Status (and fresh counts) if possible
            // The backend /reactions endpoint is POST for toggle, but usually we get status with the post.
            // If postData didn't have userReaction, we might not know it. 
            // However, let's assume for now we start neutral or use what's in postData if extended.
            // Actually, we can just fetch the post details again to get accurate counts and user reaction if the endpoint supports it.
            // Based on repository.go SELECT_ALL_POSTS, it returns user_reaction.
            // Let's try to fetch the single post to get reaction state if postData is missing it.
            // Or just proceed.
            
        } catch (e) {
            console.error("Failed to load comments", e);
        }
    };
    loadData();
  }, [postId]);

  const handleReaction = async (type) => {
    // Optimistic update
    const oldReaction = userReaction;
    const oldLikes = likesCount;
    const oldDislikes = dislikeCount;

    let newLikes = likesCount;
    let newDislikes = dislikeCount;
    let newReaction = type;

    if (oldReaction === type) {
        // Untoggle
        newReaction = null;
        if (type === 'like') newLikes--;
        else newDislikes--;
    } else {
        // Switch or new
        if (oldReaction === 'like') newLikes--;
        if (oldReaction === 'dislike') newDislikes--;
        
        if (type === 'like') newLikes++;
        else newDislikes++;
    }

    setLikesCount(newLikes);
    setDislikeCount(newDislikes);
    setUserReaction(newReaction);

    try {
        if (!submitReactionUrl) return;

        const body = { reaction: type };
        if (endpoints?.submitReaction) {
            body.group_post_id = postId; // Assuming group post endpoint expects this
            // Or if generic, we might need to adjust body structure based on endpoint
            // For now, let's assume the endpoint handles "post_id" or we pass "post_id"
            // If it's a group post, we might want to send "post_id" as the ID.
            body.post_id = postId; 
        } else {
             body.post_id = postId;
        }

        const res = await apiFetch(submitReactionUrl, {
            method: "POST",
            body: JSON.stringify(body)
        });
        if (res.ok) {
            const data = await res.json();
            setLikesCount(data.likes_count);
            setDislikeCount(data.dislikes_count);
            setUserReaction(data.userReaction || null);
        } else {
            // Revert on error
            setLikesCount(oldLikes);
            setDislikeCount(oldDislikes);
            setUserReaction(oldReaction);
        }
    } catch (e) {
        setLikesCount(oldLikes);
        setDislikeCount(oldDislikes);
        setUserReaction(oldReaction);
    }
  };

  const handleCommentSubmit = async (e) => {
    e.preventDefault();
    if (!newComment.trim() && !commentImage) return;

    const formData = new FormData();
    formData.append("content", newComment);
    if (commentImage) {
        formData.append("image", commentImage);
    }

    try {
        const res = await apiFetch(`/posts/${postId}/comments`, {
            method: "POST",
            body: formData
        });

        if (res.ok) {
            // Refresh comments
            setNewComment("");
            setCommentImage(null);
            if (fileInputRef.current) fileInputRef.current.value = "";

            const cRes = await apiFetch(`/posts/${postId}/comments`);
            if (cRes.ok) {
                const cData = await cRes.json();
                setComments(Array.isArray(cData) ? cData : []);
            }
        }
    } catch (e) {
        console.error("Failed to submit comment", e);
    }
  };

  if (!post) return null;

  // Use passed data or fallback
  const displayImage = post.imageUrl || post.image;
  const displayAvatar = post.avatar;
  const displayName = post.author || post.login;
  const displayContent = post.content;
  const displayTime = post.ago || post.time || "Just now";

  return (
    <div className={style.postContainer}>
      <div className={style.postModal}>
        <button className={style.closeBtn} onClick={onClose}>
            <X />
        </button>

        <div className={style.post}>
            <div className={style.postBgContainer}>
                {displayImage ? (
                    <img className={style.postBg} src={displayImage} alt="" />
                ) : (
                    <div className={style[bgClass]}></div>
                )}
            </div>

            <div className={style.postHeader}>
                <img className={style.avatar} src={displayAvatar} alt="" />
                <div className={style.info}>
                    <div className={style.name}>{displayName}</div>
                    <div className={style.time}>{displayTime}</div>
                </div>
            </div>

            <div style={{ flex: 1 }}></div>

            <div className={style.postContentWrapper}>
                <div className={style.postContent}>
                    {displayContent}
                </div>
                <div className={style.postFooter}>
                    <button 
                        className={`${style.iconBtn} ${userReaction === 'like' ? style.active : ''}`}
                        onClick={() => handleReaction('like')}
                    >
                        <Heart size={16} fill={userReaction === 'like' ? "currentColor" : "none"} />
                        <span>{likesCount}</span>
                    </button>
                    <button 
                        className={`${style.iconBtn} ${userReaction === 'dislike' ? style.active : ''}`}
                        onClick={() => handleReaction('dislike')}
                    >
                        <ThumbsDown size={16} fill={userReaction === 'dislike' ? "currentColor" : "none"} />
                        <span>{dislikeCount}</span>
                    </button>
                </div>
            </div>
        </div>
      </div>

      <div className={style.addComment}>
        <form onSubmit={handleCommentSubmit}>
            <div className={style.commentBoxContainer}>
                <textarea 
                    placeholder="Write a comment ..."
                    value={newComment}
                    onChange={(e) => setNewComment(e.target.value)}
                ></textarea>
                <div className={style.formActions}>
                    <input 
                        type="file" 
                        ref={fileInputRef}
                        style={{ display: 'none' }}
                        onChange={(e) => setCommentImage(e.target.files[0])}
                    />
                    <button type="button" onClick={() => fileInputRef.current?.click()} className={style.imageBtn}>
                        <ImageIcon size={20} color={commentImage ? "#007bff" : "#666"} />
                    </button>
                    <input type="submit" value="Add Comment" />
                </div>
            </div>
        </form>
      </div>

      {comments.map(comment => (
        <div key={comment.id} className={style.comment}>
            <div className={style.commentHeader}>
                <img className={style.avatar} src={toMediaUrl(comment.avatar) || "/default-avatar.png"} alt="" />
                <div className={style.info}>
                    <div className={style.name}>{comment.username || comment.login || "User"}</div>
                    <div className={style.time}>{new Date(comment.created_at).toLocaleString()}</div>
                </div>
            </div>
            <div className={style.commentContent}>
                {comment.text || comment.content || comment.body}
            </div>
            {comment.image_path && (
                <div className={style.commentImage}>
                    <img src={toMediaUrl(comment.image_path)} alt="Comment attachment" />
                </div>
            )}
        </div>
      ))}
    </div>
  );
}
