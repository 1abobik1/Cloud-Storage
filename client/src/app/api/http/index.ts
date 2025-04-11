import axios from 'axios';
import {AuthResponse} from "@/app/api/models/response/AuthResponse";
import {AUTH_API_URL, CLOUDAPI_URL} from "./urls";

export const createApiInstance = (
    baseURL: string,
    
): any => {
    const api = axios.create({
        withCredentials: true,
        baseURL: baseURL
    });

    api.interceptors.request.use((config) => {
        const token = localStorage.getItem('token');

        if (baseURL === CLOUDAPI_URL) {
            config.headers.Authorization = `Bearer ${token}`;
        } 
        return config; 

    }, (error) => {
        return Promise.reject(error);
    });
    api.interceptors.response.use((config) => {
        return config;
    }, async (error) => {
        const originalRequest = error.config;

        if (error.response.status === 401 && !originalRequest._isRetry) {
            originalRequest._isRetry = true;

            try {
                const response = await axios.post<AuthResponse>(`${AUTH_API_URL}/token/update`, {}, {
                    withCredentials: true
                });
                localStorage.setItem('token', response.data.accessToken);
                originalRequest.headers.Authorization = `Bearer ${response.data.accessToken}`;
                return api.request(originalRequest);
            } catch (e) {
                console.log('Ошибка авторизации при обновлении токена');
            }
        }
        return Promise.reject(error);
    });
    return api;
};
