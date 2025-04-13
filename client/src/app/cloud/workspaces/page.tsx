import React from 'react';
import TypeBlock from '@/app/components/TypeBlock';


const page = () => {


  return (
    <div className="flex flex-wrap gap-4 p-4">


    <div className="flex-1">
      <TypeBlock type="photo" />
    </div>
    <div className="flex-1">
      <TypeBlock type="video" />
    </div>
    <div className="flex-1">
      <TypeBlock type="text" />
    </div>
    <div className="flex-1">
      <TypeBlock type="unknown" />
    </div>

  </div>
  );
};

export default page;