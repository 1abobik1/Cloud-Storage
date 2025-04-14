import React, {useState} from 'react';
import CloudService from '../api/services/CloudServices';

import {ArrowDownTrayIcon, TrashIcon} from '@heroicons/react/24/outline';
import ModalDelete from './ModalDelete';
import {cryptoHelper} from "@/app/api/utils/CryptoHelper";
import {getMimeTypeFromName} from "@/app/api/utils/getMimeTypeFromName";
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

      const encryptedBlob = await response.blob();
      const mimeType = getMimeTypeFromName(name); // получаем тип по имени файла
      const encryptedFile = new File([encryptedBlob], name, { type: mimeType });

      const decryptedBlob = await cryptoHelper.decryptFile(encryptedFile);

      const link = document.createElement('a');
      link.href = URL.createObjectURL(new Blob([decryptedBlob], { type: mimeType }));
      link.download = name;
      link.click();
      URL.revokeObjectURL(link.href);
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

  return (
    <div className="p-4 mx-auto bg-white border-t border-b border-gray-200 w-full">

      <div className="flex justify-between items-center">

        <div><TypeFileIcon type={type}/> <Link href={url} target="_blank" rel="noopener noreferrer" >{name}</Link> </div>
        <div className="w-60% flex items-center">
  <div className="mr-5 hidden sm:block">{formatDate(created_at)}</div> {/* Скрывается на мобильных */}
  <div>

            <button
              onClick={handleDownload}
              className=" w-6 mx-1"
            >
              <DownloadIcon/>
            </button>
            <button
              onClick={handleOpenModal} // Теперь handleDelete без параметров
              className='w-6  mx-1'
            >
              <DeleteIcon/>
            </button>
         </div>

        </div>
      </div>
      {/* Модальное окно для подтверждения удаления */}
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
