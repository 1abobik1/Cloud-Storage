import React from 'react';
import { Post } from '@/app/lib/posts';
import Link from 'next/link';


interface PostCardProps {
  post: Post;
  onDelete:(id:number) =>  Promise<void>;
  onEdit:(id:number)=>Promise<void>;
}



const PostCard: React.FC<PostCardProps> = ({ post, onDelete, onEdit}) => {
  return (
  <div className="bg-white shadow-md rounded-lg p-4 mb-4 hover:bg-gray-100 transition">
    <div className='w-[50%]'>
      <Link href={`/cloud/posts/${post.id}`}>
          <h1 className="font-bold mb-2 text-fuchsia-700 text-2xl">{post.title}</h1>
          <p className="text-black line-clamp-1">{post.body}</p>
           </Link>
          
       <Link href={`/cloud/createUpdatePosts/${post.id}`}>Редактировать </Link>
          <button onClick={() => onDelete(post.id!)} className="mt-2 bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">Удалить</button>
      </div>
  </div>
      
      
   
    
  );
};

export default PostCard;