import { ApiError, ApiResponse } from '@/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiClientError extends Error {
    constructor(
        message: string,
        public status: number,
        public data?: unknown
    ) {
        super(message);
        this.name = 'ApiClientError';
    }
}

async function fetcher<T>(
    endpoint: string,
    options?: RequestInit
): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;

    try {
        const response = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options?.headers,
            },
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => null);
            throw new ApiClientError(
                errorData?.message || errorData?.error || 'An error occurred',
                response.status,
                errorData
            );
        }

        return response.json();
    } catch (error) {
        if (error instanceof ApiClientError) {
            throw error;
        }
        throw new ApiClientError('Network error', 0, error);
    }
}

export const api = {
    // Health check
    health: () => fetcher<{ status: string }>('/api/health'),

    // Employee endpoints
    employees: {
        list: (params?: {
            companyId?: string;
            trade?: string;
            status?: string;
            search?: string;
            page?: number;
            limit?: number;
        }) => {
            const queryParams = new URLSearchParams();
            if (params) {
                Object.entries(params).forEach(([key, value]) => {
                    if (value !== undefined) {
                        queryParams.append(key, String(value));
                    }
                });
            }
            const queryString = queryParams.toString();
            return fetcher(`/api/employees${queryString ? `?${queryString}` : ''}`);
        },

        get: (id: string) => fetcher(`/api/employees/${id}`),

        create: (data: FormData) =>
            fetcher('/api/employees', {
                method: 'POST',
                body: data,
                headers: {}, // Let browser set Content-Type for FormData
            }),

        update: (id: string, data: FormData) =>
            fetcher(`/api/employees/${id}`, {
                method: 'PUT',
                body: data,
                headers: {}, // Let browser set Content-Type for FormData
            }),

        delete: (id: string) =>
            fetcher(`/api/employees/${id}`, {
                method: 'DELETE',
            }),
    },

    // Document endpoints
    documents: {
        listByEmployee: (employeeId: string) =>
            fetcher(`/api/employees/${employeeId}/documents`),

        get: (id: string) => fetcher(`/api/documents/${id}`),

        create: (data: FormData) =>
            fetcher('/api/documents', {
                method: 'POST',
                body: data,
                headers: {}, // Let browser set Content-Type for FormData
            }),

        update: (id: string, data: FormData) =>
            fetcher(`/api/documents/${id}`, {
                method: 'PUT',
                body: data,
                headers: {}, // Let browser set Content-Type for FormData
            }),

        delete: (id: string) =>
            fetcher(`/api/documents/${id}`, {
                method: 'DELETE',
            }),

        getPresignedUrl: (documentId: string) =>
            fetcher<{ url: string }>(`/api/documents/${documentId}/url`),
    },

    // Dashboard endpoints
    dashboard: {
        getMetrics: () => fetcher('/api/dashboard/metrics'),
        getExpiryAlerts: () => fetcher('/api/dashboard/expiry-alerts'),
    },

    // Company endpoints
    companies: {
        list: () => fetcher('/api/companies'),
        get: (id: string) => fetcher(`/api/companies/${id}`),
    },
};

export { ApiClientError };
