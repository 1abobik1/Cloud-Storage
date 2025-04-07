import {IUser} from "../models/IUser";
import {makeAutoObservable} from "mobx";
import AuthService from "@/app/api/services/AuthServices";
import axios from 'axios';
import {AuthResponse} from "../models/response/AuthResponse";
import {API_URL} from "@/app/api/http/index";
import {jwtDecode} from 'jwt-decode';

interface JwtPayload {
    user_id: number;
    exp: number;
    is_superuser: boolean;
}

const isUserSuperuser = (): boolean => {
    const token = localStorage.getItem('token');
    if (token) {
        try {
            const decoded: JwtPayload = jwtDecode<JwtPayload>(token);
            return decoded.is_superuser;
        } catch (error) {
            console.error('Ошибка при декодировании токена:', error);
            return false;
        }
    }
    return false;
};


export default class Store {
    user = {} as IUser;
    code = 0;
    isAuth = false;
    isSuperUser = false;
    isLoading = false;

    constructor() {
        makeAutoObservable(this);
    }

    setAuth(bool: boolean) {
        this.isAuth = bool;
    }

    setUser(user: IUser) {
        this.user = user;
    }

    setCode(code: number) {
        this.code = code;
    }
    getCode() {
        return this.code;
    }
    setSuperUser(bool:boolean){
        this.isSuperUser = bool;
    }

    setLoading(bool: boolean) {
        this.isLoading = bool;
    }

    async login(email: string, password: string,platform: string) {
        try {
            const response = await AuthService.login(email, password,platform);
            console.log(response)
            localStorage.setItem('token', response.data.access);
            this.setAuth(true);
            this.setUser(response.data.user);
        } catch (e) {
            // @ts-ignore
            console.log(e.response?.data?.message);
        }
    }

    async verify(email: string) {
        try {
            const response = await AuthService.verify(email);
            console.log(response)
            this.setCode(response.data.code);
            return Promise.resolve();
        } catch (e) {
            // @ts-ignore
            return Promise.reject(e);
        }
    }

    async signup(username: string, email: string, password: string) {
        try {
            const response = await AuthService.signup(username, email, password);
            console.log(response)
            localStorage.setItem('token', response.data.access);
            this.setAuth(true);
            this.setUser(response.data.user);
        } catch (e) {
            // @ts-ignore
            console.log(e.response?.data?.message);
        
        if (e.response?.status === 409) {
           
            alert('Аккаунт на эту почту уже зарегистрированы.');
           
          } else {
           
            console.error(e);
            alert('Произошла ошибка при регистрации');
          }
        }
    }

    async logout() {
        try {
            const response = await AuthService.logout();
            console.log(response)
            localStorage.removeItem('token');
            this.setAuth(false);
            this.setUser({} as IUser);
        } catch (e: any) {
            console.log(e.response?.data?.message);
        }
    }

    async checkAuth() {
        this.setLoading(true);
        try {
            const response = await axios.post<AuthResponse>(
                `${API_URL}/token/update`,
                {},
                {
                    withCredentials: true
                }
            );
            console.log(response)
            localStorage.setItem('token', response.data.access);
            this.setAuth(true);
            
            this.setUser(response.data.user);
            return Promise.resolve();
        } catch (e: any) {
            if (e.response?.status === 401) {
                this.setAuth(false);
                localStorage.removeItem('token');
                console.log('Не удалось обновить токен: пользователь не авторизован.');
            } else {
                console.log('Ошибка авторизации:', e.response?.data?.message);
            }
            return Promise.reject(e);
        } finally {
            this.setLoading(false);
        }
    }

}