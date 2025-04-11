import React, { useState } from 'react';
import { FileData } from '@/app/api/models/FileData';

const FileUploader = () => {
  // Тип состояния изменен на массив FileData[]
  const [files, setFiles] = useState<FileData[]>([]);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = event.target.files;
    if (selectedFiles) {
      // Преобразуем файлы в массив данных типа FileData[]
      const fileArray: FileData[] = Array.from(selectedFiles).map((file) => ({
        name: file.name,
        created_at: new Date().toISOString(), // или другая логика для created_at
        obj_id: file.name, // пример obj_id, можно заменить на реальное значение
        url: URL.createObjectURL(file), // временный URL
      }));
      setFiles(fileArray); // Обновление состояния
    }
  };

  const handleUpload = async () => {
    if (files.length === 0) {
      setErrorMessage('Файл не выбран');
      return;
    }

    const formData = new FormData();
    files.forEach((file) => {
      formData.append('file', file as any); // Преобразовать в Blob, если необходимо
    });

    try {
      setIsLoading(true);
      const response = await fetch('http://localhost:8081/files/one', {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) throw new Error('Ошибка при загрузке файла');

      const data = await response.json();
      setSuccessMessage('Файл успешно загружен!');
      setErrorMessage(null);
    } catch (error) {
      setErrorMessage('Не удалось загрузить файл.');
      setSuccessMessage(null);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="p-4 max-w-md mx-auto bg-white rounded shadow">
      <h2 className="text-xl font-bold mb-4">Загрузка файлов</h2>

      <input type="file" multiple accept="image/png, image/jpeg" onChange={handleFileChange} />
      <button
        onClick={handleUpload}
        disabled={isLoading}
        className="mt-2 bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded"
      >
        {isLoading ? 'Загрузка...' : 'Загрузить'}
      </button>

      {successMessage && <p className="text-green-600 mt-2">{successMessage}</p>}
      {errorMessage && <p className="text-red-600 mt-2">{errorMessage}</p>}

      <div className="mt-4">
        <h3>Загруженные файлы:</h3>
        <ul>
          {files.map((file, index) => (
            <li key={index}>
              <p>{file.name}</p>
              <img src={file.url} alt={file.name} className="match-image" />
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default FileUploader;
