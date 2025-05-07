'use client'
import {useContext, useEffect, useState} from "react";
import {Context} from "@/app/_app";
import LoginForm from "./api/components/LoginForm";
import {observer} from "mobx-react-lite";
import {useRouter} from 'next/navigation';
import {Loader2} from 'lucide-react';

function Home() {
    const { store } = useContext(Context);
    const router = useRouter();
    const [initialCheckDone, setInitialCheckDone] = useState(false);

    useEffect(() => {
        const checkAuth = async () => {
            if (localStorage.getItem('token')) {
                await store.checkAuth();
            }
            setInitialCheckDone(true);
        };

        checkAuth();
    }, []);

    useEffect(() => {
        if (initialCheckDone && store.isAuth) {
            router.push('/cloud');
        }
    }, [store.isAuth, initialCheckDone]);

    if (!initialCheckDone || store.isLoading) {
        return (
            <div className="flex items-center justify-center h-screen">
                <Loader2 className="w-10 h-10 text-blue-500 animate-spin" />
            </div>
        );
    }

    if (store.isAuth) {
        return <div>Загрузка...</div>
    }

    return <LoginForm />;
}

export default observer(Home);