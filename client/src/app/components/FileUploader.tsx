'use client'

import React, {useRef, useState} from 'react';
import CloudService from '../api/services/CloudServices';
import {ArrowUpOnSquareIcon} from '@heroicons/react/24/outline';

const FileUploader = () => {
    const inputRef = useRef<HTMLInputElement | null>(null);
    const [toastMessage, setToastMessage] = useState<string | null>(null);
    const [toastType, setToastType] = useState<'success' | 'error' | null>(null);
    const [isLoading, setIsLoading] = useState(false);

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
        try {
            setIsLoading(true);
            setToastMessage(null);
            setToastType(null);

            const response = await CloudService.uploadFiles(formData);

            if (response.status !== 200) {
                throw new Error('Ошибка при загрузке файла');
            }

            setToastMessage('Файл успешно загружен!');
            setToastType('success');
        } catch (error) {
            console.error(error)
            setToastMessage(`Не удалось загрузить файл.`);
            setToastType('error');
        } finally {
            setIsLoading(false);
            if (inputRef.current) inputRef.current.value = '';

            // Автоматическое скрытие через 3 секунды
            setTimeout(() => {
                setToastMessage(null);
                setToastType(null);
            }, 3000);
        }
    };
    const LinkIcon = ArrowUpOnSquareIcon
    return (
        <>
            <div className=" max-w-md mx-auto  text-center">

                <input
                    type="file"
                    multiple
                    ref={inputRef}
                    style={{display: 'none'}}
                    onChange={handleFileChange}
                />

                <button
                    onClick={handleButtonClick}
                    disabled={isLoading}
                    className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded"
                >
                    <div className='flex'><LinkIcon className="w-6  text-white mr-1"/>
                        <div className='mt-1'>{isLoading ? 'Загрузка...' : 'Загрузить файл'}</div>
                    </div>
                </button>
            </div>

            {toastMessage && (
                <div
                    className={`
      fixed bottom-4 right-4 w-72 px-4 py-3 rounded shadow-lg z-50 transition-all duration-300
      text-white flex items-center justify-between
      ${toastType === 'success' ? 'bg-green-500' : toastType === 'error' ? 'bg-red-500' : 'bg-blue-500'}
    `}
                >
                    <span className="text-sm">{toastMessage}</span>

                    {isLoading && (
                        <svg
                            className="animate-spin h-5 w-5 ml-3 text-white"
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                        >
                            <circle
                                className="opacity-25"
                                cx="12"
                                cy="12"
                                r="10"
                                stroke="currentColor"
                                strokeWidth="4"
                            />
                            <path
                                className="opacity-75"
                                fill="currentColor"
                                d="M4 12a8 8 0 018-8v8H4z"
                            />
                        </svg>
                    )}
                </div>
            )}

        </>
    );
};

export default FileUploader;
