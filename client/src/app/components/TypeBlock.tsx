'use client'
import { useEffect, useState } from "react";


import PhotoCard from '@/app/ui/PhotoCard';
import { FileData } from "@/app/api/models/FileData";
import CloudService from "../api/services/CloudServices";
import FileUploader from "./FileUploader";
export default function TypeBlock() {
  const [file, setFile] = useState<FileData[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isError, setIsError] = useState<boolean>(false);
  


  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await CloudService.getAllCloud('photo');
        
        const files: FileData[] = response.data.file_data.map((file: any) => ({
          obj_id: String(file.obj_id),  
          name: String(file.name),             
          url: String(file.url), 
          created_at:String(file.created_at)             
        }));
        setFile(files);
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
      <h2 className="text-xl font-bold mb-4">Загрузка файла</h2>
<FileUploader/>
         {file.map((item) => (
  <PhotoCard obj_id={item.obj_id} name={item.name} url={item.url} created_at={item.created_at} />
))}
    </div>
  );
};
   
   