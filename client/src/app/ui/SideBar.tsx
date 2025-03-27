import Link from 'next/link';
 import NavLinks from './navlinks';
import AcmeLogo from '@/app/ui/acme-logo';

export default function SideBar() {
  return (
   <div className='flex  h-full flex-col py-6 '>
    <Link href="/" className='mb-2 flex h-20 items-end justify-start rounded-md bg-blue-700 p-4 md:h-40'> 
    
     <div className=' md:40'> <AcmeLogo/></div>
    </Link>
   <div className='flex grow flex-row justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2'>
   <NavLinks/>
     <div className='flex flex-grow bg-slate-600 rounded-md'></div>

    <Link href="/cloud/login" className='bg-slate-600 p-4 gap-2 rounded-md transition duration-300 ease-in-out transform hover:bg-slate-700 hover:scale-105'>Войти</Link>

   </div>
  </div>

  );
}
