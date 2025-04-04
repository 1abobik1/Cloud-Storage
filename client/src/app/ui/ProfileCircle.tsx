'use client';

import React, { useState } from 'react';

const ProfileCircle = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const handleMenuToggle = () => {
    setIsMenuOpen(!isMenuOpen);
  };

  const handleMenuClose = () => {
    setIsMenuOpen(false);
  };

  return (
    <div className="relative">
      {/* Здесб будет картинка профиля */}
      <div
        onClick={handleMenuToggle}
        className="w-10 h-10 rounded-full overflow-hidden cursor-pointer border-2 border-gray-300"
      >
        <img
          src="https://via.placeholder.com/150"
          alt="Profile"
          className="w-full h-full object-cover"
        />
      </div>

      {/* Контекстное меню */}
      {isMenuOpen && (
        <div
          className="absolute right-0 mt-2 w-48 bg-white shadow-lg rounded-lg border border-gray-200"
          onClick={handleMenuClose}
        >
          <ul>
            <li className="px-4 py-2 text-gray-700 hover:bg-gray-100 cursor-pointer">Profile</li>
            <li className="px-4 py-2 text-gray-700 hover:bg-gray-100 cursor-pointer">Settings</li>
            <li className="px-4 py-2 text-gray-700 hover:bg-gray-100 cursor-pointer">Logout</li>
          </ul>
        </div>
      )}
    </div>
  );
};

export default ProfileCircle;
