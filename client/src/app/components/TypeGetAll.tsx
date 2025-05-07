'use client';
import { useEffect, useState } from "react";
import { FileData } from "@/app/api/models/FileData";
import CloudService from "../api/services/CloudServices";
import FileCard from "@/app/ui/FileCard";
import TypeFileIcon from "../ui/TypeFileIcon";
import { Loader2 } from 'lucide-react';
import dynamic from 'next/dynamic';

export default function TypeBlock() {
  const [filesByType, setFilesByType] = useState<Record<string, FileData[]>>({});
  const [isLoading, setIsLoading] = useState(true);
  const [isError, setIsError] = useState(false);
  const types = ['text', 'photo', 'video', 'unknown'];
  const totalSpace = 10; 


  const DiskUsageChart = dynamic(() => import('../ui/DiskUsageChart'), {
    ssr: false,
    loading: () => <p className="text-center text-gray-500">Загрузка графика...</p>,
  });

  





  useEffect(() => {
    const fetchAllTypes = async () => {
      try {
        const result: Record<string, FileData[]> = {};

        for (const type of types) {
          const response = await CloudService.getAllCloud(type);
          const fileData = response.data.file_data;

          if (Array.isArray(fileData)) {
            result[type] = fileData.map((file: any) => ({
              obj_id: String(file.obj_id),
              name: String(file.name),
              url: String(file.url),
              created_at: String(file.created_at),
            }));
          } else {
            result[type] = [];
          }
        }

        setFilesByType(result);
      } catch (error) {
        console.error("Ошибка при получении данных:", error);
        setIsError(true);
      } finally {
        setIsLoading(false);
      }
    };

    fetchAllTypes();
  }, []);

  if (isLoading) return (
    <div className="inset-0 bg-white/70 backdrop-blur-sm z-10 flex items-center justify-center">
      <div className="flex flex-col items-center">
        <Loader2 className="w-10 h-10 text-blue-500 animate-spin mb-2" />
        <span className="text-gray-700 text-sm">Загрузка файлов...</span>
      </div>
    </div>
  );

  if (isError) return <p>Произошла ошибка при загрузке данных.</p>;

  // Рассчитываем количество файлов по типам
  const fileCounts = types.reduce((acc, type) => {
    acc[type] = filesByType[type]?.length || 0;
    return acc;
  }, {} as Record<string, number>);

  // Суммарный размер всех файлов
  const totalUsedSpace = Object.values(filesByType).flat().reduce((sum, file) => sum + (file.size || 0), 0);


  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 xl:grid-cols-2 gap-6 p-4">
      {types.map((type) => (
        <div
          key={type}
          className="bg-white border rounded-xl shadow-md p-4 flex flex-col justify-between w-full h-64"
        >
          <div className="flex items-center gap-2 mb-3">
            <TypeFileIcon type={type} />
            <h2 className="text-lg font-jetbrains text-blue-600 capitalize">{type}</h2>
          </div>

          {filesByType[type]?.length === 0 ? (
            <div className="text-gray-500 text-center flex-1 flex flex-col items-center justify-center">
              <p className="text-xl">📂 Нет файлов</p>
            </div>
          ) : (
            <div className="space-y-2 overflow-auto flex-1">
              {filesByType[type].map((item) => (
                <div
                  key={item.obj_id}
                  className="text-xl border rounded p-2 flex justify-between items-center"
                ><a href={item.url} target="_blank" rel="noopener noreferrer" className="text-blue-500 text-sm">
                  <span className="truncate max-w-[150px]">{item.name}</span>
                  </a>
                </div>
              ))}
            </div>
          )}

          
        </div>
      ))}
      
<div className="mt-10 p-4">
  <h2 className="text-xl font-semibold mb-3">📁 Все файлы</h2>
  <div className="flex flex-wrap gap-3 overflow-x-auto">
    {Object.values(filesByType).flat().map((file) => (
      <a
        key={file.obj_id}
        href={file.url}
        target="_blank"
        rel="noopener noreferrer"
        className="px-3 py-2 bg-gray-100 rounded-lg border shadow text-sm hover:bg-blue-50 transition whitespace-nowrap"
      >
        {file.name}
      </a>
    ))}
  </div>
</div>
<DiskUsageChart
        fileCounts={fileCounts}
        totalSpace={totalSpace}
        totalUsedSpace={totalUsedSpace}
      />
    </div>

  );
}
