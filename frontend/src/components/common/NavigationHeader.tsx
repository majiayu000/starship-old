'use client';

import React, { FormEvent, useEffect, useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import {
  getAuthToken,
  getAuthUsername,
  login,
  logout,
} from '@/api/auth';

const NavigationHeader: React.FC = () => {
  const pathname = usePathname();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [username, setUsername] = useState<string | null>(null);
  const [hasToken, setHasToken] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    setUsername(getAuthUsername());
    setHasToken(Boolean(getAuthToken()));
  }, []);

  const handleLogin = async (event: FormEvent) => {
    event.preventDefault();
    setIsSubmitting(true);
    try {
      const user = await login(email, password);
      const displayName =
        [user.firstName, user.lastName].filter(Boolean).join(' ') || user.email;
      setUsername(displayName);
      setHasToken(true);
      setPassword('');
      toast.success('登录成功');
    } catch (error) {
      console.error('登录失败:', error);
      toast.error('登录失败，请检查邮箱和密码');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleLogout = () => {
    logout();
    setUsername(null);
    setHasToken(false);
    toast.success('已退出登录');
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-white">
      <div className="container mx-auto px-4 py-3 flex items-center justify-between gap-4">
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

        <div className="flex items-center gap-3 text-sm text-gray-500">
          {hasToken ? (
            <>
              <span>{username || '已登录'}</span>
              <Button variant="outline" size="sm" onClick={handleLogout}>
                退出
              </Button>
            </>
          ) : (
            <form className="flex items-center gap-2" onSubmit={handleLogin}>
              <input
                type="email"
                required
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                placeholder="邮箱"
                className="h-8 rounded-md border px-2 text-sm text-gray-900"
                aria-label="登录邮箱"
              />
              <input
                type="password"
                required
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                placeholder="密码"
                className="h-8 rounded-md border px-2 text-sm text-gray-900"
                aria-label="登录密码"
              />
              <Button type="submit" size="sm" disabled={isSubmitting}>
                {isSubmitting ? '登录中…' : '登录'}
              </Button>
            </form>
          )}
          <span className="hidden lg:inline">
            API: {process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'}
          </span>
        </div>
      </div>
    </header>
  );
};

export default NavigationHeader;
