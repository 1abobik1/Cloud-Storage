'use client'
import type { AppProps } from 'next/app'
import React, { createContext } from 'react';
import Store from "@/app/api/store/store";

interface State {
  store: Store;
}
const store = new Store();
export const Context = createContext<State>({
  store, 
 })

export default function MyApp({ Component, pageProps }: AppProps) {
  
<Context.Provider value={{store}}>
  return <Component {...pageProps} />
</Context.Provider>
}