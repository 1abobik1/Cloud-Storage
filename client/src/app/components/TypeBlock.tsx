import { useEffect, useState } from "react";
import { getAllCloud, createCloud } from "@/app/lib/cloud";

import PhotoCard from '@/app/ui/PhotoCard';
import { FileData } from "@/app/lib/cloud";
export default function PhotoPack() {
  const [files, setFiles] = useState<FileData[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isError, setIsError] = useState<boolean>(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await getAllCloud('photo');
        setFiles(response); // Заполнение данных
      } catch (error) {
        console.error("Ошибка при получении данных:", error);
        setIsError(true); // Устанавливаем ошибку
      } finally {
        setIsLoading(false); // Завершаем загрузку
      }
    };

    fetchData();
  }, []); // Пустой массив зависимостей, чтобы эффект сработал только при монтировании компонента

  if (isLoading) return <p>Загрузка...</p>;
  if (isError) return <p>Произошла ошибка при загрузке данных.</p>;

  const handleFileChange = (event) => {
    setFile(event.target.files[0]);
  };

  // Обработчик отправки формы
  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!file) {
      alert('Пожалуйста, выберите файл для загрузки.');
      return;
    }
    setIsLoading(true);
    setIsError(false);

    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await fetch('http://localhost:8081/files/one', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        throw new Error('Ошибка при загрузке файла');
      }

      const result = await response.json();
      console.log('Файл успешно загружен:', result);
      // Здесь вы можете обновить состояние или выполнить другие действия
    } catch (error) {
      console.error('Ошибка при загрузке файла:', error);
      setIsError(true);
    } finally {
      setIsLoading(false);
    }
  };
  

  return (
    <form onSubmit={handleSubmit}>
      <input type="file" onChange={handleFileChange} />
      <button type="submit" disabled={isLoading}>
        {isLoading ? 'Загрузка...' : 'Загрузить'}
      </button>
      {isError && <p style={{ color: 'red' }}>Произошла ошибка при загрузке файла.</p>}
    </form>
  );
}
