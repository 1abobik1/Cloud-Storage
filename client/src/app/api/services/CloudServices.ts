import {AxiosResponse} from 'axios';
import {CloudResponse} from "@/app/api/models/response/CloudResponse";
import {cloudApi} from '@/app/api/http/cloud'

export default class CloudService{
    static async getAllCloud(type:string):Promise<AxiosResponse<CloudResponse>> {
        // @ts-ignore
        return cloudApi.get<CloudResponse>(`/files/all?type=${type}`);
    }

    static async createCloud(formData: FormData) {
        return cloudApi.post(`/files/one`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });
      }

}