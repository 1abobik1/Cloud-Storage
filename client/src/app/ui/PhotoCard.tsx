import React from 'react';
import { Photo } from '@/app/lib/photo';
import Link from 'next/link';


interface PhotoCardProps {
  photo: Photo;
}



const PhotoCard: React.FC<PhotoCardProps> = ({photo}) => {
  return (
  <div className="bg-white shadow-md rounded-lg p-4 mb-4 hover:bg-gray-100 transition border-b">
    <div className='w-[50%]'>
      <Link href={`/cloud/photo/${photo.id}`}>
          <h1 className="font-bold mb-2 text-fuchsia-700 text-2xl">{photo.data}</h1>
           </Link>


         {/* <button onClick={() => onDelete(post.id!)} className="mt-2 bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">Удалить</button>
      */}
      </div>
  </div>
      
      
   
    
  );
};

export default PhotoCard;