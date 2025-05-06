'use client'
import {useContext, useEffect} from "react";
import {Context} from "@/app/_app";
import LoginForm from "./api/components/LoginForm";
import {observer} from "mobx-react-lite";
import {useRouter} from 'next/navigation';
import {AppRouterInstance} from "next/dist/shared/lib/app-router-context.shared-runtime";

function Home() {
    const {store} = useContext(Context);
    const router: AppRouterInstance = useRouter();

    useEffect(() => {
        if (localStorage.getItem('token')) {
            store.checkAuth();
        }
    }, []);

    useEffect(() => {
        if (store.isAuth) {
            router.push('/cloud'); // переход после авторизации
        }
    }, [store.isAuth]);


    if (store.isLoading) {
        return <div>Загрузка...</div>
    }


    return (
        <div>
            <LoginForm/>
        </div>
   );

}

export default observer(Home);
