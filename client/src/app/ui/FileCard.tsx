import React, {useState} from 'react';
import CloudService from '../api/services/CloudServices';

import {ArrowDownTrayIcon, TrashIcon} from '@heroicons/react/24/outline';
import ModalDelete from './ModalDelete';

import TypeFileIcon from './TypeFileIcon';
import Link from 'next/link';




export type FileCardData = {
  name: string;
  created_at: string;
  obj_id: string;
  url: string;
  type: string;
  onDelete: (obj_id: string) => void;
};

const FileCard: React.FC<FileCardData> = ({ obj_id, created_at, name, url, type, onDelete }) => {
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const handleDownload = async () => {
    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error('Ошибка при скачивании файла');
      }

     
      
    } catch (error) {
      console.error('Ошибка при скачивании или расшифровке файла:', error);
    }
  };

  const handleDelete = async () => {
    try {
      const response = await CloudService.deleteFile(type, obj_id); // Вызов метода с type и obj_id
      console.log('Файл успешно удален:', response);
      setIsModalOpen(false);
      onDelete(obj_id); // Удаляем файл из списка в родительском компоненте
    } catch (error) {
      console.error('Ошибка при удалении файла:', error);
      // Можно показать сообщение об ошибке пользователю
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString(); // Выводит дату в формате, зависящем от локали
  };

const DeleteIcon =TrashIcon
const DownloadIcon = ArrowDownTrayIcon



const handleOpenModal = () => {
  setIsModalOpen(true); // Открываем модальное окно при нажатии "Удалить"
};

const handleCloseModal = () => {
  setIsModalOpen(false); // Закрываем модальное окно
};

function truncateText(text: string, maxLength: number): string {
  if (text.length <= maxLength) {
    return text;
  }
  return text.slice(0, maxLength) + '…'; // добавляет троеточие
}

return (
  <div className="p-4 bg-white border-t border-b border-gray-200 w-full">
    <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">

      {/* Название файла с иконкой */}
      <div className="flex items-center gap-2 w-full min-w-0">
        <TypeFileIcon type={type}  />
        <Link
          href={url}
          target="_blank"
          rel="noopener noreferrer"
          className=" break-words  leading-snug"
        >
          <>
  {/* Для экранов шире 600px — полное имя */}
  <span className="hidden sm:inline break-all  leading-snug">{name}</span>

  {/* Для экранов до 600px — сокращённое имя */}
  <span className="inline sm:hidden">{truncateText(name, 20)}</span>
</>

        </Link>
      </div>

      {/* Блок с датой и кнопками */}
      <div className="flex justify-between items-center sm:justify-end gap-4 w-full sm:w-auto">

        {/* Дата — видно только на >=sm */}
        <div className="text text-gray-500 ">
          {formatDate(created_at)}
        </div>

        {/* Кнопки */}
        <div className="flex items-center gap-2">
          <button
            onClick={handleDownload}
            className="w-6 h-6 flex items-center justify-center"
          >
            <DownloadIcon />
          </button>
          <button
            onClick={handleOpenModal}
            className="w-6 h-6 flex items-center justify-center"
          >
            <DeleteIcon />
          </button>
        </div>
      </div>
    </div>

    {isModalOpen && (
      <ModalDelete
        message="Вы уверены, что хотите удалить этот файл?"
        onConfirm={handleDelete}
        onCancel={handleCloseModal}
      />
    )}
  </div>
);

  
};

export default FileCard;
