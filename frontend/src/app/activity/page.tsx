'use client';

import { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import type { ActivityLog } from '@/types';
import { ClipboardList, User, FileText, Building2, DollarSign } from 'lucide-react';

const entityIcons: Record<string, React.ElementType> = {
    employee: User,
    document: FileText,
    company: Building2,
    salary: DollarSign,
};

const actionColors: Record<string, string> = {
    create: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400',
    update: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-400',
    delete: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-400',
};

export default function ActivityPage() {
    const [activities, setActivities] = useState<ActivityLog[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        api.activity.list(50).then(res => {
            setActivities(res.data || []);
        }).catch(() => { }).finally(() => setLoading(false));
    }, []);

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-2xl font-bold text-foreground">Activity Log</h1>
                <p className="text-sm text-muted-foreground mt-1">Audit trail of all system changes</p>
            </div>

            <div className="bg-card rounded-xl border border-border shadow-sm overflow-hidden">
                {loading ? (
                    <div className="flex items-center justify-center py-16">
                        <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin" />
                    </div>
                ) : activities.length === 0 ? (
                    <div className="flex flex-col items-center justify-center py-16 gap-3">
                        <ClipboardList className="h-10 w-10 text-muted-foreground/40" />
                        <p className="text-muted-foreground text-sm">No activity recorded yet</p>
                        <p className="text-xs text-muted-foreground">Actions like creating, updating, or deleting records will appear here</p>
                    </div>
                ) : (
                    <div className="divide-y divide-border">
                        {activities.map(a => {
                            const Icon = entityIcons[a.entityType] || ClipboardList;
                            return (
                                <div key={a.id} className="flex items-start gap-3 px-4 py-3 hover:bg-muted/20 transition-colors">
                                    <div className="mt-0.5 w-8 h-8 rounded-lg bg-muted flex items-center justify-center shrink-0">
                                        <Icon className="h-4 w-4 text-muted-foreground" />
                                    </div>
                                    <div className="flex-1 min-w-0">
                                        <p className="text-sm text-foreground">
                                            <span className="font-medium">{a.userName}</span>{' '}
                                            <span className={`inline-flex px-1.5 py-0.5 rounded text-[11px] font-medium ${actionColors[a.action] || 'bg-muted text-muted-foreground'}`}>
                                                {a.action}
                                            </span>{' '}
                                            <span className="text-muted-foreground">a {a.entityType}</span>
                                        </p>
                                        {a.details && (
                                            <p className="text-xs text-muted-foreground mt-0.5 truncate">
                                                {typeof a.details === 'object' ? JSON.stringify(a.details) : String(a.details)}
                                            </p>
                                        )}
                                    </div>
                                    <span className="text-xs text-muted-foreground whitespace-nowrap shrink-0">
                                        {new Date(a.createdAt).toLocaleString()}
                                    </span>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>
        </div>
    );
}
