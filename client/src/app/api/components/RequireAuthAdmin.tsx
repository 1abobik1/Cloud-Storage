'use client';

import React, { FC, useContext, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { observer } from 'mobx-react-lite';
import { Context } from "@/app/_app";
interface Props {
  children: React.ReactNode;
}

const RequireAuthAdmin: FC<Props> = ({ children }) => {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState<boolean | null>(null);
  const { store } = useContext(Context);  

  useEffect(() => {
    const checkAuthentication = async () => {
      try {
        if (typeof window !== 'undefined') {
          if (localStorage.getItem('token') && !store.isAuth) {
            await store.checkAuth();
          }

          if (store.isAuth ) {
            setIsAuthorized(true);
          } else {
            setIsAuthorized(false);
          }
        }
      } catch (error) {
        setIsAuthorized(false);
      } finally {
        setLoading(false);
      }
    };

    checkAuthentication();
  }, []);

  useEffect(() => {
    if (!loading && isAuthorized === false) {
      router.push('/login');
    }
  }, [loading, isAuthorized, router]);

  if (loading || isAuthorized === false) {
    return <div className="loading-spinner"></div>; // можно заменить на что-то красивее
  }

  return <>{children}</>;
};

export default observer(RequireAuthAdmin);
