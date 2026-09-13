'use client';

import {
  createContext,
  createElement,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from 'react';
import {
  AUTH_CHANGE_EVENT,
  fetchAuthStatus,
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
  // Wait for /auth/status before granting access so hydration does not start
  // protected loaders that the probe would immediately invalidate/duplicate.
  const access = authStatusLoaded
    ? deriveAccess(isAuthenticated, isAdmin, authRequired)
    : { canAccessReview: false, canSubmitReview: false };
  return {
    token,
    username: getAuthUsername(),
    role,
    isAuthenticated,
    isAdmin,
    authRequired,
    authStatusLoaded,
    ...access,
  };
}

function sessionAccessFingerprint(
  s: Omit<AuthSessionState, 'authEpoch'>
): string {
  return [
    s.token ?? '',
    s.role ?? '',
    s.isAuthenticated,
    s.isAdmin,
    s.authRequired,
    s.authStatusLoaded,
    s.canAccessReview,
    s.canSubmitReview,
  ].join('|');
}

const AuthSessionContext = createContext<AuthSessionState | null>(null);

/** Single shared auth session so pages/forms do not each probe /auth/status. */
function useAuthSessionState(): AuthSessionState {
  // Always start anonymous so SSR markup matches the first client render.
  const [session, setSession] = useState<AuthSessionState>(() => ({
    authEpoch: 0,
    ...ANONYMOUS_SESSION,
  }));

  useEffect(() => {
    let cancelled = false;

    const applySession = (
      nextBase: Omit<AuthSessionState, 'authEpoch'>,
      forceEpochBump: boolean
    ) => {
      setSession((prev) => {
        const changed =
          forceEpochBump ||
          sessionAccessFingerprint(prev) !== sessionAccessFingerprint(nextBase);
        return {
          authEpoch: changed ? prev.authEpoch + 1 : prev.authEpoch,
          ...nextBase,
        };
      });
    };

    const syncFromStorage = () => {
      setSession((prev) => {
        const next = readSession(prev.authRequired, prev.authStatusLoaded);
        const changed =
          sessionAccessFingerprint(prev) !== sessionAccessFingerprint(next);
        return {
          authEpoch: changed ? prev.authEpoch + 1 : prev.authEpoch,
          ...next,
        };
      });
    };

    // Hydrate token/role from localStorage after mount, but keep access gated
    // until /auth/status completes (authStatusLoaded stays false here).
    setSession((prev) => ({
      ...prev,
      ...readSession(prev.authRequired, prev.authStatusLoaded),
    }));

    (async () => {
      const status = await fetchAuthStatus();
      if (cancelled) return;
      // fetchAuthStatus may have updated AUTH_ROLE_KEY; re-read after probe.
      applySession(readSession(status.enabled, true), false);
    })();

    const onFocus = () => {
      // Re-probe on focus so a promoted user picks up the new role without logout.
      void (async () => {
        const status = await fetchAuthStatus();
        if (cancelled) return;
        applySession(readSession(status.enabled, true), false);
      })();
    };

    window.addEventListener(AUTH_CHANGE_EVENT, syncFromStorage);
    window.addEventListener('storage', syncFromStorage);
    window.addEventListener('focus', onFocus);
    return () => {
      cancelled = true;
      window.removeEventListener(AUTH_CHANGE_EVENT, syncFromStorage);
      window.removeEventListener('storage', syncFromStorage);
      window.removeEventListener('focus', onFocus);
    };
  }, []);

  return session;
}

export function AuthSessionProvider({ children }: { children: ReactNode }) {
  const session = useAuthSessionState();
  return createElement(AuthSessionContext.Provider, { value: session }, children);
}

/** Subscribe to the shared auth session from AuthSessionProvider. */
export function useAuthSession(): AuthSessionState {
  const shared = useContext(AuthSessionContext);
  if (!shared) {
    throw new Error('useAuthSession must be used within AuthSessionProvider');
  }
  return shared;
}
