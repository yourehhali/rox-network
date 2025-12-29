"use client";
import style from "./feed.module.css";
import PostCard from "../../post/card/card";
import OpenPost from "../../post/open/openPost";
import { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from "next/navigation";
import { apiFetch, toMediaUrl } from '../../../../lib/api';
import { usePosts } from '../../../../context/PostsContext';
import { useAuth } from '../../../../context/AuthContext';

const LIMIT = 10;

export default function PostsFeed({ userId: propUserId, posts: propPosts, endpoints, groupId }) {
  const [internalPosts, setInternalPosts] = useState([]);
  const [loading, setLoading] = useState(!propPosts);
  const [selectedPost, setSelectedPost] = useState(null);
  const params = useSearchParams();
  const { refreshTrigger } = usePosts();
  const { user } = useAuth();
  const [locked, setLocked] = useState(false);


  const [cursor, setCursor] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);

  const posts = propPosts || internalPosts;

  const loadPosts = useCallback(async (currentCursor, isReset = false) => {
    if (propPosts) {
      setLoading(false);
      return;
    }

    if (isReset) {
      setLoading(true);
      setLocked(false);
      setHasMore(true);
    } else {
      setLoadingMore(true);
    }

    try {
      const paramId = params?.get?.("id");
      const uid = propUserId || (paramId ? parseInt(paramId, 10) : null);

      // Profile privacy check for user feeds (only on reset/first load)
      if (isReset && !groupId && uid) {
        try {
          const profRes = await apiFetch(`/profile?id=${uid}`);
          if (profRes.ok) {
            const profData = await profRes.json();
            const u = profData.user || profData;
            const isPrivate = !!u.is_private;
            const isFriend = !!profData.isfriend;
            const hasPosts = !!profData.has_posts;
            const isMyProfile = user && user.id === uid;
            
            if (!isMyProfile && isPrivate && !isFriend && hasPosts) {
              setLocked(true);
              setInternalPosts([]);
              setLoading(false);
              return;
            }
          }
        } catch (_) {
          // ignore privacy check errors; fall through to posts fetch
        }
      }

      let endpoint = '/getposts';
      if (groupId) {
        endpoint = `/GetGroupPosts?group_id=${groupId}`;
      } else if (uid) {
        endpoint = `/GetCUser/?id=${uid}`;
      }

      // Append pagination params
      const separator = endpoint.includes('?') ? '&' : '?';
      endpoint += `${separator}cursor=${currentCursor}&limit=${LIMIT}`;

      const res = await apiFetch(endpoint, { cache: "no-store" });
      if (res.ok) {
        const data = await res.json();
        const mappedPosts = (data || []).map(post => {
          const imageUrl = toMediaUrl(post.image_path || post.image) || undefined;
          const avatarUrl = toMediaUrl(post.avatar) || "/default-avatar.png";
          
          // Handle differences between group post and feed post
          let content = post.content;
          if (post.title) {
            content = `**${post.title}**\n\n${post.content}`;
          }

          return {
            id: post.id,
            author: post.full_name || post.username,
            avatar: avatarUrl,
            content: content,
            imageUrl,
            ago: new Date(post.created_at).toLocaleDateString(),
            likesCount: post.likes_count || 0,
            dislikeCount: post.dislikes_count || 0,
            userReaction: post.userReaction,
          };
        });

        if (isReset) {
          setInternalPosts(mappedPosts);
        } else {
          setInternalPosts(prev => {
            const existingIds = new Set(prev.map(p => p.id));
            const uniqueNew = mappedPosts.filter(p => !existingIds.has(p.id));
            return [...prev, ...uniqueNew];
          });
        }

        if (mappedPosts.length < LIMIT) {
          setHasMore(false);
        } else {
            setHasMore(true);
        }
      } else {
        // If fetch failed (e.g. 403 Forbidden), stop trying to load more
        setHasMore(false);
      }
    } catch (err) {
      console.error("Failed to fetch posts", err);
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  }, [propPosts, params, propUserId, groupId, user]);

  useEffect(() => {
    // Reset and load cursor 0
    setCursor(0);
    loadPosts(0, true);
  }, [loadPosts, refreshTrigger]);

  const handleLoadMore = () => {
    const lastPost = internalPosts[internalPosts.length - 1];
    const nextCursor = lastPost ? lastPost.id : 0;
    setCursor(nextCursor);
    loadPosts(nextCursor, false);
  };

  const handleUpdatePost = (updatedPost) => {
    setInternalPosts(prev => prev.map(p => p.id === updatedPost.id ? updatedPost : p));
    if (selectedPost && selectedPost.id === updatedPost.id) {
        setSelectedPost(updatedPost);
    }
  };

  if (loading && !posts.length) {
    return <div className={style.postsFeed}>Loading...</div>;
  }

  if (locked) {
    return (
      <div className={style.postsFeed}>
        <div className={style.postsFeedContainer}>
          <div style={{ width: "100%", padding: "24px", borderRadius: "12px", background: "#121212", border: "1px solid #222", color: "#aaa", textAlign: "center" }}>
            This profile is private. Follow to see posts.
          </div>
        </div>
      </div>
    );
  }

  const handleOpenPost = (post) => {
    setSelectedPost(post);
  };

  const handleClosePost = () => {
    setSelectedPost(null);
  };

  const buildPosts = posts.map((postData) => (
    <PostCard 
      key={postData.id} 
      postData={postData} 
      onOpen={() => handleOpenPost(postData)}
    />
  ));
  
  return (
    <div className={style.postsFeed}>
      <div className={style.postsFeedContainer}>{buildPosts}</div>
      
      {!loading && !propPosts && hasMore && (
        <div style={{ display: 'flex', justifyContent: 'center', padding: '20px' }}>
            <button 
                onClick={handleLoadMore} 
                disabled={loadingMore}
                style={{
                    padding: '10px 20px',
                    borderRadius: '20px',
                    background: '#333',
                    color: '#fff',
                    border: 'none',
                    cursor: loadingMore ? 'not-allowed' : 'pointer',
                    opacity: loadingMore ? 0.7 : 1
                }}
            >
                {loadingMore ? 'Loading...' : 'Load More'}
            </button>
        </div>
      )}

      {selectedPost && (
        <OpenPost 
          postData={selectedPost} 
          onClose={handleClosePost}
          endpoints={endpoints}
          onUpdatePost={handleUpdatePost}
        />
      )}
    </div>
  );
}
