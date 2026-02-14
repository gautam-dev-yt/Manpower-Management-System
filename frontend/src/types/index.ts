// ── Core Entities ─────────────────────────────────────────────

export interface Company {
    id: string;
    name: string;
    createdAt: string;
    updatedAt: string;
}

export interface CompanySummary {
    id: string;
    name: string;
    employeeCount: number;
}

export interface Employee {
    id: string;
    companyId: string;
    name: string;
    trade: string;
    mobile: string;
    joiningDate: string;
    photoUrl?: string | null;
    gender?: string | null;
    dateOfBirth?: string | null;
    nationality?: string | null;
    passportNumber?: string | null;
    nativeLocation?: string | null;
    currentLocation?: string | null;
    salary?: number | null;
    status: string; // active, inactive, on_leave
    createdAt: string;
    updatedAt: string;
}

export interface EmployeeWithCompany extends Employee {
    companyName: string;
    docStatus: 'valid' | 'expiring' | 'expired' | 'none';
    expiryDaysLeft?: number | null;   // days until primary doc expires (negative = overdue)
    primaryDocType?: string | null;   // e.g. "Work Permit"
}

export interface Document {
    id: string;
    employeeId: string;
    documentType: string;
    expiryDate?: string | null;   // nullable — docs without expiry are allowed
    isPrimary: boolean;           // only one per employee — the tracked doc
    fileUrl: string;
    fileName: string;
    fileSize: number;
    fileType: string;
    lastUpdated: string;
    createdAt: string;
}

// ── Salary ────────────────────────────────────────────────────

export interface SalaryRecord {
    id: string;
    employeeId: string;
    month: number;
    year: number;
    amount: number;
    status: 'pending' | 'paid' | 'partial';
    paidDate?: string | null;
    notes?: string | null;
    createdAt: string;
    updatedAt: string;
}

export interface SalaryRecordWithEmployee extends SalaryRecord {
    employeeName: string;
    companyName: string;
}

export interface SalarySummary {
    totalAmount: number;
    paidAmount: number;
    pendingCount: number;
    paidCount: number;
    partialCount: number;
    totalCount: number;
}

// ── Notifications ─────────────────────────────────────────────

export interface Notification {
    id: string;
    userId: string;
    title: string;
    message: string;
    type: string;
    read: boolean;
    entityType?: string | null;
    entityId?: string | null;
    createdAt: string;
}

// ── Activity ──────────────────────────────────────────────────

export interface ActivityLog {
    id: string;
    userId: string;
    userName: string;
    action: string;
    entityType: string;
    entityId: string;
    details?: Record<string, unknown>;
    createdAt: string;
}

// ── Dashboard ─────────────────────────────────────────────────

export interface DashboardMetrics {
    totalEmployees: number;
    activeDocuments: number;
    expiringSoon: number;
    expired: number;
}

export interface ExpiryAlert {
    documentId: string;
    employeeId: string;
    employeeName: string;
    companyName: string;
    documentType: string;
    expiryDate: string;
    daysLeft: number;
    status: 'expired' | 'urgent' | 'warning';
}

// ── API Requests ──────────────────────────────────────────────

export interface CreateEmployeeRequest {
    companyId: string;
    name: string;
    trade: string;
    mobile: string;
    joiningDate: string;
    photoUrl?: string;
    gender?: string;
    dateOfBirth?: string;
    nationality?: string;
    passportNumber?: string;
    nativeLocation?: string;
    currentLocation?: string;
    salary?: number;
    status?: string;
}

export interface CreateDocumentRequest {
    documentType: string;
    expiryDate?: string;
    fileUrl: string;
    fileName: string;
    fileSize: number;
    fileType: string;
}

// ── API Responses ─────────────────────────────────────────────

export interface ApiResponse<T> {
    data: T;
    message?: string;
}

export interface ApiError {
    error: string;
    message: string;
    status: number;
}

export interface EmployeeFilters {
    company_id?: string;
    trade?: string;
    status?: string;         // document status: valid, expiring, expired
    emp_status?: string;     // employee status: active, inactive, on_leave
    nationality?: string;
    search?: string;
    sort_by?: string;
    sort_order?: string;
    page?: number;
    limit?: number;
}
