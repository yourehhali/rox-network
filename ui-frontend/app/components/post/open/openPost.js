import style from "./openPost.module.css";
import { X, Heart, ThumbsDown, Image as ImageIcon } from "../../icons";
import { useState, useEffect, useRef } from "react";
import { apiFetch, toMediaUrl } from "../../../../lib/api";
import { useAuth } from "../../../../context/AuthContext";

export default function OpenPost({ postId: propPostId, onClose, postData, endpoints, onUpdatePost }) {
  const { user } = useAuth();
  const postId = propPostId || postData?.id;
  
  // Default endpoints
  const commentsUrl = endpoints?.comments || `/posts/${postId}/comments`;
  const submitCommentUrl = endpoints?.submitComment || `/posts/${postId}/comments`;
  const submitReactionUrl = endpoints?.submitReaction || "/reactions";

  const [post, setPost] = useState(null);
  const [comments, setComments] = useState([]);
  const [commentCursor, setCommentCursor] = useState(0);
  const [hasMoreComments, setHasMoreComments] = useState(true);
  const [loadingComments, setLoadingComments] = useState(false);
  const [bgClass, setBgClass] = useState("");
  const [newComment, setNewComment] = useState("");
  const [commentImage, setCommentImage] = useState(null);
  const [likesCount, setLikesCount] = useState(0);
  const [dislikeCount, setDislikeCount] = useState(0);
  const [userReaction, setUserReaction] = useState(null);
  const fileInputRef = useRef(null);

  const loadComments = async (cursor = 0, reset = false) => {
    if (loadingComments && !reset) return;
    setLoadingComments(true);
    try {
        const res = await apiFetch(`${commentsUrl}?cursor=${cursor}&limit=10`);
        if (res.ok) {
            const data = await res.json();
            const newComments = Array.isArray(data) ? data : [];
            
            if (reset) {
                setComments(newComments);
            } else {
                setComments(prev => {
                    // Filter duplicates just in case
                    const existingIds = new Set(prev.map(c => c.id));
                    const uniqueNew = newComments.filter(c => !existingIds.has(c.id));
                    return [...prev, ...uniqueNew];
                });
            }
            
            if (newComments.length < 10) {
                setHasMoreComments(false);
            } else {
                setHasMoreComments(true);
            }
        }
    } catch (e) {
        console.error("Failed to load comments", e);
    } finally {
        setLoadingComments(false);
    }
  };

  const handleLoadMoreComments = () => {
    if (loadingComments) return;
    const lastComment = comments[comments.length - 1];
    const nextCursor = lastComment ? lastComment.id : 0;
    setCommentCursor(nextCursor);
    loadComments(nextCursor, false);
  };

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
                const res = await apiFetch(`/posts/${postId}`, { cache: "no-store" }); // Assuming this endpoint exists
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
            setComments([]);
            setCommentCursor(0);
            setHasMoreComments(true);
            await loadComments(0, true);

            // 2. Fetch Fresh Post Details (for reaction state)
            const pRes = await apiFetch(`/GetPost?id=${postId}`, { cache: "no-store" });
            if (pRes.ok) {
                const pData = await pRes.json();
                if (pData) {
                    setLikesCount(pData.likes_count || 0);
                    setDislikeCount(pData.dislikes_count || 0);
                    setUserReaction(pData.userReaction || null); // null if no reaction
                    
                    // Update parent if data changed
                    if (onUpdatePost) {
                         onUpdatePost({
                             ...postData,
                             likesCount: pData.likes_count || 0,
                             dislikeCount: pData.dislikes_count || 0,
                             userReaction: pData.userReaction || null
                         });
                    }
                }
            }
            
        } catch (e) {
            console.error("Failed to load data", e);
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
    
    // Notify parent optimistically (optional, but good for responsiveness)
    if (onUpdatePost) {
        onUpdatePost({
            ...postData,
            likesCount: newLikes,
            dislikeCount: newDislikes,
            userReaction: newReaction
        });
    }

    try {
        if (!submitReactionUrl) return;

        const body = { reaction: type };
        if (endpoints?.submitReaction) {
            body.group_post_id = postId; 
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
            
            // Notify parent with confirmed data
            if (onUpdatePost) {
                onUpdatePost({
                    ...postData,
                    likesCount: data.likes_count,
                    dislikeCount: data.dislikes_count,
                    userReaction: data.userReaction || null
                });
            }
        } else {
            // Revert on error
            setLikesCount(oldLikes);
            setDislikeCount(oldDislikes);
            setUserReaction(oldReaction);
            if (onUpdatePost) {
                onUpdatePost({
                    ...postData,
                    likesCount: oldLikes,
                    dislikeCount: oldDislikes,
                    userReaction: oldReaction
                });
            }
        }
    } catch (e) {
        setLikesCount(oldLikes);
        setDislikeCount(oldDislikes);
        setUserReaction(oldReaction);
        if (onUpdatePost) {
            onUpdatePost({
                ...postData,
                likesCount: oldLikes,
                dislikeCount: oldDislikes,
                userReaction: oldReaction
            });
        }
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

            setCommentCursor(0);
            setHasMoreComments(true);
            await loadComments(0, true);
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
                    maxLength={1000}
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

      {hasMoreComments && (
        <div style={{ textAlign: "center", padding: "10px" }}>
            <button 
                onClick={handleLoadMoreComments} 
                disabled={loadingComments}
                style={{
                    padding: "8px 16px",
                    borderRadius: "20px",
                    border: "none",
                    background: "#e4e6eb",
                    cursor: loadingComments ? "default" : "pointer",
                    fontWeight: "600",
                    color: "#050505"
                }}
            >
                {loadingComments ? "Loading..." : "Load more comments"}
            </button>
        </div>
      )}
    </div>
  );
}
