import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const AUTH_TOKEN_KEY = 'auth_token';
export const AUTH_USERNAME_KEY = 'username';

export type AuthUser = {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
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

export function setAuthSession(token: string, user: AuthUser): void {
  localStorage.setItem(AUTH_TOKEN_KEY, token);
  const displayName =
    [user.firstName, user.lastName].filter(Boolean).join(' ') || user.email;
  localStorage.setItem(AUTH_USERNAME_KEY, displayName);
}

export function clearAuthSession(): void {
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem(AUTH_USERNAME_KEY);
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
