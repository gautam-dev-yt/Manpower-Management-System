export interface User {
    id: string;
    email: string;
    name: string;
    role: string;
}

export function getUser(): User | null {
    if (typeof window === 'undefined') return null;
    const token = localStorage.getItem('token');
    if (!token) return null;
    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return {
            id: payload.sub || payload.id, // standardized
            email: payload.email,
            name: payload.name,
            role: payload.role || 'viewer',
        };
    } catch {
        return null;
    }
}

export function isAdmin(): boolean {
    const user = getUser();
    return user?.role === 'admin';
}
