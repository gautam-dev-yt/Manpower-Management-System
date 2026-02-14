'use client';

import { useEffect, useState, useCallback } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
    AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent,
    AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import {
    ArrowLeft, Pencil, Trash2, Phone, Building2, Calendar,
    Briefcase, FileText, Plus, Loader2, AlertTriangle, XCircle, CheckCircle, Star, RefreshCw,
} from 'lucide-react';
import { api } from '@/lib/api';
import type { EmployeeWithCompany, Document as DocType } from '@/types';
import { toast } from 'sonner';
import { AddDocumentDialog, EditDocumentDialog } from '@/components/documents/add-document-dialog';
import { RenewDocumentDialog } from '@/components/documents/renew-document-dialog';
import { useUser } from '@/hooks/use-user';

/** Calculate document expiry status and days remaining */
function getDocStatus(expiryDate: string) {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const expiry = new Date(expiryDate);
    const diffDays = Math.ceil((expiry.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));

    if (diffDays < 0) return { label: 'Expired', color: 'bg-red-100 dark:bg-red-950/40 text-red-800 dark:text-red-400', icon: XCircle, days: diffDays };
    if (diffDays <= 7) return { label: 'Urgent', color: 'bg-orange-100 dark:bg-orange-950/40 text-orange-800 dark:text-orange-400', icon: AlertTriangle, days: diffDays };
    if (diffDays <= 30) return { label: 'Warning', color: 'bg-yellow-100 dark:bg-yellow-950/40 text-yellow-800 dark:text-yellow-400', icon: AlertTriangle, days: diffDays };
    return { label: 'Valid', color: 'bg-green-100 dark:bg-green-950/40 text-green-800 dark:text-green-400', icon: CheckCircle, days: diffDays };
}

export default function EmployeeDetailPage() {
    const params = useParams();
    const router = useRouter();
    const id = params.id as string;

    const [employee, setEmployee] = useState<EmployeeWithCompany | null>(null);
    const [documents, setDocuments] = useState<DocType[]>([]);
    const [loading, setLoading] = useState(true);
    const [deleting, setDeleting] = useState(false);
    const { isAdmin } = useUser();
    const [showAddDoc, setShowAddDoc] = useState(false);
    const [editingDoc, setEditingDoc] = useState<DocType | null>(null);
    const [renewingDoc, setRenewingDoc] = useState<DocType | null>(null);

    const fetchEmployee = useCallback(async () => {
        try {
            setLoading(true);
            const res = await api.employees.get(id);
            setEmployee(res.data);
            setDocuments(res.documents || []);
        } catch {
            toast.error('Failed to load employee details');
        } finally {
            setLoading(false);
        }
    }, [id]);

    useEffect(() => { fetchEmployee(); }, [fetchEmployee]);

    const handleDelete = async () => {
        try {
            setDeleting(true);
            await api.employees.delete(id);
            toast.success('Employee deleted successfully');
            router.push('/employees');
        } catch {
            toast.error('Failed to delete employee');
        } finally {
            setDeleting(false);
        }
    };

    const handleDeleteDocument = async (docId: string) => {
        try {
            await api.documents.delete(docId);
            toast.success('Document deleted');
            fetchEmployee();
        } catch {
            toast.error('Failed to delete document');
        }
    };

    if (loading) {
        return (
            <div className="flex justify-center py-20">
                <Loader2 className="h-10 w-10 animate-spin text-blue-600" />
            </div>
        );
    }

    if (!employee) {
        return (
            <div className="max-w-2xl mx-auto text-center py-20">
                <h2 className="text-xl font-semibold text-muted-foreground">Employee not found</h2>
                <Link href="/employees">
                    <Button className="mt-4">Back to Employees</Button>
                </Link>
            </div>
        );
    }

    return (
        <div className="max-w-5xl mx-auto space-y-6">
            {/* Back button */}
            <Link href="/employees" className="inline-flex items-center text-sm text-muted-foreground hover:text-foreground">
                <ArrowLeft className="h-4 w-4 mr-1" /> Back to employees
            </Link>

            {/* Employee Info Card */}
            <Card className="border-border/60">
                <CardContent className="pt-6">
                    <div className="flex flex-col sm:flex-row gap-6">
                        {/* Avatar */}
                        <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-blue-100 to-indigo-100 dark:from-blue-950 dark:to-indigo-950 flex items-center justify-center flex-shrink-0">
                            <span className="text-blue-700 dark:text-blue-300 font-bold text-2xl">
                                {employee.name.split(' ').map((n) => n[0]).join('').toUpperCase().slice(0, 2)}
                            </span>
                        </div>

                        {/* Details */}
                        <div className="flex-1 min-w-0">
                            <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
                                <div>
                                    <h1 className="text-2xl font-bold text-foreground">{employee.name}</h1>
                                    <div className="flex flex-wrap items-center gap-x-4 gap-y-2 mt-2 text-sm text-muted-foreground">
                                        <span className="flex items-center gap-1.5"><Briefcase className="h-4 w-4" /> {employee.trade}</span>
                                        <span className="flex items-center gap-1.5"><Building2 className="h-4 w-4" /> {employee.companyName}</span>
                                        <span className="flex items-center gap-1.5"><Phone className="h-4 w-4" /> {employee.mobile}</span>
                                        <span className="flex items-center gap-1.5"><Calendar className="h-4 w-4" /> Joined {employee.joiningDate}</span>
                                    </div>
                                </div>

                                {isAdmin && (
                                    <div className="flex gap-2 flex-shrink-0">
                                        <Link href={`/employees/${id}/edit`}>
                                            <Button variant="outline" size="sm"><Pencil className="h-4 w-4 mr-1" /> Edit</Button>
                                        </Link>
                                        <AlertDialog>
                                            <AlertDialogTrigger asChild>
                                                <Button variant="outline" size="sm" className="text-red-600 dark:text-red-400 hover:text-red-700">
                                                    <Trash2 className="h-4 w-4 mr-1" /> Delete
                                                </Button>
                                            </AlertDialogTrigger>
                                            <AlertDialogContent>
                                                <AlertDialogHeader>
                                                    <AlertDialogTitle>Delete {employee.name}?</AlertDialogTitle>
                                                    <AlertDialogDescription>
                                                        This will permanently delete this employee and all their documents.
                                                    </AlertDialogDescription>
                                                </AlertDialogHeader>
                                                <AlertDialogFooter>
                                                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                                                    <AlertDialogAction onClick={handleDelete} disabled={deleting} className="bg-red-600 hover:bg-red-700">
                                                        {deleting ? 'Deleting...' : 'Delete'}
                                                    </AlertDialogAction>
                                                </AlertDialogFooter>
                                            </AlertDialogContent>
                                        </AlertDialog>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                </CardContent>
            </Card>

            {/* Documents Section */}
            <Card className="border-border/60">
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <div>
                            <CardTitle className="flex items-center gap-2"><FileText className="h-5 w-5" /> Documents</CardTitle>
                            <CardDescription>{documents.length} document{documents.length !== 1 ? 's' : ''}</CardDescription>
                        </div>
                        {isAdmin && (
                            <Button size="sm" onClick={() => setShowAddDoc(true)}>
                                <Plus className="h-4 w-4 mr-1" /> Add Document
                            </Button>
                        )}
                    </div>
                </CardHeader>
                <CardContent>
                    {documents.length === 0 ? (
                        <div className="text-center py-8">
                            <FileText className="h-10 w-10 mx-auto text-muted-foreground/40 mb-3" />
                            <p className="text-muted-foreground">No documents yet. Add one to start tracking.</p>
                        </div>
                    ) : (
                        <div className="space-y-3">
                            {documents.map((doc) => {
                                const hasExpiry = !!doc.expiryDate;
                                const status = hasExpiry ? getDocStatus(doc.expiryDate!) : null;
                                const StatusIcon = status?.icon || FileText;
                                const isPrimary = doc.isPrimary;

                                return (
                                    <div
                                        key={doc.id}
                                        className={`flex items-center justify-between p-4 rounded-lg border transition-shadow hover:shadow-sm ${isPrimary
                                            ? 'border-amber-400 dark:border-amber-600 bg-amber-50/50 dark:bg-amber-950/20'
                                            : 'border-border/60'
                                            }`}
                                    >
                                        <div className="flex items-center gap-4 flex-1 min-w-0">
                                            <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${status?.color || 'bg-muted text-muted-foreground'}`}>
                                                <StatusIcon className="h-5 w-5" />
                                            </div>
                                            <div className="min-w-0">
                                                <div className="flex items-center gap-2">
                                                    <h4 className="font-medium text-foreground">{doc.documentType}</h4>
                                                    {isPrimary && (
                                                        <Badge variant="outline" className="bg-amber-100 dark:bg-amber-950/40 text-amber-700 dark:text-amber-400 border-amber-300 dark:border-amber-700 text-xs">
                                                            <Star className="h-3 w-3 mr-1 fill-current" /> Primary
                                                        </Badge>
                                                    )}
                                                </div>
                                                <p className="text-sm text-muted-foreground">
                                                    {hasExpiry ? `Expires: ${doc.expiryDate}` : 'No expiry date'}
                                                </p>
                                            </div>
                                        </div>

                                        <div className="flex items-center gap-3 ml-4">
                                            {/* Per-doc expiry countdown */}
                                            {hasExpiry && status && (
                                                <div className="text-right hidden sm:block">
                                                    <div className={`text-lg font-bold ${status.days < 0 ? 'text-red-600' : status.days <= 7 ? 'text-orange-600' : status.days <= 30 ? 'text-yellow-600' : 'text-green-600'}`}>
                                                        {Math.abs(status.days)}
                                                    </div>
                                                    <div className="text-xs text-muted-foreground">
                                                        {status.days < 0 ? 'days overdue' : 'days left'}
                                                    </div>
                                                </div>
                                            )}

                                            {/* Status badge */}
                                            <Badge variant="outline" className={status?.color || 'bg-muted text-muted-foreground border-border'}>
                                                {status?.label || 'No Expiry'}
                                            </Badge>

                                            {/* Admin Actions */}
                                            {isAdmin && (
                                                <>
                                                    {/* Toggle primary (only for docs that have an expiry date) */}
                                                    {hasExpiry && (
                                                        <Button
                                                            variant={isPrimary ? 'default' : 'outline'}
                                                            size="icon"
                                                            title={isPrimary ? 'Unset as primary' : 'Set as primary document'}
                                                            className={isPrimary ? 'bg-amber-500 hover:bg-amber-600 text-white h-8 w-8' : 'h-8 w-8'}
                                                            onClick={async () => {
                                                                try {
                                                                    await api.documents.togglePrimary(doc.id);
                                                                    toast.success(isPrimary ? 'Primary document unset' : 'Set as primary document');
                                                                    fetchEmployee();
                                                                } catch {
                                                                    toast.error('Failed to toggle primary');
                                                                }
                                                            }}
                                                        >
                                                            <Star className={`h-4 w-4 ${isPrimary ? 'fill-current' : ''}`} />
                                                        </Button>
                                                    )}

                                                    {/* Renew button */}
                                                    {hasExpiry && (
                                                        <Button
                                                            variant="ghost" size="icon"
                                                            title="Renew Document"
                                                            className="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
                                                            onClick={() => setRenewingDoc(doc)}
                                                        >
                                                            <RefreshCw className="h-4 w-4" />
                                                        </Button>
                                                    )}

                                                    {/* Edit button */}
                                                    <Button
                                                        variant="ghost" size="icon"
                                                        className="text-muted-foreground hover:text-foreground"
                                                        onClick={() => setEditingDoc(doc)}
                                                    >
                                                        <Pencil className="h-4 w-4" />
                                                    </Button>

                                                    {/* Delete button */}
                                                    <AlertDialog>
                                                        <AlertDialogTrigger asChild>
                                                            <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-red-600 dark:hover:text-red-400">
                                                                <Trash2 className="h-4 w-4" />
                                                            </Button>
                                                        </AlertDialogTrigger>
                                                        <AlertDialogContent>
                                                            <AlertDialogHeader>
                                                                <AlertDialogTitle>Delete {doc.documentType}?</AlertDialogTitle>
                                                                <AlertDialogDescription>This will permanently delete this document record.</AlertDialogDescription>
                                                            </AlertDialogHeader>
                                                            <AlertDialogFooter>
                                                                <AlertDialogCancel>Cancel</AlertDialogCancel>
                                                                <AlertDialogAction onClick={() => handleDeleteDocument(doc.id)} className="bg-red-600 hover:bg-red-700">
                                                                    Delete
                                                                </AlertDialogAction>
                                                            </AlertDialogFooter>
                                                        </AlertDialogContent>
                                                    </AlertDialog>
                                                </>
                                            )}
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Add Document Dialog */}
            <AddDocumentDialog
                employeeId={id}
                open={showAddDoc}
                onOpenChange={setShowAddDoc}
                onSuccess={fetchEmployee}
            />

            {/* Edit Document Dialog */}
            {editingDoc && (
                <EditDocumentDialog
                    document={editingDoc}
                    open={!!editingDoc}
                    onOpenChange={(open) => !open && setEditingDoc(null)}
                    onSuccess={fetchEmployee}
                />
            )}

            {/* Renew Document Dialog */}
            {renewingDoc && (
                <RenewDocumentDialog
                    document={renewingDoc}
                    open={!!renewingDoc}
                    onOpenChange={(open) => !open && setRenewingDoc(null)}
                    onSuccess={fetchEmployee}
                />
            )}
        </div>
    );
}
