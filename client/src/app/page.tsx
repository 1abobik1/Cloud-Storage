'use client'
import Link from "next/link";
import { useContext, useState, useEffect } from "react";
import UserService from "@/app/api/services/UserServices";
import { Context } from "@/app/_app";
import LoginForm from "./api/components/LoginForm";
import { IUser } from '@/app/api/models/IUser'
import { observer } from "mobx-react-lite";
import { useRouter } from 'next/navigation';

function Home() {
  const { store } = useContext(Context);
  const router = useRouter(); // ✅ выносим вверх, чтобы вызывался всегда
  const [users, setUsers] = useState<IUser[]>([]);

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

  async function getUsers() {
    try {
      const response = await UserService.fetchUsers();
      setUsers(response.data);
    } catch (e) {
      console.log(e);
    }
  }

  if (store.isLoading) {
    return <div>Загрузка...</div>
  }

  if (!store.isAuth) {
    return (
      <div>
        <LoginForm />
        <button onClick={getUsers}>Получить пользователей</button>
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
