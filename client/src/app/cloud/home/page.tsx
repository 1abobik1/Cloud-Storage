'use client'
import React from 'react';

import dynamic from 'next/dynamic';

const page = () => {

  const TypeGetAll = dynamic(() => import('@/app/components/TypeGetAll'), {
    ssr: false,
    loading: () => <p className="text-center text-gray-500">Загрузка графика...</p>,
  });



  return (
    <div className="flex flex-wrap gap-4 p-4">


    <div className="flex-1">
      <TypeGetAll/>
    </div>
    

  </div>
  );
};
export default page;

