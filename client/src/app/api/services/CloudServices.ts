import {AxiosResponse} from 'axios';
import {CloudResponse} from "@/app/api/models/response/CloudResponse";
import $api from "@/app/api/http/index";

export default class CloudService{
    static async getAllCloud() {
        // @ts-ignore
        return $api.get<CloudResponse[]>('/files/all');
    }

    static async createCloud(file: string,) {
        return $api.post(`/files/one`, {file});
    }
   
}