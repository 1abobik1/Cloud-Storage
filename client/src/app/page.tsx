'use client'

import { useContext, useState, useEffect } from "react";
import UserService from "@/app/api/services/UserServices";
import { Context } from "@/app/_app";
import LoginForm from "./api/components/LoginForm";
import { IUser } from '@/app/api/models/IUser'
import { observer } from "mobx-react-lite";
import { useRouter } from 'next/navigation';

function Home() {
  const { store } = useContext(Context);
  const router = useRouter(); 


  useEffect(() => {
    if (localStorage.getItem('token')) {
      store.checkAuth();
    }
  }, []);

  useEffect(() => {
    if (store.isAuth) {
      router.push('/cloud/home'); // переход после авторизации
    }
  }, [store.isAuth]);

 

  if (store.isLoading) {
    return <div>Загрузка...</div>
  }

  if (!store.isAuth) {
    return (
      <div>
        <LoginForm />
       
      </div>
    );
  }

  return (
    <div>
      <button onClick={() => store.logout()}>Выйти</button>
    </div>
  );
}

export default observer(Home);
