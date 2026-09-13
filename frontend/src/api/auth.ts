import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const AUTH_TOKEN_KEY = 'auth_token';
export const AUTH_USERNAME_KEY = 'username';
export const AUTH_ROLE_KEY = 'auth_role';
export const AUTH_CHANGE_EVENT = 'auth-session-changed';

export type AuthUser = {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
};

export type AuthStatus = {
  enabled: boolean;
  /** Present when a valid Bearer token was accepted; reflects the current DB role. */
  role?: string | null;
};

type AuthResponse = {
  token: string;
  user: AuthUser;
};

function unwrapAuthPayload(data: unknown): AuthResponse {
  if (data && typeof data === 'object' && 'data' in data) {
    const nested = (data as { data: unknown }).data;
    if (nested && typeof nested === 'object' && 'token' in nested) {
      return nested as AuthResponse;
    }
  }
  return data as AuthResponse;
}

function notifyAuthChange(): void {
  if (typeof window === 'undefined') {
    return;
  }
  window.dispatchEvent(new Event(AUTH_CHANGE_EVENT));
}

export function getAuthToken(): string | null {
  if (typeof window === 'undefined') {
    return null;
  }
  return localStorage.getItem(AUTH_TOKEN_KEY);
}

export function getAuthUsername(): string | null {
  if (typeof window === 'undefined') {
    return null;
  }
  return localStorage.getItem(AUTH_USERNAME_KEY);
}

export function getAuthRole(): string | null {
  if (typeof window === 'undefined') {
    return null;
  }
  return localStorage.getItem(AUTH_ROLE_KEY);
}

export function isAdmin(): boolean {
  return getAuthRole() === 'admin';
}

/** Update only the cached role when the server reports a fresher value. */
export function updateCachedAuthRole(role: string, options?: { notify?: boolean }): void {
  if (typeof window === 'undefined') {
    return;
  }
  const next = role || 'user';
  if (localStorage.getItem(AUTH_ROLE_KEY) === next) {
    return;
  }
  localStorage.setItem(AUTH_ROLE_KEY, next);
  if (options?.notify === false) {
    return;
  }
  notifyAuthChange();
}

export function setAuthSession(token: string, user: AuthUser): void {
  localStorage.setItem(AUTH_TOKEN_KEY, token);
  const displayName =
    [user.firstName, user.lastName].filter(Boolean).join(' ') || user.email;
  localStorage.setItem(AUTH_USERNAME_KEY, displayName);
  localStorage.setItem(AUTH_ROLE_KEY, user.role || 'user');
  notifyAuthChange();
}

export function clearAuthSession(): void {
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem(AUTH_USERNAME_KEY);
  localStorage.removeItem(AUTH_ROLE_KEY);
  notifyAuthChange();
}

/**
 * Clear the session only when the stored token is still the one that failed.
 * Prevents an in-flight 401 for token A from wiping a newer session for token B.
 */
export function clearAuthSessionIfCurrent(requestToken: string | null | undefined): void {
  if (!requestToken) {
    return;
  }
  if (getAuthToken() === requestToken) {
    clearAuthSession();
  }
}

/** Headers for authenticated API calls. Omits Authorization when no token is stored. */
export function authHeaders(): Record<string, string> {
  const token = getAuthToken();
  if (!token) {
    return {};
  }
  return { Authorization: `Bearer ${token}` };
}

export async function login(email: string, password: string): Promise<AuthUser> {
  const response = await axios.post(`${API_BASE_URL}/auth/login`, {
    email,
    password,
  });
  const payload = unwrapAuthPayload(response.data);
  if (!payload?.token) {
    throw new Error('Login response missing token');
  }
  setAuthSession(payload.token, payload.user);
  return payload.user;
}

export function logout(): void {
  clearAuthSession();
}

function parseAuthStatusPayload(payload: unknown): AuthStatus | null {
  const body =
    payload && typeof payload === 'object' && 'data' in payload
      ? (payload as { data: unknown }).data
      : payload;
  if (!body || typeof body !== 'object') {
    return null;
  }
  const record = body as { enabled?: boolean; role?: string };
  if (typeof record.enabled !== 'boolean') {
    return null;
  }
  return {
    enabled: record.enabled,
    role: typeof record.role === 'string' ? record.role : null,
  };
}

/**
 * Probe whether the backend requires Bearer auth, and refresh the cached role
 * when a stored token is still valid (covers admin promotions without re-login).
 */
export async function fetchAuthStatus(): Promise<AuthStatus> {
  try {
    const response = await axios.get(`${API_BASE_URL}/auth/status`, {
      headers: authHeaders(),
    });
    const status = parseAuthStatusPayload(response.data);
    if (status) {
      if (status.role) {
        // Silent write: callers re-read localStorage and decide whether to bump epoch.
        updateCachedAuthRole(status.role, { notify: false });
      }
      return status;
    }
  } catch (error) {
    console.error('Failed to fetch auth status:', error);
  }
  // Fail closed: assume auth is required if the probe fails.
  return { enabled: true, role: null };
}

/** @deprecated Prefer fetchAuthStatus; kept for call-site compatibility. */
export async function fetchAuthEnabled(): Promise<boolean> {
  const status = await fetchAuthStatus();
  return status.enabled;
}
