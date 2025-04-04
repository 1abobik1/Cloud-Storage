'use client'
import React, { createContext, useContext, useState } from 'react';

interface ActiveTabContextType {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

const ActiveTabContext = createContext<ActiveTabContextType | undefined>(undefined);

interface ActiveTabProviderProps {
  children: React.ReactNode; // Указываем тип для children
}

export const ActiveTabProvider: React.FC<ActiveTabProviderProps> = ({ children }) => {
  const [activeTab, setActiveTab] = useState<string>('/cloud/home');

  return (
    <ActiveTabContext.Provider value={{ activeTab, setActiveTab }}>
      {children}
    </ActiveTabContext.Provider>
  );
};

export const useActiveTab = () => {
  const context = useContext(ActiveTabContext);
  if (!context) {
    throw new Error('useActiveTab must be used within an ActiveTabProvider');
  }
  return context;
};
