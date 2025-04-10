import React from 'react';
import { FileData } from '@/app/lib/cloud';
import Link from 'next/link';


interface FileDataCardProps {
  file: FileData;
}



const PhotoCard: React.FC<FileDataCardProps> = ({file}) => {
  return (
  <div className="bg-white shadow-md rounded-lg p-4 mb-4 hover:bg-gray-100 transition border-b">
    <div className='w-[50%]'>
      <Link href={`/cloud/photo/${file.obj_id}`}>
          <h1 className="font-bold mb-2 text-fuchsia-700 text-2xl">{file.url}</h1>
           </Link>


         {/* <button onClick={() => onDelete(post.id!)} className="mt-2 bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">Удалить</button>
      */}
      </div>
  </div>
      
      
   
    
  );
};

export default PhotoCard;