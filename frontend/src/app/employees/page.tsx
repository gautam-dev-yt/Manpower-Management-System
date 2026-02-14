'use client';

import { useEffect, useState, useCallback } from 'react';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import {
    Users,
    Search,
    Plus,
    ChevronLeft,
    ChevronRight,
    Loader2,
    Phone,
} from 'lucide-react';
import { api, type PaginatedResponse } from '@/lib/api';
import type { EmployeeWithCompany, Company } from '@/types';
import { useUser } from '@/hooks/use-user';

const STATUS_COLORS: Record<string, string> = {
    valid: 'bg-green-100 dark:bg-green-950/40 text-green-800 dark:text-green-400 border-green-200 dark:border-green-800',
    expiring: 'bg-yellow-100 dark:bg-yellow-950/40 text-yellow-800 dark:text-yellow-400 border-yellow-200 dark:border-yellow-800',
    expired: 'bg-red-100 dark:bg-red-950/40 text-red-800 dark:text-red-400 border-red-200 dark:border-red-800',
    none: 'bg-muted text-muted-foreground border-border',
};

const STATUS_LABELS: Record<string, string> = {
    valid: 'Valid',
    expiring: 'Expiring',
    expired: 'Expired',
    none: 'No Docs',
};

export default function EmployeeListPage() {
    const searchParams = useSearchParams();
    const [employees, setEmployees] = useState<EmployeeWithCompany[]>([]);
    const [pagination, setPagination] = useState({ page: 1, limit: 20, total: 0, totalPages: 0 });
    const [companies, setCompanies] = useState<Company[]>([]);
    const [loading, setLoading] = useState(true);
    const { isAdmin } = useUser();

    // Filters — statusFilter is pre-filled from URL ?status=xxx (e.g. from dashboard cards)
    const [search, setSearch] = useState('');
    const [companyFilter, setCompanyFilter] = useState('');
    const [statusFilter, setStatusFilter] = useState(searchParams.get('status') || '');
    const [page, setPage] = useState(1);

    const fetchEmployees = useCallback(async () => {
        try {
            setLoading(true);
            const res = await api.employees.list({
                search: search || undefined,
                company_id: companyFilter || undefined,
                status: statusFilter || undefined,
                page,
                limit: 20,
            }) as PaginatedResponse<EmployeeWithCompany>;
            setEmployees(res.data || []);
            setPagination(res.pagination);
        } catch (err) {
            console.error('Failed to fetch employees:', err);
        } finally {
            setLoading(false);
        }
    }, [search, companyFilter, statusFilter, page]);

    const fetchCompanies = useCallback(async () => {
        try {
            const res = await api.companies.list();
            setCompanies(res.data || []);
        } catch (err) {
            console.error('Failed to fetch companies:', err);
        }
    }, []);

    useEffect(() => {
        fetchCompanies();
    }, [fetchCompanies]);

    useEffect(() => {
        fetchEmployees();
    }, [fetchEmployees]);

    // Reset page when filters change
    useEffect(() => {
        setPage(1);
    }, [search, companyFilter, statusFilter]);

    return (
        <div className="max-w-7xl mx-auto space-y-6">
            {/* Header */}
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                <div>
                    <h1 className="text-2xl sm:text-3xl font-bold text-foreground">Employees</h1>
                    <p className="text-muted-foreground text-sm mt-1">
                        {pagination.total} employee{pagination.total !== 1 ? 's' : ''} total
                    </p>
                </div>
                {isAdmin && (
                    <Link href="/employees/new">
                        <Button className="w-full sm:w-auto">
                            <Plus className="h-4 w-4 mr-2" />
                            Add Employee
                        </Button>
                    </Link>
                )}
            </div>

            {/* Filters */}
            <Card>
                <CardContent className="pt-6">
                    <div className="flex flex-col sm:flex-row gap-3">
                        <div className="relative flex-1">
                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                            <Input
                                placeholder="Search by name..."
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                                className="pl-10"
                            />
                        </div>
                        <Select value={companyFilter} onValueChange={setCompanyFilter}>
                            <SelectTrigger className="w-full sm:w-[200px]">
                                <SelectValue placeholder="All Companies" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="all">All Companies</SelectItem>
                                {companies.map((c) => (
                                    <SelectItem key={c.id} value={c.id}>
                                        {c.name}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                        <Select value={statusFilter} onValueChange={setStatusFilter}>
                            <SelectTrigger className="w-full sm:w-[180px]">
                                <SelectValue placeholder="All Statuses" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="all">All Statuses</SelectItem>
                                <SelectItem value="valid">Valid</SelectItem>
                                <SelectItem value="expiring">Expiring Soon</SelectItem>
                                <SelectItem value="expired">Expired</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>
                </CardContent>
            </Card>

            {/* Employee List */}
            {loading ? (
                <div className="flex justify-center py-12">
                    <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
                </div>
            ) : employees.length === 0 ? (
                <Card>
                    <CardContent className="py-12 text-center">
                        <Users className="h-12 w-12 mx-auto text-muted-foreground/50 mb-4" />
                        <h3 className="text-lg font-medium text-foreground">No employees found</h3>
                        <p className="text-muted-foreground text-sm mt-1">
                            {search || companyFilter || statusFilter
                                ? 'Try adjusting your filters.'
                                : 'Get started by adding your first employee.'}
                        </p>
                    </CardContent>
                </Card>
            ) : (
                <div className="space-y-3">
                    {employees.map((emp) => (
                        <Link key={emp.id} href={`/employees/${emp.id}`}>
                            <Card className="hover:shadow-md transition-shadow cursor-pointer mb-3">
                                <CardContent className="py-4 px-5">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center gap-4 flex-1 min-w-0">
                                            {/* Avatar */}
                                            <div className="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-950/40 flex items-center justify-center flex-shrink-0">
                                                <span className="text-blue-700 dark:text-blue-400 font-semibold text-sm">
                                                    {emp.name
                                                        .split(' ')
                                                        .map((n) => n[0])
                                                        .join('')
                                                        .toUpperCase()
                                                        .slice(0, 2)}
                                                </span>
                                            </div>

                                            {/* Info */}
                                            <div className="min-w-0 flex-1">
                                                <h3 className="font-semibold text-foreground truncate">{emp.name}</h3>
                                                <div className="flex items-center gap-2 text-sm text-muted-foreground flex-wrap">
                                                    <span>{emp.trade}</span>
                                                    <span className="hidden sm:inline">•</span>
                                                    <span className="hidden sm:inline">{emp.companyName}</span>
                                                </div>
                                            </div>
                                        </div>

                                        {/* Right side — expiry countdown + status badge */}
                                        <div className="flex items-center gap-3 ml-4">
                                            <div className="hidden md:flex items-center gap-1.5 text-sm text-muted-foreground">
                                                <Phone className="h-3.5 w-3.5" />
                                                <span>{emp.mobile}</span>
                                            </div>
                                            {/* Expiry countdown — derived from primary document */}
                                            {emp.expiryDaysLeft != null && (
                                                <span className={`hidden sm:inline text-xs font-medium ${emp.expiryDaysLeft < 0
                                                    ? 'text-red-600 dark:text-red-400'
                                                    : emp.expiryDaysLeft <= 30
                                                        ? 'text-amber-600 dark:text-amber-400'
                                                        : 'text-emerald-600 dark:text-emerald-400'
                                                    }`}>
                                                    {emp.expiryDaysLeft < 0
                                                        ? `${Math.abs(emp.expiryDaysLeft)}d ago`
                                                        : `${emp.expiryDaysLeft}d left`}
                                                </span>
                                            )}
                                            <Badge
                                                variant="outline"
                                                className={STATUS_COLORS[emp.docStatus] || STATUS_COLORS.none}
                                            >
                                                {STATUS_LABELS[emp.docStatus] || 'No Docs'}
                                            </Badge>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        </Link>
                    ))}
                </div>
            )}

            {/* Pagination */}
            {pagination.totalPages > 1 && (
                <div className="flex items-center justify-between pt-2">
                    <p className="text-sm text-muted-foreground">
                        Page {pagination.page} of {pagination.totalPages}
                    </p>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={page <= 1}
                            onClick={() => setPage((p) => Math.max(1, p - 1))}
                        >
                            <ChevronLeft className="h-4 w-4 mr-1" /> Previous
                        </Button>
                        <Button
                            variant="outline"
                            size="sm"
                            disabled={page >= pagination.totalPages}
                            onClick={() => setPage((p) => p + 1)}
                        >
                            Next <ChevronRight className="h-4 w-4 ml-1" />
                        </Button>
                    </div>
                </div>
            )}
        </div>
    );
}
