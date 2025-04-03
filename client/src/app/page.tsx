'use client'
import Link from "next/link";
import { useContext } from "react";

import { Context } from "@/app/_app";
import { useEffect } from "react";

export default function Home() {
  const {store} = useContext(Context);

  useEffect(() => {
      if (localStorage.getItem('token')) {
          store.checkAuth();
      }
  }, []);

  if (store.isLoading) {
      return <div className="loading-spinner"></div>;
  }
  return (
    <div>
      <Link href="/cloud"> Start cloud doing!!</Link>
    </div>
   
  );
}
