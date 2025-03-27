import { makeAutoObservable } from "mobx";
import { IUser } from "../models/IUser";
import AuthServices from "../services/AuthServices";



export default class Store{
    user = {} as IUser
    isAuth = false


constructor(){
    makeAutoObservable(this)
}

setAuth(bool: boolean){
    this.isAuth = bool
}

setUser(user: IUser){
    this.user = user
}       

async login(email: string, password: string){
    try{
        const response = await AuthServices.login(email, password) 
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
    
    }catch(e){
        console.log(e.response?.data?.message)
    }

}

async registration(username: string,email: string, password: string){
    try{
        const response = await AuthServices.registration(username,email, password) 
        localStorage.setItem('token', response.data.accessToken)
        console.log(response)
        this.setAuth(true)
        this.setUser(response.data.user)
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
    }catch(e){
        console.log(e.response?.data?.message)
    }
}


    async logout(){
        try{
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            const response = await AuthServices.logout() 
            localStorage.removeItemItem('token')
            this.setAuth(false)
            this.setUser({} as IUser)
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
        }catch(e){
            console.log(e.response?.data?.message)
        }

}



}