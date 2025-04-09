"use client"; // Делаем компонент клиентским

import { useEffect, useState } from "react";
import { deletePost, getPosts, Post } from "@/app/lib/posts";
import PostCard from '@/app/ui/PostCard';
import UpdatePostForm from "@/app/ui/UpdatePostForm";

export default function PostPage() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [editingPostId, setEditingPostId] = useState<number | null>(null);

  useEffect(() => {
    getPosts()
      .then(setPosts)
      .catch((err) => console.error(err));
  }, []);

  if (posts.length === 0) return <p>Загрузка...</p>;

const handleDeletePost = async (id:number) =>{
  try{
    await deletePost(id);
    setPosts((prevPosts)=> prevPosts.filter((post)=> post.id !== id))
  }catch(error){
    console.error('Failed to delete post:', error)
  }
}

const handleEditPost = (id: number): void => {
  setEditingPostId(id);
};

return (
  <div>
    {editingPostId ? (
      <UpdatePostForm postId={editingPostId} />
    ) : (
      posts.map((post) => (
        <PostCard key={post.id} post={post} onDelete={handleDeletePost} onEdit={handleEditPost} />
      ))
    )}
  </div>
);
}
