'use client';

import { useState, useEffect, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger, DialogFooter,
} from '@/components/ui/dialog';
import {
    AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent,
    AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Building2, Plus, Pencil, Trash2, Users, Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { api } from '@/lib/api';

interface CompanyWithCount {
    id: string;
    name: string;
    employeeCount: number;
    createdAt: string;
    updatedAt: string;
}

export default function CompaniesPage() {
    const [companies, setCompanies] = useState<CompanyWithCount[]>([]);
    const [loading, setLoading] = useState(true);

    // Add/Edit dialog state
    const [dialogOpen, setDialogOpen] = useState(false);
    const [editingId, setEditingId] = useState<string | null>(null);
    const [companyName, setCompanyName] = useState('');
    const [saving, setSaving] = useState(false);

    const fetchCompanies = useCallback(async () => {
        try {
            const res = await api.companies.list();
            setCompanies(res.data as unknown as CompanyWithCount[]);
        } catch {
            toast.error('Failed to load companies');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => { fetchCompanies(); }, [fetchCompanies]);

    const handleSave = async () => {
        if (!companyName.trim()) return;
        setSaving(true);

        try {
            if (editingId) {
                await api.companies.update(editingId, { name: companyName.trim() });
                toast.success('Company updated');
            } else {
                await api.companies.create({ name: companyName.trim() });
                toast.success('Company created');
            }

            setDialogOpen(false);
            setCompanyName('');
            setEditingId(null);
            fetchCompanies();
        } catch (err) {
            toast.error(err instanceof Error ? err.message : 'Failed to save');
        } finally {
            setSaving(false);
        }
    };

    const handleDelete = async (id: string) => {
        try {
            await api.companies.delete(id);
            toast.success('Company deleted');
            fetchCompanies();
        } catch {
            toast.error('Failed to delete company');
        }
    };

    const openEdit = (company: CompanyWithCount) => {
        setEditingId(company.id);
        setCompanyName(company.name);
        setDialogOpen(true);
    };

    const openAdd = () => {
        setEditingId(null);
        setCompanyName('');
        setDialogOpen(true);
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center py-20">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold text-foreground">Companies</h1>
                    <p className="text-muted-foreground mt-1">
                        {companies.length} compan{companies.length === 1 ? 'y' : 'ies'} registered
                    </p>
                </div>

                <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
                    <DialogTrigger asChild>
                        <Button onClick={openAdd} className="gap-2">
                            <Plus className="h-4 w-4" /> Add Company
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>{editingId ? 'Edit Company' : 'Add Company'}</DialogTitle>
                            <DialogDescription>
                                {editingId ? 'Update the company name.' : 'Enter the name for the new company.'}
                            </DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4 py-4">
                            <div className="space-y-2">
                                <Label htmlFor="company-name">Company Name</Label>
                                <Input
                                    id="company-name"
                                    placeholder="ABC Construction Co."
                                    value={companyName}
                                    onChange={(e) => setCompanyName(e.target.value)}
                                    autoFocus
                                    onKeyDown={(e) => e.key === 'Enter' && handleSave()}
                                />
                            </div>
                        </div>
                        <DialogFooter>
                            <Button variant="outline" onClick={() => setDialogOpen(false)}>Cancel</Button>
                            <Button onClick={handleSave} disabled={saving || !companyName.trim()}>
                                {saving ? <Loader2 className="h-4 w-4 mr-2 animate-spin" /> : null}
                                {editingId ? 'Save Changes' : 'Create Company'}
                            </Button>
                        </DialogFooter>
                    </DialogContent>
                </Dialog>
            </div>

            {/* Company Cards Grid */}
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {companies.map((company) => (
                    <Card key={company.id} className="group hover:shadow-md transition-shadow border-border/60">
                        <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
                            <div className="flex items-center gap-3">
                                <div className="w-10 h-10 rounded-lg bg-blue-100 dark:bg-blue-950/40 flex items-center justify-center">
                                    <Building2 className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                                </div>
                                <CardTitle className="text-base font-semibold">{company.name}</CardTitle>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                    <Users className="h-4 w-4" />
                                    <span>{company.employeeCount} employee{company.employeeCount !== 1 ? 's' : ''}</span>
                                </div>

                                <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                                    <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => openEdit(company)}>
                                        <Pencil className="h-3.5 w-3.5" />
                                    </Button>

                                    <AlertDialog>
                                        <AlertDialogTrigger asChild>
                                            <Button variant="ghost" size="icon" className="h-8 w-8 text-red-600 dark:text-red-400">
                                                <Trash2 className="h-3.5 w-3.5" />
                                            </Button>
                                        </AlertDialogTrigger>
                                        <AlertDialogContent>
                                            <AlertDialogHeader>
                                                <AlertDialogTitle>Delete {company.name}?</AlertDialogTitle>
                                                <AlertDialogDescription>
                                                    This will also delete all {company.employeeCount} employees
                                                    and their documents. This action cannot be undone.
                                                </AlertDialogDescription>
                                            </AlertDialogHeader>
                                            <AlertDialogFooter>
                                                <AlertDialogCancel>Cancel</AlertDialogCancel>
                                                <AlertDialogAction
                                                    onClick={() => handleDelete(company.id)}
                                                    className="bg-red-600 hover:bg-red-700"
                                                >
                                                    Delete
                                                </AlertDialogAction>
                                            </AlertDialogFooter>
                                        </AlertDialogContent>
                                    </AlertDialog>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Empty state */}
            {companies.length === 0 && (
                <div className="text-center py-16 space-y-3">
                    <Building2 className="h-12 w-12 text-muted-foreground/50 mx-auto" />
                    <h2 className="text-lg font-semibold text-foreground">No companies yet</h2>
                    <p className="text-muted-foreground">Add your first company to start managing employees.</p>
                    <Button onClick={openAdd} className="gap-2 mt-2">
                        <Plus className="h-4 w-4" /> Add Company
                    </Button>
                </div>
            )}
        </div>
    );
}
