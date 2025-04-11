'use client';

import { useState, useEffect } from 'react';
import NavLinks from './navlinks';
import { MenuIcon, XIcon } from 'lucide-react';
import { usePathname } from 'next/navigation';
import clsx from 'clsx';
export default function SideBar() {
  const [isOpen, setIsOpen] = useState(false);
  const pathname = usePathname();

  // Закрывать бургер-меню при изменении пути
  useEffect(() => {
    setIsOpen(false);
  }, [pathname]);

  return (
    <>
      {/* Бургер кнопка - видна только на маленьких экранах */}
      <div className="md:hidden flex justify-between items-center px-4 py-2 border-b border-gray-200">
        <button onClick={() => setIsOpen(!isOpen)} className="text-gray-600">
          {isOpen ? <XIcon className="w-6 h-6" /> : <MenuIcon className="w-6 h-6" />}
        </button>
        <span className="text-lg font-semibold">Menu</span>
      </div>

      {/* Боковая панель */}
      <div className={clsx(
        'flex-col border-r border-gray-200 bg-white w-60 h-full px-4 py-6 md:flex transition-transform duration-300',
        {
          'hidden md:flex': !isOpen,
          'absolute z-50 left-0 top-0 h-screen shadow-lg md:relative': isOpen,
        }
      )}>
        <NavLinks />
      </div>
    </>
  );
}
