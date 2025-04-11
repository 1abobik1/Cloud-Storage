import {AxiosResponse} from 'axios';
import {CloudResponse} from "@/app/api/models/response/CloudResponse";
import {FileData} from "@/app/api/models/FileData";
import {cloudApi} from '@/app/api/http/cloud'

export default class CloudService{
    static async getAllCloud(type:string):Promise<AxiosResponse<CloudResponse>> {
        // @ts-ignore
        return cloudApi.get<CloudResponse>(`/files/all?type=${type}`);
    }

    static async createCloud(file: string,) {
        return cloudApi.post(`/files/one`, {file});
    }
   
}