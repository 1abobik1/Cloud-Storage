'use client';
import React, { useState } from 'react';
import { usePathname } from 'next/navigation';
import {
    HomeIcon,
    Squares2X2Icon,
    BellIcon,
    MagnifyingGlassIcon,
} from '@heroicons/react/24/solid';
import {
    ArrowUpOnSquareIcon
} from '@heroicons/react/24/outline';
import ProfileCircle from './ProfileCircle';
import FileUploader from '../components/FileUploader';
const links = [
    { name: 'Home', href: '/cloud/home', icon: HomeIcon },
    { name: 'Workspaces', href: '/cloud/workspaces', icon: Squares2X2Icon },
    { name: 'Search', href: '/cloud/search', icon: MagnifyingGlassIcon },
    { name: 'Notifications', href: '/cloud/notifications', icon: BellIcon },
];
const Header = () => {
    const pathname = usePathname();
    const activeLink = links.find((link) => pathname.startsWith(link.href));
    if (!activeLink) return null; // Если путь не найден, не рендерим ничего
    const Icon = activeLink.icon;
    const LinkIcon = ArrowUpOnSquareIcon

    return (
        <div className='flex flex-row justify-between px-5 '>
            <div className="flex flex-row  p-4 my-4 bg-white ">
                <Icon className="w-6 text-blue-600" />
                <span className="ml-2 mt-1 text-sm font-semibold text-gray-900">{activeLink.name}</span>
            </div>
            <div className='flex flex-row py-6 mx-5'>
            <FileUploader/>
            <div className='mx-5'><ProfileCircle /></div>
            </div>
        </div>
    );
};

export default Header;
