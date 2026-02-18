'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useUser } from '@/hooks/use-user';
import { api } from '@/lib/api';
import type { AdminUser } from '@/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Trash2, Shield, Eye } from 'lucide-react';
import { toast } from 'sonner';

export default function UsersPage() {
    const { user, isAdmin, loading: authLoading } = useUser();
    const router = useRouter();
    const [users, setUsers] = useState<AdminUser[]>([]);
    const [loading, setLoading] = useState(true);

    const fetchUsers = useCallback(async () => {
        try {
            const res = await api.users.list();
            setUsers(res.data);
        } catch {
            toast.error('Failed to fetch users');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        if (!authLoading && !isAdmin) {
            router.push('/');
            return;
        }
        if (!authLoading && isAdmin) {
            fetchUsers();
        }
    }, [authLoading, isAdmin, router, fetchUsers]);

    const handleRoleChange = async (targetId: string, newRole: string) => {
        try {
            await api.users.updateRole(targetId, newRole);
            setUsers(prev => prev.map(u =>
                u.id === targetId ? { ...u, role: newRole } : u
            ));
            toast.success('Role updated successfully');
        } catch (err: unknown) {
            const message = err instanceof Error ? err.message : 'Failed to update role';
            toast.error(message);
        }
    };

    const handleDelete = async (targetId: string) => {
        try {
            await api.users.delete(targetId);
            setUsers(prev => prev.filter(u => u.id !== targetId));
            toast.success('User deleted successfully');
        } catch (err: unknown) {
            const message = err instanceof Error ? err.message : 'Failed to delete user';
            toast.error(message);
        }
    };

    if (authLoading || loading) {
        return (
            <div className="flex items-center justify-center py-20">
                <div className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin" />
            </div>
        );
    }

    if (!isAdmin) return null;

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-2xl font-bold text-foreground">User Management</h1>
                <p className="text-sm text-muted-foreground mt-1">
                    Manage user accounts and roles. New users register as viewers — promote to admin here.
                </p>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle className="text-lg">Users ({users.length})</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="overflow-x-auto">
                        <table className="w-full text-sm">
                            <thead>
                                <tr className="border-b border-border">
                                    <th className="text-left py-3 px-4 font-medium text-muted-foreground">Name</th>
                                    <th className="text-left py-3 px-4 font-medium text-muted-foreground">Email</th>
                                    <th className="text-left py-3 px-4 font-medium text-muted-foreground">Role</th>
                                    <th className="text-left py-3 px-4 font-medium text-muted-foreground">Joined</th>
                                    <th className="text-right py-3 px-4 font-medium text-muted-foreground">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {users.map((u) => {
                                    const isSelf = u.id === user?.id;
                                    return (
                                        <tr key={u.id} className="border-b border-border/50 hover:bg-accent/30">
                                            <td className="py-3 px-4">
                                                <div className="flex items-center gap-3">
                                                    <div className="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white text-xs font-bold">
                                                        {u.name.charAt(0).toUpperCase()}
                                                    </div>
                                                    <span className="font-medium text-foreground">
                                                        {u.name}
                                                        {isSelf && <span className="text-xs text-muted-foreground ml-2">(you)</span>}
                                                    </span>
                                                </div>
                                            </td>
                                            <td className="py-3 px-4 text-muted-foreground">{u.email}</td>
                                            <td className="py-3 px-4">
                                                {isSelf ? (
                                                    <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-950/40 text-blue-700 dark:text-blue-400">
                                                        <Shield className="h-3 w-3" /> {u.role}
                                                    </span>
                                                ) : (
                                                    <Select
                                                        value={u.role}
                                                        onValueChange={(val) => handleRoleChange(u.id, val)}
                                                    >
                                                        <SelectTrigger className="w-32 h-8 text-xs">
                                                            <SelectValue />
                                                        </SelectTrigger>
                                                        <SelectContent>
                                                            <SelectItem value="admin">
                                                                <span className="flex items-center gap-1.5">
                                                                    <Shield className="h-3 w-3" /> Admin
                                                                </span>
                                                            </SelectItem>
                                                            <SelectItem value="viewer">
                                                                <span className="flex items-center gap-1.5">
                                                                    <Eye className="h-3 w-3" /> Viewer
                                                                </span>
                                                            </SelectItem>
                                                        </SelectContent>
                                                    </Select>
                                                )}
                                            </td>
                                            <td className="py-3 px-4 text-muted-foreground">
                                                {new Date(u.createdAt).toLocaleDateString()}
                                            </td>
                                            <td className="py-3 px-4 text-right">
                                                {!isSelf && (
                                                    <AlertDialog>
                                                        <AlertDialogTrigger asChild>
                                                            <Button variant="ghost" size="icon" className="h-8 w-8 text-red-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950/30">
                                                                <Trash2 className="h-4 w-4" />
                                                            </Button>
                                                        </AlertDialogTrigger>
                                                        <AlertDialogContent>
                                                            <AlertDialogHeader>
                                                                <AlertDialogTitle>Delete User</AlertDialogTitle>
                                                                <AlertDialogDescription>
                                                                    Are you sure you want to delete <strong>{u.name}</strong> ({u.email})?
                                                                    This action cannot be undone.
                                                                </AlertDialogDescription>
                                                            </AlertDialogHeader>
                                                            <AlertDialogFooter>
                                                                <AlertDialogCancel>Cancel</AlertDialogCancel>
                                                                <AlertDialogAction
                                                                    onClick={() => handleDelete(u.id)}
                                                                    className="bg-red-600 hover:bg-red-700"
                                                                >
                                                                    Delete
                                                                </AlertDialogAction>
                                                            </AlertDialogFooter>
                                                        </AlertDialogContent>
                                                    </AlertDialog>
                                                )}
                                            </td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                        {users.length === 0 && (
                            <div className="text-center py-10 text-muted-foreground">No users found</div>
                        )}
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
