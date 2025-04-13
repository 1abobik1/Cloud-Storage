'use client';
import {useEffect, useState} from "react";
import {FileData} from "@/app/api/models/FileData";
import CloudService from "../api/services/CloudServices";
import FileCard from "@/app/ui/FileCard";


export default function TypeBlock({ type }) {
  const [file, setFile] = useState<FileData[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isError, setIsError] = useState<boolean>(false);
  const [timeSort,setTimeSort] = useState<boolean>(false)
  // Состояния для фильтрации и сортировки
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc'); // сортировка по возрастанию или убыванию
  const [filteredFiles, setFilteredFiles] = useState<FileData[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await CloudService.getAllCloud(type);
        const fileData = response.data.file_data;

        if (Array.isArray(fileData)) {
          const files: FileData[] = fileData.map((file: any) => ({
            obj_id: String(file.obj_id),
            name: String(file.name),
            url: String(file.url),
            created_at: String(file.created_at),
          }));
          setFile(files);
          setFilteredFiles(files); // Изначально отображаем все файлы
        } else {
          setFile([]);
        }
      } catch (error) {
        console.error("Ошибка при получении данных:", error);
        setIsError(true);
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, [type]);

  // Функция сортировки по дате
  const sortFiles = (order: 'asc' | 'desc') => {
    const sortedFiles = [...filteredFiles].sort((a, b) => {
      const dateA = new Date(a.created_at).getTime();
      const dateB = new Date(b.created_at).getTime();
      return order === 'asc' ? dateA - dateB : dateB - dateA;
    });
    setFilteredFiles(sortedFiles);
  };

  // Обработчик изменения сортировки
  const handleSortChange = (order: 'asc' | 'desc') => {
    setSortOrder(order);
    sortFiles(order); // Пересортировываем файлы
  };

  if (isLoading) return <p>Загрузка...</p>;
  if (isError) return <p>Произошла ошибка при загрузке данных.</p>;

  const handleDelete = (id: string) => {
    setFilteredFiles(prevFiles => prevFiles.filter(file => file.obj_id !== id)); // Убираем удаленный файл из состояния
  };

  return (
    <div className="p-4 mx-auto bg-white rounded shadow w-100vw">
      <h2 className="text-xl font-bold mb-4">{type}</h2>


      {timeSort &&(<div className="mb-4">
        <button
          onClick={() =>{ handleSortChange('desc');
            setTimeSort(true)}
          }
          className="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded">
          Сортировать по дате (по убыванию)
        </button>
      </div>)}

      {/* Отображаем отсортированные файлы */}
      {filteredFiles.map((item) => (
        <FileCard
          key={item.obj_id}
          obj_id={item.obj_id}
          name={item.name}
          url={item.url}
          created_at={item.created_at}
          type={type}
          onDelete={handleDelete}
        />
      ))}
    </div>
  );
};
