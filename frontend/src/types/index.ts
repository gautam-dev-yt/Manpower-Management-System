export interface Company {
    id: string;
    name: string;
    createdAt: string;
    updatedAt: string;
}

export interface Employee {
    id: string;
    companyId: string;
    name: string;
    trade: string;
    mobile: string;
    joiningDate: string; // ISO date string
    photoUrl?: string;
    createdAt: string;
    updatedAt: string;
}

export interface Document {
    id: string;
    employeeId: string;
    documentType: string;
    expiryDate: string; // ISO date string
    fileUrl: string;
    fileName: string;
    fileSize: number;
    fileType: string;
    lastUpdated: string;
    createdAt: string;
}

export interface DashboardMetrics {
    totalEmployees: number;
    activeDocuments: number;
    expiringSoon: number;
    expired: number;
}

export interface ExpiryAlert {
    id: string;
    employeeName: string;
    companyName: string;
    documentType: string;
    expiryDate: string;
    daysLeft: number;
    status: 'expired' | 'urgent' | 'warning';
}

// API Request/Response Types
export interface CreateEmployeeRequest {
    companyId: string;
    name: string;
    trade: string;
    mobile: string;
    joiningDate: string;
    photo?: File;
}

export interface CreateDocumentRequest {
    employeeId: string;
    documentType: string;
    expiryDate: string;
    file: File;
}

export interface EmployeeFilters {
    companyId?: string;
    trade?: string;
    status?: 'all' | 'valid' | 'expiring' | 'expired';
    search?: string;
}

export interface ApiResponse<T> {
    data: T;
    message?: string;
}

export interface ApiError {
    error: string;
    message: string;
    status: number;
}
