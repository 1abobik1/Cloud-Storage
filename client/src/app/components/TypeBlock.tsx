'use client'
import { useEffect, useState } from "react";
import { FileData } from "@/app/api/models/FileData";
import CloudService from "../api/services/CloudServices";
import FileUploader from "./FileUploader";
import FileCard from "@/app/ui/FileCard";

export default function TypeBlock({type}) {
  const [file, setFile] = useState<FileData[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isError, setIsError] = useState<boolean>(false);
  


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
        } else {
          console.warn("file_data не является массивом:", fileData);
          setFile([]); // или можно показать сообщение пользователю
        }
  
      } catch (error) {
        console.error("Ошибка при получении данных:", error);
        setIsError(true);
      } finally {
        setIsLoading(false);
      }
    };
  
    fetchData();
  }, []);
  

  if (isLoading) return <p>Загрузка...</p>;
  if (isError) return <p>Произошла ошибка при загрузке данных.</p>;

  


  
  return (
    <div className="p-4 max-w-md mx-auto bg-white rounded shadow">
      <h2 className="text-xl font-bold mb-4">Тип фалов: {type}</h2>
      
      {file.map((item) => (
  <FileCard
    key={item.obj_id} 
    obj_id={item.obj_id}
    name={item.name}
    url={item.url}
    created_at={item.created_at}
  />
))}

    </div>
  );
};
   
   