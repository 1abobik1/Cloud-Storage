"use client"; // Делаем компонент клиентским

import { useEffect, useState } from "react";
import {  getPhoto, Photo } from "@/app/lib/photo";
import PhotoCard from '@/app/ui/PhotoCard';


export default function PhotoPack() {
  const [photo] = useState<Photo[]>([]);


  useEffect(() => {
    getPhoto()
      .then()
      .catch((err) => console.error(err));
  }, []);

  if (photo.length === 0) return <p>Загрузка...</p>;

// const handleDeletePost = async (id:number) =>{
//   try{
//     await deletePost(id);
//     setPosts((prevPosts)=> prevPosts.filter((post)=> post.id !== id))
//   }catch(error){
//     console.error('Failed to delete post:', error)
//   }
// }



return (
  <div>
     {photo.map((photo) => (
        <PhotoCard key={photo.id} photo={photo}  />
      ))
    }
  </div>
);
}
