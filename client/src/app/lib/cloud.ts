export type FileData = {
    name: string;
    created_at: string;
    obj_id: string;
    url: string;
  };
  
  export type CloudResponse = {
    file_data: FileData[];
    message: string;
    status: number;
  };
  
  
export async function getAllCloud(type:string): Promise<FileData[]> {
    try {
        const response = await fetch(`localhost:8081/files/all?type=${type}`);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    } catch (error) {
        console.error('Failed to fetch posts:', error);
        throw error;
    }
}
export async function getCloudById(id:number,type:string ){
    try{
        const response = await fetch(`localhost:8081/files/one?id=${id}&type=${type}`) 
        if (!response.ok){
            throw new Error('Network response was not ok')
        } 
                    return response.json()
    }catch(error){
        console.error('Failde to fetch post:', error)
        throw error
    }
}

export async function createCloud(file: File): Promise<void> {
    try {
      const formData = new FormData();
      formData.append('file', file);
  
      const response = await fetch('http://localhost:8081/files/one', {
        method: 'POST',
        body: formData,
      });
  
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
  
      console.log('Файл успешно загружен');
    } catch (error) {
      console.error('Ошибка при загрузке файла:', error);
      throw error;
    }
  }
  


// export async function deletePost(id:number){
//     try{
//         const response =await fetch(`https://jsonplaceholder.typicode.com/posts/${id}`,
//         {method:'DELETE'})
//         if(!response.ok){
//             throw new Error('Network response was not ok')
//         }
//     }catch(error){
//         console.log('failed to delete your post',error)
//         throw error
//     }
// }