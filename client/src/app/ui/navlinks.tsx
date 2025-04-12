'use client';

import {
  HomeIcon,
  Squares2X2Icon,
  BellIcon,
  MagnifyingGlassIcon
} from '@heroicons/react/24/outline';


import clsx from 'clsx';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const links = [
  { name: 'Home', href: '/cloud/home', icon: HomeIcon },
  {name: 'Workspaces',href: '/cloud/workspaces',icon: Squares2X2Icon,},
  
  
];

export default function NavLinks() {
  const pathname = usePathname();
 
 
  return (
    <>
      {links.map((link) => {
        const LinkIcon = link.icon;
        return (
          <Link
          key={link.name}
          href={link.href}
          
          className={clsx(
            'flex h-[36px]  grow items-center justify-center text-gray-500 gap-2 rounded-md text-sm font-medium  hover:text-blue-500 md:flex-none md:justify-start md:p-2 md:px-3',
            {
              'bg-blue-300 text-blue-600': pathname === link.href,
            },
          )}
        >
            <LinkIcon className="w-6" />
            <p className="hidden md:block">{link.name}</p>
          </Link>
        );
      })}
    </>
  );
}
