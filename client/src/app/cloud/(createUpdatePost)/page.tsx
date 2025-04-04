'use client';

import CreatePostForm from '@/app/ui/CreatePostForm';
import UpdatePostForm from '@/app/ui/UpdatePostForm';
import React, { useState } from 'react';

const Page = () => {
  const [showCreateForm, setShowCreateForm] = useState(true);

  const toggleForm = () => {
    setShowCreateForm(!showCreateForm);
  };

  return (
    <div>
      <button
        onClick={toggleForm}
        className="mb-4 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
      >
        {showCreateForm ? 'Показать форму обновления' : 'Показать форму создания'}
      </button>
      {showCreateForm ? <CreatePostForm /> : <UpdatePostForm />}
    </div>
  );
};

export default Page;