'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Button } from '@/components/ui/button';

const NavigationHeader: React.FC = () => {
  const pathname = usePathname();
  
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-white">
      <div className="container mx-auto px-4 py-3 flex items-center justify-between">
        <div className="flex items-center">
          <h1 className="text-xl font-bold mr-8">SAT题目审核系统</h1>
          <nav className="flex space-x-2">
            <Link href="/review">
              <Button
                variant={pathname === '/review' ? 'default' : 'ghost'}
                className="text-sm font-medium"
              >
                审核页面
              </Button>
            </Link>
            <Link href="/questions">
              <Button
                variant={pathname === '/questions' ? 'default' : 'ghost'}
                className="text-sm font-medium"
              >
                题目审核
              </Button>
            </Link>
          </nav>
        </div>
        
        <div className="text-sm text-gray-500">
          <span>API: {process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}</span>
        </div>
      </div>
    </header>
  );
};

export default NavigationHeader; 