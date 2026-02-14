'use client';

import { createContext, useContext, useEffect, useState, useCallback } from 'react';
import { useRouter, usePathname } from 'next/navigation';

interface User {
    id: string;
    email: string;
    name: string;
    role: string;
}

interface AuthContextValue {
    user: User | null;
    token: string | null;
    loading: boolean;
    login: (email: string, password: string) => Promise<void>;
    register: (name: string, email: string, password: string) => Promise<void>;
    logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Pages that don't require authentication
const PUBLIC_PATHS = ['/login', '/register'];

/**
 * AuthProvider manages authentication state across the entire app.
 * On mount, checks for a stored JWT token and validates it.
 * Redirects unauthenticated users to /login automatically.
 */
export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const router = useRouter();
    const pathname = usePathname();

    // Validate stored token on initial load
    useEffect(() => {
        const storedToken = localStorage.getItem('token');
        if (!storedToken) {
            setLoading(false);
            return;
        }

        // Verify token is still valid by hitting /auth/me
        fetch(`${API_BASE}/api/auth/me`, {
            headers: { Authorization: `Bearer ${storedToken}` },
        })
            .then((res) => {
                if (!res.ok) throw new Error('Invalid token');
                return res.json();
            })
            .then((userData: User) => {
                setUser(userData);
                setToken(storedToken);
            })
            .catch(() => {
                // Token expired or invalid — clear it
                localStorage.removeItem('token');
            })
            .finally(() => setLoading(false));
    }, []);

    // Redirect unauthenticated users to login
    useEffect(() => {
        if (loading) return;
        if (!user && !PUBLIC_PATHS.includes(pathname)) {
            router.replace('/login');
        }
    }, [user, loading, pathname, router]);

    const login = useCallback(async (email: string, password: string) => {
        const res = await fetch(`${API_BASE}/api/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        });

        if (!res.ok) {
            const err = await res.json().catch(() => ({ message: 'Login failed' }));
            throw new Error(err.message || 'Invalid credentials');
        }

        const data = await res.json();
        localStorage.setItem('token', data.token);
        setToken(data.token);
        setUser(data.user);
        router.push('/');
    }, [router]);

    const register = useCallback(async (name: string, email: string, password: string) => {
        const res = await fetch(`${API_BASE}/api/auth/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, email, password }),
        });

        if (!res.ok) {
            const err = await res.json().catch(() => ({ message: 'Registration failed' }));
            throw new Error(err.message || 'Could not create account');
        }

        const data = await res.json();
        localStorage.setItem('token', data.token);
        setToken(data.token);
        setUser(data.user);
        router.push('/');
    }, [router]);

    const logout = useCallback(() => {
        localStorage.removeItem('token');
        setToken(null);
        setUser(null);
        router.push('/login');
    }, [router]);

    return (
        <AuthContext.Provider value={{ user, token, loading, login, register, logout }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) throw new Error('useAuth must be used within AuthProvider');
    return context;
}
