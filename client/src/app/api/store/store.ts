
import {makeAutoObservable} from "mobx";
import AuthService from "@/app/api/services/AuthServices";
import axios from 'axios';
import {AuthResponse} from "../models/response/AuthResponse";
import {AUTH_API_URL} from "@/app/api/http/urls";



export default class Store {
   
    code = 0;
    isAuth = false;
   
    isLoading = false;

    constructor() {
        makeAutoObservable(this);
    }

    setAuth(bool: boolean) {
        this.isAuth = bool;
    }

    

    setCode(code: number) {
        this.code = code;
    }
    getCode() {
        return this.code;
    }
    

    setLoading(bool: boolean) {
        this.isLoading = bool;
    }

    async login(email: string, password: string,platform: string) {
        try {
            const response = await AuthService.login(email, password,platform);
            console.log(response)
            localStorage.setItem('token', response.data.accessToken);
            this.setAuth(true);
            
        } catch (e) {
            // @ts-ignore
            console.log(e.response?.data?.message);
        }
    }

    // async verify(email: string) {
    //     try {
    //         const response = await AuthService.verify(email);
    //         console.log(response)
    //         this.setCode(response.data.code);
    //         return Promise.resolve();
    //     } catch (e) {
    //         // @ts-ignore
    //         return Promise.reject(e);
    //     }
    // }

    async signup(username: string, email: string, password: string) {
        try {
            const response = await AuthService.signup(username, email, password);
            console.log(response)
            localStorage.setItem('token', response.data.accessToken);
            this.setAuth(true);
            
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
            
        } catch (e: any) {
            console.log(e.response?.data?.message);
        }
    }

    async checkAuth() {
        this.setLoading(true);
        try {
            const response = await axios.post<AuthResponse>(
                `${AUTH_API_URL}/token/update/`,
                {},
                {
                    withCredentials: true
                }
            );
            console.log(response)
            localStorage.setItem('token', response.data.accessToken);
            this.setAuth(true);
            
          
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