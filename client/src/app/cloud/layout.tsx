'use client'

import SideBar from '@/app/ui/SideBar';
import React from 'react';
import Header from '../ui/Header';






export default function Layout({ children }: { children: React.ReactNode }) {
  
  return (
    
    <div className="flex h-screen h-[90vh] flex-col md:flex-row md:overflow-hidden ">
       
      <div className="w-full flex-none md:w-64 bg-gray-100">
        <SideBar />
      </div>
      
      <div className='flex flex-col flex-grow w-100%'>
        <Header/>
        
      
      <div className="flex-grow p-6 md:overflow-y-auto md:p-12">{children}</div>
      </div>
    </div>
    
  );
}




