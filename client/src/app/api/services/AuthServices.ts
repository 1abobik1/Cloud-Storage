
import $api from '../http';
import { AxiosResponse} from 'axios';
import { AuthResponce } from '../models/response/AuthResponse';

export default class AuthServices{
    static async login(email:string,password:string):Promise<AxiosResponse<AuthResponce>>{
        return $api.post<AuthResponce>('/login/',{email,password})
        
    }
    static async registration (username:string,email:string,password:string):Promise<AxiosResponse<AuthResponce>>{
        return $api.post<AuthResponce>('/signup/',{username,email,password})
        
    }
    static async logout ():Promise<void>{
        return $api.post('/logout/')
        
    }



}
