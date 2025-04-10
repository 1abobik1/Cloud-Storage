'use client'
import {useContext, useEffect} from "react";
import {Context} from "@/app/_app";
import LoginForm from "./api/components/LoginForm";
import {observer} from "mobx-react-lite";
import {useRouter} from 'next/navigation';

function Home() {
  const { store } = useContext(Context);
  const router = useRouter(); // ✅ выносим вверх, чтобы вызывался всегда

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
