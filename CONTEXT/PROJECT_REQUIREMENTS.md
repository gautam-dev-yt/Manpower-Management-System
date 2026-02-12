# 📋 Manpower Management System - Complete Requirements

> **For AI Editors:** This document contains complete project requirements, features, flows, and database schema. Read this to understand what we're building.

---

## 🎯 Project Overview

### **What We're Building**

A web-based **Manpower Management System** to help a business owner track employees across companies and manage document expiries with automated alerts.

### **The Problem**

The business owner manages employees working at different companies (like a contracting/manpower supply business). Each employee has multiple documents (Visa, Passport, Emirates ID, Labor Card, etc.) that expire at different times. He needs:
- A centralized system to track all employees and their documents
- Automatic alerts when documents are about to expire (30 days, 7 days)
- Easy way to see which documents need renewal
- File storage for document copies

### **The Solution**

A clean, simple web application with:
- **Dashboard:** Overview of all employees, companies, and expiring documents
- **Employee Management:** Add, edit, view, delete employees
- **Document Management:** Upload documents with expiry dates, get automatic alerts
- **Notifications:** Daily email alerts for expiring/expired documents
- **File Storage:** Secure cloud storage for document files

### **Who Will Use It**

- **Primary User:** Business owner (single admin)
- **Usage:** Office/desktop environment
- **Scale:** ~500-1000 employees max, <15GB storage
- **Future:** May expand to multiple users/companies

---

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Frontend (Web)                    │
│              Next.js 15 + TypeScript                │
│         Dashboard | Employees | Documents           │
└────────────────────┬────────────────────────────────┘
                     │ HTTPS/REST API
┌────────────────────▼────────────────────────────────┐
│                  Backend (Go)                       │
│          Chi Router + Clean Architecture            │
│    ┌──────────┬──────────┬──────────┬──────────┐   │
│    │ Employee │ Document │Dashboard │  Cron    │   │
│    │ Service  │ Service  │ Service  │  Job     │   │
│    └──────────┴──────────┴──────────┴──────────┘   │
└───────┬─────────────┬────────────────┬──────────────┘
        │             │                │
┌───────▼──────┐ ┌────▼─────────┐ ┌───▼──────────┐
│  PostgreSQL  │ │   AWS S3     │ │  AWS SES     │
│   Database   │ │ File Storage │ │    Email     │
└──────────────┘ └──────────────┘ └──────────────┘
```

---

## 📊 Database Schema Design

### **Tables Overview**

We have 3 main tables:
1. **companies** - Store company information
2. **employees** - Store employee details
3. **documents** - Store document metadata (files in S3)

### **1. Companies Table**

```sql
CREATE TABLE companies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL UNIQUE,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_companies_name ON companies(name);

-- Sample Data
INSERT INTO companies (id, name) VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'Default Company');
```

**Purpose:** Store company information. For now, we have one company, but designed for future expansion.

---

### **2. Employees Table**

```sql
CREATE TABLE employees (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    trade           VARCHAR(100) NOT NULL,  -- Job role (e.g., "Welder", "Engineer")
    mobile          VARCHAR(20) NOT NULL,
    joining_date    DATE NOT NULL,
    photo_url       TEXT,  -- S3 URL for employee photo (optional)
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_employees_company_id ON employees(company_id);
CREATE INDEX idx_employees_name ON employees(name);
CREATE INDEX idx_employees_trade ON employees(trade);
CREATE INDEX idx_employees_joining_date ON employees(joining_date);

-- Constraints
ALTER TABLE employees ADD CONSTRAINT chk_mobile_format 
    CHECK (mobile ~ '^\+?[1-9]\d{1,14}$');
```

**Fields Explanation:**
- `id` - Unique identifier (UUID)
- `company_id` - Which company this employee belongs to
- `name` - Employee full name
- `trade` - Job role/designation (Welder, Engineer, Safety Officer, etc.)
- `mobile` - Contact number (international format)
- `joining_date` - When employee joined
- `photo_url` - Optional employee photo stored in S3
- `created_at/updated_at` - Audit timestamps

**Business Rules:**
- Name must be at least 2 characters
- Mobile must be in international format (E.164)
- Trade/job role is required
- Company reference is required

---

### **3. Documents Table**

```sql
CREATE TABLE documents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id     UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    document_type   VARCHAR(100) NOT NULL,  -- Flexible: "Visa", "Passport", "ID Card", etc.
    expiry_date     DATE NOT NULL,
    file_url        TEXT NOT NULL,  -- S3 URL for the document file
    file_name       VARCHAR(255) NOT NULL,  -- Original filename
    file_size       BIGINT NOT NULL,  -- File size in bytes
    file_type       VARCHAR(50) NOT NULL,  -- MIME type (image/jpeg, application/pdf)
    last_updated    TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_documents_employee_id ON documents(employee_id);
CREATE INDEX idx_documents_expiry_date ON documents(expiry_date);  -- Critical for alerts
CREATE INDEX idx_documents_type ON documents(document_type);

-- Composite index for expiry alerts query
CREATE INDEX idx_documents_expiry_alerts ON documents(expiry_date, employee_id);
```

**Fields Explanation:**
- `id` - Unique identifier
- `employee_id` - Which employee this document belongs to
- `document_type` - Flexible string (Visa, Passport, Emirates ID, Bank Slip, Labor Card, etc.)
- `expiry_date` - When document expires (critical field)
- `file_url` - S3 URL where actual file is stored
- `file_name` - Original filename for display
- `file_size` - Size in bytes (for validation)
- `file_type` - MIME type (for validation and display)
- `last_updated` - When document was last modified
- `created_at` - When document was added

**Business Rules:**
- Each document must belong to an employee
- Document type is flexible (not predefined enum)
- Expiry date is required and must be a future date (on creation)
- File must be uploaded to S3 before saving metadata
- Supported file types: PDF, JPG, JPEG, PNG
- Max file size: 10MB

**Document Status Logic (Calculated):**
- **Expired:** `expiry_date < TODAY`
- **Urgent (7 days):** `expiry_date BETWEEN TODAY AND TODAY + 7 DAYS`
- **Warning (30 days):** `expiry_date BETWEEN TODAY AND TODAY + 30 DAYS`
- **Valid:** `expiry_date > TODAY + 30 DAYS`

---

### **Database Relationships**

```
companies (1) ──────< employees (many)
                          │
                          │
                          └──────< documents (many)

One company has many employees
One employee has many documents
```

### **Sample Data**

```sql
-- Company
INSERT INTO companies (id, name) VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'ABC Construction Co.');

-- Employee
INSERT INTO employees (id, company_id, name, trade, mobile, joining_date) VALUES 
('660e8400-e29b-41d4-a716-446655440001', 
 '550e8400-e29b-41d4-a716-446655440000',
 'John Doe',
 'Welder',
 '+971501234567',
 '2024-01-15');

-- Documents
INSERT INTO documents (employee_id, document_type, expiry_date, file_url, file_name, file_size, file_type) VALUES 
('660e8400-e29b-41d4-a716-446655440001',
 'Visa',
 '2026-12-31',
 's3://bucket/documents/visa_john_doe.pdf',
 'visa_john_doe.pdf',
 1048576,
 'application/pdf');
```

---

## 🎯 Core Features

### **Feature 1: Dashboard**

**Purpose:** Single-page overview of entire system

**Displays:**

1. **Metric Cards** (Top Row)
   - Total Employees (count)
   - Active Documents (valid, not expired)
   - Expiring Soon (within 30 days, warning color)
   - Expired Documents (past expiry, red color)

2. **Company Summary** (If multiple companies)
   ```
   Company A: 120 employees
   Company B: 85 employees
   Total: 205 employees
   ```

3. **Critical Expiry Alerts Table**
   - Shows documents expiring in next 30 days
   - Columns: Employee Name, Company, Document Type, Expiry Date, Days Left, Actions
   - Color coding:
     - Red: Expired (already past)
     - Orange: Expiring in 7 days
     - Yellow: Expiring in 30 days
   - Sorted by: Expiry date (nearest first)
   - Filterable by: Company, Document Type

4. **Data Visualizations** (Optional/Future)
   - Pie chart: Document status distribution
   - Bar chart: Company-wise employee count
   - Line chart: Expiry trends

**Business Logic:**
```javascript
// Metric calculations
totalEmployees = COUNT(employees)
activeDocuments = COUNT(documents WHERE expiry_date > TODAY)
expiringSoon = COUNT(documents WHERE expiry_date BETWEEN TODAY AND TODAY + 30)
expired = COUNT(documents WHERE expiry_date < TODAY)
```

---

### **Feature 2: Employee Management**

#### **2.1 Employee List View**

**URL:** `/employees`

**Layout:**
- Search bar (by name)
- Filter dropdowns:
  - Company
  - Trade/Job Role
  - Document Status (All, Valid, Expiring, Expired)
- "Add New Employee" button
- Employee table/grid

**Table Columns:**
- Photo (thumbnail)
- Name
- Employee ID
- Trade
- Company
- Mobile
- Joining Date
- Document Status Badge
- Actions (Edit, Delete, View Details)

**Pagination:**
- 20 items per page (default)
- Next/Previous buttons
- Page numbers

**Business Rules:**
- Default sort: Name (A-Z)
- Can sort by any column
- Search matches partial name
- Filters can be combined

---

#### **2.2 Add/Edit Employee**

**URL:** `/employees/new` or `/employees/{id}/edit`

**Form Fields:**
1. **Name** (required, 2-100 chars)
2. **Trade** (required, dropdown + custom input)
   - Common: Welder, Electrician, Plumber, Engineer, Technician, Driver, Helper
   - Allow custom entry
3. **Company** (required, dropdown)
   - For now: Single company (auto-selected)
   - Future: Multiple companies
4. **Mobile** (required, international format with validation)
   - Format: +971501234567
   - Validation: E.164 format
5. **Joining Date** (required, date picker)
   - Cannot be future date
6. **Photo** (optional, file upload)
   - Formats: JPG, PNG
   - Max size: 2MB
   - Displays preview after upload

**Validation Rules:**
```typescript
name: required, min 2 chars, max 100 chars
trade: required, min 2 chars
company_id: required, valid UUID
mobile: required, E.164 format regex
joining_date: required, not future date
photo: optional, max 2MB, JPG/PNG only
```

**Flow:**
1. User fills form
2. Frontend validates input
3. If photo uploaded:
   - Upload to S3 first
   - Get S3 URL
4. Submit employee data with photo URL to API
5. Backend validates and saves to database
6. Redirect to employee details page
7. Show success toast: "Employee created successfully"

**Error Handling:**
- Show validation errors below each field
- If photo upload fails, show error but allow retry
- If API fails, show error toast with retry option

---

#### **2.3 Employee Details View**

**URL:** `/employees/{id}`

**Layout:**

**Section 1: Employee Info** (Top)
- Photo (large)
- Name
- Trade
- Company
- Mobile (with click-to-call icon)
- Joining Date
- Actions: Edit, Delete

**Section 2: Documents** (Main)
- "Add Document" button
- Documents grid/list
- Each document card shows:
  - Document icon (based on type)
  - Document type
  - Expiry date
  - Status badge (Expired/Expiring/Valid)
  - Days until expiry (or days since expired)
  - Actions: View File, Edit, Delete

**Business Logic:**
- Documents sorted by: Expiry date (nearest first)
- Expired documents at top with red highlight
- Click "View File" opens document in new tab (pre-signed S3 URL)

---

#### **2.4 Delete Employee**

**Flow:**
1. User clicks "Delete" button
2. Show confirmation modal:
   ```
   Are you sure you want to delete John Doe?
   This will also delete all associated documents.
   This action cannot be undone.
   [Cancel] [Delete]
   ```
3. If confirmed:
   - Call DELETE API
   - Delete employee from database
   - Cascade delete all documents (database handles this)
   - Delete document files from S3 (background job)
   - Show success toast
   - Redirect to employee list

---

### **Feature 3: Document Management**

#### **3.1 Add Document**

**URL:** Modal/Drawer from employee details page

**Form Fields:**
1. **Employee** (auto-filled, read-only in this context)
2. **Document Type** (required, custom text input)
   - Suggestions: Visa, Passport, Emirates ID, Labor Card, Bank Slip, etc.
   - User can type custom type
3. **Expiry Date** (required, date picker)
   - Must be future date (on creation)
   - Clear date display format
4. **File Upload** (required)
   - Drag & drop area
   - Or click to browse
   - Show file name and size after selection
   - Formats: PDF, JPG, PNG
   - Max size: 10MB

**Validation:**
```typescript
document_type: required, 2-100 chars
expiry_date: required, must be future date
file: required, max 10MB, PDF/JPG/PNG only
```

**Flow:**
1. User fills form and selects file
2. Frontend validates input
3. Click "Upload Document"
4. Show upload progress bar
5. Upload file to S3 (get pre-signed URL from backend)
6. Once uploaded, submit document metadata to API
7. Backend saves metadata to database
8. Show success toast
9. Document appears in employee's document list
10. If expiring soon, immediately shows in dashboard alerts

**Error Handling:**
- File too large: "File must be less than 10MB"
- Invalid format: "Only PDF, JPG, PNG files allowed"
- Upload fails: "Upload failed. Please try again."
- Network error: "Connection error. Check your internet."

---

#### **3.2 Edit Document**

**What Can Be Edited:**
- Document type (text)
- Expiry date (date)
- Replace file (file upload)

**Flow:**
1. Open edit modal with current values
2. User modifies fields
3. If file replaced:
   - Upload new file to S3
   - Delete old file from S3 (background)
   - Update file_url in database
4. Update metadata in database
5. Update `last_updated` timestamp
6. Show success toast

---

#### **3.3 Delete Document**

**Flow:**
1. User clicks "Delete" on document
2. Show confirmation:
   ```
   Delete Visa document?
   This will permanently delete the file.
   [Cancel] [Delete]
   ```
3. If confirmed:
   - Delete metadata from database
   - Queue file deletion from S3 (background job)
   - Show success toast
   - Remove from UI

---

### **Feature 4: Notification System**

#### **4.1 Email Notifications**

**When:** Daily at 9:00 AM (configurable timezone)

**Who Receives:** Admin email (configurable, can be multiple)

**Email Types:**

**A. 30-Day Warning Email**
```
Subject: Document Expiring Soon - [Employee Name]

Hi,

This is a reminder that the following document will expire in 30 days:

Employee: John Doe
Company: ABC Construction Co.
Document Type: Visa
Expiry Date: March 15, 2026
Days Remaining: 30

Please take necessary action for renewal.

[View Employee] button

---
Manpower Management System
```

**B. 7-Day Urgent Email**
```
Subject: URGENT: Document Expiring in 7 Days - [Employee Name]

Hi,

⚠️ URGENT: The following document will expire in 7 days:

Employee: John Doe
Company: ABC Construction Co.
Document Type: Passport
Expiry Date: February 17, 2026
Days Remaining: 7

IMMEDIATE ACTION REQUIRED.

[View Employee] button

---
Manpower Management System
```

**C. Expired Email**
```
Subject: EXPIRED: Document Renewal Required - [Employee Name]

Hi,

❌ The following document has EXPIRED:

Employee: John Doe
Company: ABC Construction Co.
Document Type: Emirates ID
Expiry Date: February 9, 2026
Days Overdue: 1

Please renew immediately to avoid compliance issues.

[View Employee] button

---
Manpower Management System
```

**Business Logic:**
```sql
-- Query for 30-day alerts (run daily at 9 AM)
SELECT e.name, e.company_id, d.document_type, d.expiry_date
FROM documents d
JOIN employees e ON d.employee_id = e.id
WHERE d.expiry_date = CURRENT_DATE + INTERVAL '30 days';

-- Query for 7-day alerts
WHERE d.expiry_date = CURRENT_DATE + INTERVAL '7 days';

-- Query for expired alerts (1 day old)
WHERE d.expiry_date = CURRENT_DATE - INTERVAL '1 day';
```

**Cron Job Implementation:**
1. Runs daily at configured time (9:00 AM)
2. Queries database for documents meeting criteria
3. Groups by employee if multiple documents
4. Sends email for each document
5. Logs all sent notifications
6. Retries failed emails (max 3 attempts)

---

#### **4.2 SMS Notifications (Future)**

Same triggers, shorter message format:
```
URGENT: Your Visa expires in 7 days (Feb 17).
Renew immediately. -Manpower System
```

---

### **Feature 5: File Management (S3)**

#### **Upload Flow**

```
User selects file
     ↓
Frontend validates (size, type)
     ↓
Request pre-signed URL from backend
     ↓
Backend generates S3 pre-signed URL
     ↓
Frontend uploads directly to S3
     ↓
Success → Submit metadata to backend
     ↓
Backend saves file_url to database
```

#### **Download/View Flow**

```
User clicks "View Document"
     ↓
Frontend requests file access
     ↓
Backend generates pre-signed URL (expires in 5 min)
     ↓
Return URL to frontend
     ↓
Frontend opens in new tab
     ↓
User views/downloads file
```

**S3 Bucket Structure:**
```
manpower-files/
├── employees/
│   └── photos/
│       └── {employee-id}/
│           └── profile.jpg
└── documents/
    └── {employee-id}/
        ├── visa_{timestamp}.pdf
        ├── passport_{timestamp}.pdf
        └── emirates_id_{timestamp}.jpg
```

**Security:**
- Bucket is private (no public access)
- All access via pre-signed URLs
- URLs expire after 5 minutes
- Files encrypted at rest (S3 default)

---

## 🔄 User Flows

### **Flow 1: Adding a New Employee with Documents**

```
1. User clicks "Add New Employee"
2. Fill employee form (name, trade, company, mobile, joining date)
3. Optional: Upload employee photo
4. Click "Save Employee"
5. System saves employee → Redirects to employee details page
6. User sees empty documents section
7. Click "Add Document"
8. Select document type (e.g., "Visa")
9. Select expiry date
10. Upload document file
11. Click "Upload"
12. System uploads file to S3
13. System saves document metadata
14. Document appears in list
15. Repeat steps 7-14 for other documents (Passport, ID, etc.)
16. All documents now visible on employee page
17. Dashboard automatically updates with new employee count
18. If any document expires within 30 days, appears in alerts
```

---

### **Flow 2: Daily Notification Check**

```
Every day at 9:00 AM:

1. Cron job wakes up
2. Query database for documents expiring in 30 days (exactly)
3. Query database for documents expiring in 7 days (exactly)
4. Query database for documents expired yesterday
5. For each document found:
   a. Fetch employee details
   b. Compose email with document info
   c. Send email via AWS SES
   d. Log notification sent
6. If email fails:
   a. Retry after 1 minute
   b. Retry after 5 minutes
   c. If still fails, log error
7. Send summary email to admin:
   - X emails sent successfully
   - Y emails failed
8. Job completes
```

---

### **Flow 3: Document Renewal**

```
1. User receives email: "Visa expiring in 30 days"
2. User clicks [View Employee] in email
3. Opens employee details page
4. User sees red "Expiring Soon" badge on Visa
5. User arranges visa renewal with authorities
6. Once renewed:
   a. Click "Edit" on Visa document
   b. Update expiry date to new date
   c. Upload new visa copy (replace file)
   d. Click "Save"
7. System updates document
8. Status badge changes to "Valid"
9. Document removed from dashboard alerts
10. User receives confirmation toast
```

---

### **Flow 4: Weekly Review**

```
1. User logs into system
2. Views dashboard
3. Sees metric cards:
   - 248 Total Employees
   - 212 Active Documents
   - 28 Expiring Soon ⚠️
   - 8 Expired ❌
4. Reviews "Critical Expiry Alerts" table
5. Sees list of upcoming expiries
6. Takes action on urgent items:
   - Click employee name
   - Review document
   - Contact employee/company
   - Arrange renewals
7. User can filter by company or document type
8. Export list to Excel (future feature)
```

---

## 🔌 API Endpoints

### **Employee Endpoints**

```
POST   /api/employees
GET    /api/employees
GET    /api/employees/{id}
PUT    /api/employees/{id}
DELETE /api/employees/{id}

Query Parameters for GET /api/employees:
- company_id (UUID)
- trade (string)
- search (string) - searches in name
- status (string) - all, valid, expiring, expired
- page (int) - default 1
- limit (int) - default 20, max 100
- sort_by (string) - name, joining_date
- sort_order (string) - asc, desc
```

**Example Requests:**

```bash
# Create employee
POST /api/employees
{
  "name": "John Doe",
  "trade": "Welder",
  "company_id": "550e8400-e29b-41d4-a716-446655440000",
  "mobile": "+971501234567",
  "joining_date": "2024-01-15",
  "photo_url": "https://s3.../photo.jpg"  // optional
}

# List employees with filters
GET /api/employees?company_id=550e8400...&status=expiring&page=1&limit=20

# Get employee details
GET /api/employees/660e8400-e29b-41d4-a716-446655440001

# Update employee
PUT /api/employees/660e8400-e29b-41d4-a716-446655440001
{
  "name": "John Smith",
  "mobile": "+971509876543"
}

# Delete employee
DELETE /api/employees/660e8400-e29b-41d4-a716-446655440001
```

---

### **Document Endpoints**

```
POST   /api/employees/{employee_id}/documents
GET    /api/employees/{employee_id}/documents
GET    /api/documents/{id}
PUT    /api/documents/{id}
DELETE /api/documents/{id}
GET    /api/documents/{id}/download  (returns pre-signed URL)
POST   /api/documents/upload-url      (get pre-signed URL for upload)
```

**Example Requests:**

```bash
# Get upload URL (before uploading file)
POST /api/documents/upload-url
{
  "employee_id": "660e8400...",
  "file_name": "visa.pdf",
  "file_type": "application/pdf",
  "file_size": 1048576
}
Response:
{
  "upload_url": "https://s3.amazonaws.com/...",
  "file_key": "documents/660e8400.../visa_1234567890.pdf"
}

# Create document (after file uploaded to S3)
POST /api/employees/660e8400.../documents
{
  "document_type": "Visa",
  "expiry_date": "2026-12-31",
  "file_url": "s3://bucket/documents/660e8400.../visa_1234567890.pdf",
  "file_name": "visa.pdf",
  "file_size": 1048576,
  "file_type": "application/pdf"
}

# Get download URL
GET /api/documents/770e8400.../download
Response:
{
  "download_url": "https://s3.amazonaws.com/...?signature=...",
  "expires_in": 300  // seconds
}
```

---

### **Dashboard Endpoints**

```
GET    /api/dashboard/stats
GET    /api/dashboard/expiring
GET    /api/dashboard/company-summary
```

**Example Requests:**

```bash
# Get dashboard stats
GET /api/dashboard/stats
Response:
{
  "total_employees": 248,
  "active_documents": 212,
  "expiring_soon": 28,
  "expired": 8,
  "by_company": [
    {
      "company_id": "550e8400...",
      "company_name": "ABC Construction",
      "employee_count": 120
    }
  ]
}

# Get expiring documents
GET /api/dashboard/expiring?days=30
Response:
{
  "documents": [
    {
      "document_id": "770e8400...",
      "employee_id": "660e8400...",
      "employee_name": "John Doe",
      "company_name": "ABC Construction",
      "document_type": "Visa",
      "expiry_date": "2026-03-15",
      "days_left": 30,
      "status": "warning"  // expired, urgent, warning, valid
    }
  ],
  "total": 28
}
```

---

### **Company Endpoints (Future)**

```
POST   /api/companies
GET    /api/companies
GET    /api/companies/{id}
PUT    /api/companies/{id}
DELETE /api/companies/{id}
```

---

## 📏 Business Rules & Validations

### **Employee Rules**

1. **Name:**
   - Required
   - Minimum 2 characters
   - Maximum 100 characters
   - Allowed: Letters, spaces, hyphens, apostrophes

2. **Trade:**
   - Required
   - Minimum 2 characters
   - Maximum 100 characters

3. **Mobile:**
   - Required
   - Must be E.164 format: `^\+?[1-9]\d{1,14}$`
   - Example: +971501234567

4. **Joining Date:**
   - Required
   - Cannot be future date
   - Format: YYYY-MM-DD

5. **Photo:**
   - Optional
   - Max size: 2MB
   - Formats: JPG, JPEG, PNG
   - Dimensions: No restriction (will be resized if needed)

---

### **Document Rules**

1. **Document Type:**
   - Required
   - Minimum 2 characters
   - Maximum 100 characters
   - Case-insensitive (stored as-is)

2. **Expiry Date:**
   - Required
   - Must be future date on creation
   - Can be past date on edit (for updating expired docs)
   - Format: YYYY-MM-DD

3. **File:**
   - Required
   - Max size: 10MB
   - Formats: PDF, JPG, JPEG, PNG
   - Validated on both frontend and backend

4. **Status Calculation:**
   ```javascript
   if (expiryDate < today) return "expired"
   if (expiryDate <= today + 7 days) return "urgent"
   if (expiryDate <= today + 30 days) return "warning"
   return "valid"
   ```

---

### **Notification Rules**

1. **Email Timing:**
   - Daily at 9:00 AM (configurable)
   - Timezone: Asia/Dubai (configurable)

2. **Alert Triggers:**
   - 30-day alert: Document expires exactly 30 days from today
   - 7-day alert: Document expires exactly 7 days from today
   - Expired alert: Document expired exactly 1 day ago

3. **Email Recipients:**
   - Primary: Admin email (from config)
   - CC: Additional emails (comma-separated in config)

4. **Retry Logic:**
   - Retry failed emails 3 times
   - Wait 1 minute between retries
   - Log all failures

---

## 🚀 Development Phases

### **Phase 1: Foundation (Week 1-2)**

**Backend:**
- [ ] Setup project structure
- [ ] Database schema & migrations
- [ ] Domain entities
- [ ] Repository layer (PostgreSQL)
- [ ] Employee CRUD service
- [ ] Employee HTTP handlers
- [ ] Basic error handling
- [ ] Logging setup

**Frontend:**
- [ ] Setup Next.js project
- [ ] Install dependencies (Tailwind, shadcn/ui)
- [ ] Create layout (header, sidebar)
- [ ] Setup API client
- [ ] Type definitions

**Testing:**
- [ ] Can create/read/update/delete employee
- [ ] API responds correctly
- [ ] Frontend displays data

---

### **Phase 2: Core Features (Week 2-3)**

**Backend:**
- [ ] Document service (CRUD)
- [ ] S3 integration (upload/download)
- [ ] Pre-signed URL generation
- [ ] Document HTTP handlers
- [ ] Dashboard service
- [ ] Dashboard endpoints

**Frontend:**
- [ ] Dashboard page with metrics
- [ ] Employee list page
- [ ] Employee details page
- [ ] Add/Edit employee form
- [ ] Document upload component
- [ ] Document list component

**Testing:**
- [ ] Can upload documents
- [ ] Files stored in S3
- [ ] Dashboard shows correct stats
- [ ] Can view employee documents

---

### **Phase 3: Notifications (Week 3-4)**

**Backend:**
- [ ] Notification service interface
- [ ] Email notifier (AWS SES)
- [ ] Cron job scheduler
- [ ] Expiry checker logic
- [ ] Email templates
- [ ] Notification logging

**Frontend:**
- [ ] Expiry alerts table on dashboard
- [ ] Status badges on documents
- [ ] Filter by status

**Testing:**
- [ ] Cron job runs daily
- [ ] Emails sent for expiring docs
- [ ] Correct documents identified
- [ ] Email content accurate

---

### **Phase 4: Polish & Deploy (Week 4)**

**Backend:**
- [ ] Input validation
- [ ] Error handling
- [ ] Rate limiting
- [ ] Security headers
- [ ] Documentation

**Frontend:**
- [ ] Loading states
- [ ] Error handling
- [ ] Toast notifications
- [ ] Responsive design
- [ ] Accessibility

**Deployment:**
- [ ] Backend on AWS EC2/ECS
- [ ] Frontend on Vercel
- [ ] Database on AWS RDS
- [ ] S3 bucket setup
- [ ] SES email verification
- [ ] Environment variables
- [ ] CI/CD pipeline

---

## 🎨 UI/UX Guidelines

### **Color Scheme**

```css
/* Status Colors */
--success: #10b981     /* Green - Valid */
--warning: #f59e0b     /* Yellow - Expiring 30 days */
--urgent: #f97316      /* Orange - Expiring 7 days */
--danger: #ef4444      /* Red - Expired */

/* Neutral */
--background: #ffffff
--card: #f9fafb
--border: #e5e7eb
--text: #111827
--muted: #6b7280
```

### **Typography**

- **Headings:** Bold, clear hierarchy
- **Body:** 14-16px, readable
- **Tables:** Monospace for IDs, regular for text

### **Components**

- **Cards:** Rounded corners, subtle shadows
- **Buttons:** Clear primary/secondary distinction
- **Forms:** Inline validation, clear errors
- **Tables:** Zebra striping, hover effects
- **Status Badges:** Color-coded, rounded pills

### **Responsive Design**

- **Desktop (1024px+):** Sidebar + main content
- **Tablet (768-1023px):** Collapsible sidebar
- **Mobile (<768px):** Bottom nav, cards instead of table

---

## 🔒 Security Considerations

1. **Authentication (Future):**
   - JWT-based auth
   - Secure password hashing (bcrypt)
   - Session management

2. **API Security:**
   - HTTPS only
   - CORS configured
   - Rate limiting
   - Input validation
   - SQL injection prevention (parameterized queries)

3. **File Security:**
   - Private S3 bucket
   - Pre-signed URLs (5-min expiry)
   - File type validation
   - Size validation
   - Virus scanning (future)

4. **Data Security:**
   - Encrypted at rest (database encryption)
   - Encrypted in transit (HTTPS)
   - Backup strategy
   - Audit logs

---

## 📝 Notes for AI Editor

When generating code for this project:

1. **Always refer to database schema** when creating queries
2. **Follow Clean Architecture** layers (domain → repository → service → handler)
3. **Handle ALL errors** - never ignore errors in Go
4. **Validate input** on both frontend and backend
5. **Use TypeScript** - no `any` types
6. **Mobile-first** responsive design
7. **Accessible** components (ARIA labels)
8. **Logging** for important operations
9. **Comments** for complex business logic
10. **Tests** for critical functionality

---

## 🎯 Success Criteria

The project is successful when:

✅ Admin can add/edit/delete employees  
✅ Admin can upload documents with expiry dates  
✅ Dashboard shows accurate counts and alerts  
✅ Email alerts sent daily at 9 AM  
✅ Files securely stored in S3  
✅ System handles 500+ employees smoothly  
✅ Mobile-responsive interface  
✅ No critical bugs in production  

---

**Document Version:** 1.0  
**Last Updated:** February 2026  
**Status:** Final - Ready for Development
