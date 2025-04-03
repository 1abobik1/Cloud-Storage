import $api from "@/app/api/http/index";
import {AxiosResponse} from 'axios';
import {AuthResponse} from "@/app/api/models/response/AuthResponse";


export default class AuthService {
    static async login(email: string, password: string,platform:string): Promise<AxiosResponse<AuthResponse>> {
        return $api.post<AuthResponse>('/api/user/login', {email, password,platform})
    }

    static async signup(email: string, password: string,platform:string): Promise<AxiosResponse<AuthResponse>> {
        return $api.post<AuthResponse>('/api/user/signup', {email, password,platform})
    }

    static async logout(): Promise<void> {
        return $api.post('/api/user/logout')
    }
    static async verify(email: string) {
        // @ts-ignore
        return authApi.post<CodeResponse>('/verify-email/', {email})
    }

}

