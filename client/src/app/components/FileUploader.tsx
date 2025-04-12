'use client'

import React, { useState } from 'react';
import CloudService from '../api/services/CloudServices';

const FileUploader = () => {
  const [files, setFiles] = useState<File[]>([]);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = event.target.files;
    if (selectedFiles && selectedFiles.length > 0) {
      setFiles(Array.from(selectedFiles));
    }
  };

  const handleUpload = async () => {
    if (files.length === 0) {
      setErrorMessage('Файл не выбран');
      return;
    }

    const formData = new FormData();
    files.forEach((file) => {
      formData.append('file', file); // именно оригинальный File
    });

    try {
      setIsLoading(true);
      const response = await CloudService.createCloud(formData);

      if (response.status !== 200) {
        throw new Error('Ошибка при загрузке файла');
      }
      setSuccessMessage('Файл успешно загружен!');
      setErrorMessage(null);
    } catch (error) {
      setErrorMessage('Не удалось загрузить файл.');
      setSuccessMessage(null);
    } finally {
      setIsLoading(false);
    }
  };

  const renderFilePreview = (file: File) => {
    // Если файл изображение, показываем картинку
    if (file.type.startsWith('image/')) {
      return <img src={URL.createObjectURL(file)} alt={file.name} className="match-image" />;
    }
    // Для других типов файлов показываем ссылку для скачивания
    return <a href={URL.createObjectURL(file)} download={file.name}>Скачать {file.name}</a>;
  };

  return (
    <div className="p-4 max-w-md mx-auto bg-white rounded shadow">
      <h2 className="text-xl font-bold mb-4">Загрузка файлов</h2>

      <input
        type="file"
        multiple
        onChange={handleFileChange}
      />

      <button
        onClick={handleUpload}
        disabled={isLoading}
        className="mt-2 bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded"
      >
        {isLoading ? 'Загрузка...' : 'Загрузить'}
      </button>

      {successMessage && <p className="text-green-600 mt-2">{successMessage}</p>}
      {errorMessage && <p className="text-red-600 mt-2">{errorMessage}</p>}

      {files.length > 0 && (
        <div className="mt-4">
          <h3>Файлы для загрузки:</h3>
          <ul>
            {files.map((file, index) => (
              <li key={index}>
                <p>{file.name}</p>
                {renderFilePreview(file)} {/* Универсальная отрисовка файла */}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default FileUploader;
