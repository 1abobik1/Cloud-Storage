import {AxiosResponse} from 'axios';
import {CloudResponse} from "@/app/api/models/response/CloudResponse";
import {cloudApi} from '@/app/api/http/cloud';
import {cryptoHelper} from '@/app/api/utils/CryptoHelper';
import {FileData} from "@/app/api/models/FileData";

export default class CloudService {
    static async getAllCloud(type: string): Promise<AxiosResponse<CloudResponse>> {
        const response = await cloudApi.get<CloudResponse>(`/files/all?type=${type}`);

        const fileData = response.data.file_data;

        if (Array.isArray(fileData) && fileData.length > 0) {
            response.data.file_data = await Promise.all(
                fileData.map(async (file: FileData) => {
                    try {
                        const res = await fetch(file.url);
                        const blob = await res.blob();
                        const decryptedFile = await cryptoHelper.decryptFile(new File([blob], file.name));
                        return {
                            ...file,
                            decryptedFile,
                        };
                    } catch (error) {
                        console.error(`Ошибка при расшифровке файла ${file.name}:`, error);
                        return {
                            ...file,
                            decryptedFile: null,
                        };
                    }
                })
            );
        }

        return response;
    }


    static async uploadFiles(formData: FormData, config = {}) {
        return await cloudApi.post(`/files/many`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
            ...config, // <- прокидываем onUploadProgress и другие настройки
        });
    }


    static async deleteFile(type: string, obj_id: string) {
        return await cloudApi.delete(`files/one?id=${obj_id}&type=${type}`);
    }
}