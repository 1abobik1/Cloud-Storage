import Link from 'next/link';
 import NavLinks from './navlinks';
import AcmeLogo from '@/app/ui/acme-logo';

export default function SideBar() {
  return (
   <div className='flex  h-full flex-col py-6 border-r border-gray-200'>
   
   <div className=' flex-row justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2 mx-2'>
   <NavLinks/>
   

    

   </div>
  </div>

  );
}
