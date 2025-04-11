'use client';
import React, { useState, useEffect } from 'react';
import { getPostById, updatePost } from '@/app/lib/cloud';

interface UpdatePostFormProps {
  postId: number;
}

const UpdatePostForm: React.FC<UpdatePostFormProps> = ({ postId }) => {
  const [title, setTitle] = useState('');
  const [body, setBody] = useState('');

  useEffect(() => {
    const fetchPost = async () => {
      try {
        const post = await getPostById(postId);
        setTitle(post.title);
        setBody(post.body);
      } catch (error) {
        console.error('Failed to fetch post:', error);
      }
    };

    fetchPost();
  }, [postId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const updatedPost = await updatePost({ id: postId, title, body });
      console.log('Post updated:', updatedPost);
      // Очистить форму после успешного обновления поста
      setTitle('');
      setBody('');
    } catch (error) {
      console.error('Failed to update post:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-white shadow-md rounded-lg p-4 mb-4">
      <h1>Редактирование поста</h1>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="title">
          Title
        </label>
        <input
          id="title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          required
        />
      </div>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="body">
          Body
        </label>
        <textarea
          id="body"
          value={body}
          onChange={(e) => setBody(e.target.value)}
          className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          required
        />
      </div>
      <button
        type="submit"
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
      >
        Update Post
      </button>
    </form>
  );
};

export default UpdatePostForm;