'use client';

import { useEffect, useState } from 'react';
import {
  AUTH_CHANGE_EVENT,
  getAuthRole,
  getAuthToken,
  getAuthUsername,
} from '@/api/auth';

export type AuthSessionState = {
  authEpoch: number;
  token: string | null;
  username: string | null;
  role: string | null;
  isAuthenticated: boolean;
  isAdmin: boolean;
};

function readSession(): Omit<AuthSessionState, 'authEpoch'> {
  const token = getAuthToken();
  const role = getAuthRole();
  return {
    token,
    username: getAuthUsername(),
    role,
    isAuthenticated: Boolean(token),
    isAdmin: role === 'admin',
  };
}

/** Subscribe to login/logout so protected pages can reload or clear state. */
export function useAuthSession(): AuthSessionState {
  const [session, setSession] = useState<AuthSessionState>(() => ({
    authEpoch: 0,
    ...readSession(),
  }));

  useEffect(() => {
    const sync = () => {
      setSession((prev) => ({
        authEpoch: prev.authEpoch + 1,
        ...readSession(),
      }));
    };

    // Hydrate from localStorage after mount (SSR-safe initial nulls).
    setSession((prev) => ({
      ...prev,
      ...readSession(),
    }));

    window.addEventListener(AUTH_CHANGE_EVENT, sync);
    window.addEventListener('storage', sync);
    return () => {
      window.removeEventListener(AUTH_CHANGE_EVENT, sync);
      window.removeEventListener('storage', sync);
    };
  }, []);

  return session;
}
