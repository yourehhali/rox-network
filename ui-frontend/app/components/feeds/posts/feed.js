"use client";
import style from "./feed.module.css";
import PostCard from "../../post/card/card";
import OpenPost from "../../post/open/openPost";
import { useState, useEffect } from 'react';
import { useSearchParams } from "next/navigation";
import { apiFetch, toMediaUrl } from '../../../../lib/api';

export default function PostsFeed({ userId: propUserId }) {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedPost, setSelectedPost] = useState(null);
  const params = useSearchParams();

  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const paramId = params?.get?.("id");
        const uid = propUserId || (paramId ? parseInt(paramId, 10) : null);
        const endpoint = uid ? `/GetCUser/?id=${uid}` : '/getposts';
        const res = await apiFetch(endpoint);
        if (res.ok) {
          const data = await res.json();
          const mappedPosts = (data || []).map(post => {
            const imageUrl = toMediaUrl(post.image_path) || undefined;
            const avatarUrl = toMediaUrl(post.avatar) || "/default-avatar.png";
            return {
              id: post.id,
              author: post.full_name || post.username,
              avatar: avatarUrl,
              content: post.content,
              imageUrl,
              ago: new Date(post.created_at).toLocaleDateString(),
              likesCount: post.likes_count || 0,
              dislikeCount: post.dislikes_count || 0,
              userReaction: post.userReaction,
            };
          });
          setPosts(mappedPosts);
        }
      } catch (err) {
        console.error("Failed to fetch posts", err);
      } finally {
        setLoading(false);
      }
    };

    fetchPosts();
  }, [propUserId, params]);

  if (loading) {
    return <div className={style.postsFeed}>Loading...</div>;
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
      {selectedPost && (
        <OpenPost 
          postData={selectedPost} 
          onClose={handleClosePost} 
        />
      )}
    </div>
  );
}
