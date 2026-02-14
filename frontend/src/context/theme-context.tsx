'use client';

import { createContext, useContext, useEffect, useState, useCallback } from 'react';

type Theme = 'light' | 'dark';

interface ThemeContextValue {
    theme: Theme;
    toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);

/**
 * ThemeProvider manages light/dark mode state.
 * Persists preference in localStorage and applies .dark class to <html>.
 */
export function ThemeProvider({ children }: { children: React.ReactNode }) {
    const [theme, setTheme] = useState<Theme>('light');

    // Load saved preference on mount
    useEffect(() => {
        const saved = localStorage.getItem('theme') as Theme | null;
        if (saved) {
            setTheme(saved);
            document.documentElement.classList.toggle('dark', saved === 'dark');
        }
    }, []);

    const toggleTheme = useCallback(() => {
        setTheme((current) => {
            const next = current === 'light' ? 'dark' : 'light';
            localStorage.setItem('theme', next);
            document.documentElement.classList.toggle('dark', next === 'dark');
            return next;
        });
    }, []);

    return (
        <ThemeContext.Provider value={{ theme, toggleTheme }}>
            {children}
        </ThemeContext.Provider>
    );
}

export function useTheme() {
    const context = useContext(ThemeContext);
    if (!context) throw new Error('useTheme must be used within ThemeProvider');
    return context;
}
