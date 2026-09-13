'use client';

import { useEffect, useState } from 'react';
import {
  AUTH_CHANGE_EVENT,
  fetchAuthEnabled,
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
  /** True when the server requires Bearer auth (default until probe returns). */
  authRequired: boolean;
  /** True after /auth/status has been resolved (or failed closed). */
  authStatusLoaded: boolean;
  /** Review pages may load data when authenticated or when auth is disabled. */
  canAccessReview: boolean;
  /** Mutations are allowed for admins, or for anyone when auth is disabled. */
  canSubmitReview: boolean;
};

const ANONYMOUS_SESSION: Omit<AuthSessionState, 'authEpoch'> = {
  token: null,
  username: null,
  role: null,
  isAuthenticated: false,
  isAdmin: false,
  authRequired: true,
  authStatusLoaded: false,
  canAccessReview: false,
  canSubmitReview: false,
};

function deriveAccess(
  isAuthenticated: boolean,
  isAdmin: boolean,
  authRequired: boolean
): Pick<AuthSessionState, 'canAccessReview' | 'canSubmitReview'> {
  return {
    canAccessReview: !authRequired || isAuthenticated,
    canSubmitReview: !authRequired || isAdmin,
  };
}

function readSession(
  authRequired: boolean,
  authStatusLoaded: boolean
): Omit<AuthSessionState, 'authEpoch'> {
  const token = getAuthToken();
  const role = getAuthRole();
  const isAuthenticated = Boolean(token);
  const isAdmin = role === 'admin';
  return {
    token,
    username: getAuthUsername(),
    role,
    isAuthenticated,
    isAdmin,
    authRequired,
    authStatusLoaded,
    ...deriveAccess(isAuthenticated, isAdmin, authRequired),
  };
}

/** Subscribe to login/logout so protected pages can reload or clear state. */
export function useAuthSession(): AuthSessionState {
  // Always start anonymous so SSR markup matches the first client render.
  const [session, setSession] = useState<AuthSessionState>(() => ({
    authEpoch: 0,
    ...ANONYMOUS_SESSION,
  }));

  useEffect(() => {
    let cancelled = false;

    const sync = () => {
      setSession((prev) => ({
        authEpoch: prev.authEpoch + 1,
        ...readSession(prev.authRequired, prev.authStatusLoaded),
      }));
    };

    // Hydrate from localStorage only after mount.
    setSession((prev) => ({
      ...prev,
      ...readSession(prev.authRequired, prev.authStatusLoaded),
    }));

    (async () => {
      const enabled = await fetchAuthEnabled();
      if (cancelled) return;
      setSession((prev) => {
        const next = readSession(enabled, true);
        return {
          authEpoch: prev.authEpoch + 1,
          ...next,
        };
      });
    })();

    window.addEventListener(AUTH_CHANGE_EVENT, sync);
    window.addEventListener('storage', sync);
    return () => {
      cancelled = true;
      window.removeEventListener(AUTH_CHANGE_EVENT, sync);
      window.removeEventListener('storage', sync);
    };
  }, []);

  return session;
}
