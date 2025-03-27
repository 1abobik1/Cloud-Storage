'use client';

import {
  UserGroupIcon,
  HomeIcon,
  DocumentDuplicateIcon,
} from '@heroicons/react/24/outline';

import clsx from 'clsx';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const links = [
  { name: 'Home', href: '/cloud', icon: HomeIcon },
  {name: 'Clouds',href: '/cloud/posts',icon: DocumentDuplicateIcon,},
  { name: 'Settings', href: '/cloud/settings', icon: UserGroupIcon },
  { name: 'Your posts', href: '/cloud/', icon: UserGroupIcon },
  
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
            'flex h-[48px] grow items-center justify-center gap-2 rounded-md bg-slate-600 p-3 text-sm font-medium hover:bg-slate-300 hover:text-white md:flex-none md:justify-start md:p-2 md:px-3',
            {
              'bg-slate-900 text-white': pathname === link.href,
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
