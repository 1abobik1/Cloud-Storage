import React from 'react';
import {FileData} from '@/app/api/models/FileData'
import Link from 'next/link';


const FileCard: React.FC<FileData> = ({ obj_id, created_at, name, url }) => {
  const handleDownload = async () => {
    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error('Ошибка при скачивании файла');
      }
  
      
      const blob = await response.blob();
      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = name;
      link.click();
    } catch (error) {
      console.error('Ошибка скачивания файла:', error);
    }
  };

  return (
    <div className="bg-white shadow-md rounded-lg p-4 mb-4 hover:bg-gray-100 transition border-b">
      <div className="w-[50%]">
      <button
    onClick={handleDownload}
    className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded"
  >
    Скачать {name}
  </button>
      </div>
    </div>
  );
};




export default FileCard;