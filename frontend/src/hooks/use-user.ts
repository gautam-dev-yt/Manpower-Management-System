'use client';

import { useState, useEffect } from 'react';
import { getUser, User } from '@/lib/auth';

export function useUser() {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // eslint-disable-next-line react-hooks/exhaustive-deps
        setUser(getUser());
        setLoading(false);
    }, []);

    return { user, loading, isAdmin: user?.role === 'admin' };
}
