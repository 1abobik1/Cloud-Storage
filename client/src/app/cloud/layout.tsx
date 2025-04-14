'use client';

import React, {useContext, useEffect} from 'react';
import {useRouter} from 'next/navigation';
import {Context} from '@/app/_app';
import SideBar from '@/app/ui/SideBar';
import Header from '../ui/Header';
import {observer} from 'mobx-react-lite';

function Layout({ children }: { children: React.ReactNode }) {
    const { store } = useContext(Context);
    const router = useRouter();

    useEffect(() => {
        if (!store.isLoading && !store.isAuth) {
            router.push('/'); // редиректим на главную/логин
        }
    }, [store.isAuth, store.isLoading]);

    if (store.isLoading) {
        return <div className="p-12">Загрузка...</div>;
    }

    if (!store.isAuth) {
        return null; // пока идёт редирект
    }

    return (
        <div className="flex h-screen h-[90vh] flex-col md:flex-row md:overflow-hidden">
            <div className="w-full flex-none md:w-64 bg-gray-100">
                <SideBar />
            </div>

            <div className="flex flex-col flex-grow w-100%">
                <Header />
                <div className="flex-grow p-6 md:overflow-y-auto md:p-12">{children}</div>
            </div>
        </div>
    );
}

export default observer(Layout);
