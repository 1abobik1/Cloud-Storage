'use client';

import React, { useRef, useState } from 'react';
import CloudService from '../api/services/CloudServices';
import { ArrowUpOnSquareIcon } from '@heroicons/react/24/outline';

const FileUploader = () => {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [toastMessage, setToastMessage] = useState<string | null>(null);
  const [toastType, setToastType] = useState<'success' | 'error' | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState<number>(0);

  const handleButtonClick = () => {
    inputRef.current?.click();
  };

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = event.target.files;
    if (!selectedFiles || selectedFiles.length === 0) return;

    const formData = new FormData();
    Array.from(selectedFiles).forEach((file) => {
      formData.append('files', file);
    });

    const MIN_SPINNER_DURATION = 2000;
    const startTime = Date.now();

    try {
      setIsLoading(true);
      setToastMessage('Загрузка файла...');
      setToastType(null);
      setUploadProgress(0);

      const response = await CloudService.uploadFiles(formData, {
        onUploadProgress: (progressEvent: ProgressEvent) => {
          const progress = Math.round((progressEvent.loaded * 100) / (progressEvent.total || 1));
          setUploadProgress(progress);
        },
        timeout: 10000, // 10 секунд
      });

      const elapsedTime = Date.now() - startTime;
      const remainingTime = Math.max(0, MIN_SPINNER_DURATION - elapsedTime);
      await new Promise((resolve) => setTimeout(resolve, remainingTime));

      if (response.status !== 200) {
        throw new Error('Ошибка при загрузке файла');
      }

      setToastMessage('Файл успешно загружен!');
      setToastType('success');
    } catch (error: any) {
      console.error(error);
      if (error.code === 'ECONNABORTED') {
        setToastMessage('Время ожидания запроса истекло.');
      } else {
        setToastMessage('Не удалось загрузить файл.');
      }
      setToastType('error');
    } finally {
      setIsLoading(false);
      setUploadProgress(0);
      if (inputRef.current) inputRef.current.value = '';
      setTimeout(() => {
        setToastMessage(null);
        setToastType(null);
      }, 3000);
    }
  };

  const LinkIcon = ArrowUpOnSquareIcon;

  return (
    <>
      <div className="max-w-md mx-auto text-center">
        <input
          type="file"
          multiple
          ref={inputRef}
          style={{ display: 'none' }}
          onChange={handleFileChange}
        />

        <button
          onClick={handleButtonClick}
          disabled={isLoading}
          className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded"
        >
          <div className="flex items-center justify-center">
            <LinkIcon className="w-6 text-white mr-1" />
            <span>{isLoading ? 'Загрузка...' : 'Загрузить файл'}</span>
          </div>
        </button>

        {isLoading && (
          <div className="mt-4">
            <div className="w-full bg-gray-200 h-2 rounded">
              <div
                className="h-2 bg-blue-500 rounded transition-all duration-300"
                style={{ width: `${uploadProgress}%` }}
              />
            </div>
            <div className="text-sm text-gray-700 mt-1">{uploadProgress}%</div>
          </div>
        )}
      </div>

      {toastMessage && (
        <div
          className={`fixed bottom-4 right-4 w-80 px-4 py-3 rounded shadow-lg z-50 transition-all duration-300 text-white
          ${toastType === 'success' ? 'bg-green-500' : toastType === 'error' ? 'bg-red-500' : 'bg-blue-500'}`}
        >
          <div className="flex items-center space-x-3">
            {isLoading && (
              <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            )}
            <div className="text-sm">{toastMessage}</div>
          </div>
        </div>
      )}
    </>
  );
};

export default FileUploader;
