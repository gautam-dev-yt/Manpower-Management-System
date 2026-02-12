# 🎯 Manpower Management System - Development Guidelines

> **AI-Assisted Development Ready**  
> This document provides comprehensive guidelines for AI editors (Cursor, Windsurf, GitHub Copilot) to generate clean, scalable, and maintainable code.

---

## 📋 Table of Contents

1. [Tech Stack & Architecture](#tech-stack--architecture)
2. [Project Structure](#project-structure)
3. [Go Backend Guidelines](#go-backend-guidelines)
4. [Next.js Frontend Guidelines](#nextjs-frontend-guidelines)
5. [Code Quality Standards](#code-quality-standards)
6. [Setup Instructions](#setup-instructions)
7. [AI Editor Prompts](#ai-editor-prompts)

---

## 🏗️ Tech Stack & Architecture

### **Backend Stack (Go)**

| Component | Technology | Justification |
|-----------|-----------|---------------|
| **Language** | Go 1.22+ | Simplicity, performance, excellent concurrency, strong typing |
| **Framework** | Chi Router | Minimalist, idiomatic Go, stdlib-compatible, middleware support |
| **Database** | PostgreSQL 15+ | Robust ACID compliance, JSON support (JSONB), excellent date handling |
| **ORM/Query** | sqlc | Type-safe SQL, no reflection overhead, learn SQL properly, performance |
| **Migration** | golang-migrate | Industry standard, version control for schema |
| **Validation** | go-playground/validator | Struct tag validation, comprehensive rules |
| **Config** | godotenv + viper | Environment management, 12-factor app principles |
| **File Storage** | AWS S3 SDK | Secure, scalable object storage, pre-signed URLs |
| **Email** | AWS SES | Reliable, cost-effective, 62k emails/month free |
| **SMS (Future)** | AWS SNS | Scalable notification service |
| **Scheduler** | gocron | Simple cron job management, no external dependencies |
| **Logging** | zerolog | Fast, structured logging, zero allocation |
| **Testing** | testify | Assertions, mocking, test suites |

### **Frontend Stack (Next.js)**

| Component | Technology | Justification |
|-----------|-----------|---------------|
| **Framework** | Next.js 15 (App Router) | SSR, file-based routing, React Server Components |
| **Language** | TypeScript 5+ | Type safety, better DX, catch errors early |
| **Styling** | Tailwind CSS 3 | Utility-first, responsive, consistent design system |
| **UI Components** | shadcn/ui | Accessible, customizable, Radix UI primitives |
| **State Management** | Zustand (if needed) | Minimal, no boilerplate, simple API |
| **Data Fetching** | TanStack Query | Caching, automatic refetching, optimistic updates |
| **Forms** | React Hook Form | Minimal re-renders, validation, great DX |
| **Tables** | TanStack Table | Headless, powerful, responsive |
| **Date Handling** | date-fns | Lightweight, modular, tree-shakeable |
| **HTTP Client** | Fetch API | Native, no dependencies |
| **Validation** | Zod | Type-safe schema validation, works with RHF |

### **DevOps & Deployment**

| Component | Technology | Justification |
|-----------|-----------|---------------|
| **Hosting** | AWS (EC2/ECS) | Flexible, free tier eligible, full control |
| **Frontend** | Vercel | Optimized for Next.js, automatic deployments |
| **Database** | AWS RDS (PostgreSQL) | Managed, automated backups, multi-AZ |
| **CI/CD** | GitHub Actions | Free, integrated, powerful |
| **Containerization** | Docker | Consistent environments, easy deployment |

### **Architecture Pattern**

**Clean/Hexagonal Architecture** with the following layers:

```
┌─────────────────────────────────────┐
│   Delivery Layer (API Handlers)    │  ← HTTP/JSON interface
├─────────────────────────────────────┤
│   Service Layer (Business Logic)   │  ← Use cases, orchestration
├─────────────────────────────────────┤
│   Repository Layer (Data Access)   │  ← Database operations
├─────────────────────────────────────┤
│   Domain Layer (Entities/Models)   │  ← Pure business objects
└─────────────────────────────────────┘
```

**Dependency Rule:** Dependencies flow inward (Delivery → Service → Repository → Domain)

---

## 📁 Project Structure

### **Backend (Go) - Recommended Structure**

```
manpower-backend/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
│
├── internal/                          # Private application code
│   ├── domain/                        # Business entities (no dependencies)
│   │   ├── employee.go
│   │   ├── document.go
│   │   └── company.go
│   │
│   ├── repository/                    # Data access layer (interfaces + implementations)
│   │   ├── interfaces.go              # Repository interfaces
│   │   ├── employee_postgres.go
│   │   ├── document_postgres.go
│   │   └── company_postgres.go
│   │
│   ├── service/                       # Business logic layer
│   │   ├── employee_service.go
│   │   ├── document_service.go
│   │   ├── notification_service.go    # Interface-based
│   │   └── dashboard_service.go
│   │
│   ├── handler/                       # HTTP handlers (delivery layer)
│   │   ├── employee_handler.go
│   │   ├── document_handler.go
│   │   ├── dashboard_handler.go
│   │   └── response.go                # Standard responses
│   │
│   ├── middleware/                    # HTTP middleware
│   │   ├── logger.go
│   │   ├── cors.go
│   │   ├── rate_limit.go
│   │   └── error.go
│   │
│   ├── notification/                  # Notification implementations
│   │   ├── email_notifier.go          # AWS SES implementation
│   │   ├── sms_notifier.go            # AWS SNS implementation
│   │   └── multi_notifier.go          # Composite pattern
│   │
│   ├── storage/                       # File storage
│   │   └── s3_storage.go
│   │
│   ├── scheduler/                     # Cron jobs
│   │   └── expiry_checker.go
│   │
│   └── pkg/                           # Shared utilities (optional)
│       ├── validator/
│       ├── logger/
│       └── errors/
│
├── migrations/                        # Database migrations
│   ├── 001_create_companies.up.sql
│   ├── 001_create_companies.down.sql
│   ├── 002_create_employees.up.sql
│   └── ...
│
├── config/                            # Configuration files
│   ├── config.go
│   └── database.go
│
├── scripts/                           # Helper scripts
│   ├── seed.go
│   └── migration.sh
│
├── docker/                            # Docker files
│   └── Dockerfile
│
├── .env.example                       # Environment template
├── .gitignore
├── go.mod
├── go.sum
├── Makefile                           # Common commands
└── README.md
```

### **Frontend (Next.js) - Recommended Structure**

```
manpower-frontend/
├── src/
│   ├── app/                           # App Router (Next.js 15)
│   │   ├── layout.tsx                 # Root layout
│   │   ├── page.tsx                   # Home/Dashboard
│   │   ├── loading.tsx                # Loading states
│   │   ├── error.tsx                  # Error boundaries
│   │   │
│   │   ├── dashboard/                 # Dashboard routes
│   │   │   └── page.tsx
│   │   │
│   │   ├── employees/                 # Employee routes
│   │   │   ├── page.tsx               # List view
│   │   │   ├── [id]/
│   │   │   │   └── page.tsx           # Detail view
│   │   │   └── new/
│   │   │       └── page.tsx           # Create view
│   │   │
│   │   └── api/                       # API routes (if needed)
│   │       └── health/
│   │           └── route.ts
│   │
│   ├── components/                    # React components
│   │   ├── ui/                        # shadcn/ui components
│   │   │   ├── button.tsx
│   │   │   ├── input.tsx
│   │   │   ├── card.tsx
│   │   │   └── ...
│   │   │
│   │   ├── layout/                    # Layout components
│   │   │   ├── header.tsx
│   │   │   ├── sidebar.tsx
│   │   │   └── footer.tsx
│   │   │
│   │   ├── dashboard/                 # Dashboard-specific
│   │   │   ├── metric-card.tsx
│   │   │   ├── expiry-alerts.tsx
│   │   │   └── company-summary.tsx
│   │   │
│   │   ├── employees/                 # Employee-specific
│   │   │   ├── employee-table.tsx
│   │   │   ├── employee-form.tsx
│   │   │   └── employee-filters.tsx
│   │   │
│   │   └── documents/                 # Document-specific
│   │       ├── document-list.tsx
│   │       ├── document-upload.tsx
│   │       └── document-card.tsx
│   │
│   ├── lib/                           # Utility libraries
│   │   ├── api.ts                     # API client
│   │   ├── utils.ts                   # Helper functions
│   │   └── constants.ts               # Constants
│   │
│   ├── hooks/                         # Custom React hooks
│   │   ├── use-employees.ts
│   │   ├── use-documents.ts
│   │   └── use-dashboard.ts
│   │
│   ├── types/                         # TypeScript types
│   │   ├── employee.ts
│   │   ├── document.ts
│   │   └── api.ts
│   │
│   ├── store/                         # State management (if needed)
│   │   └── auth-store.ts
│   │
│   └── styles/                        # Global styles
│       └── globals.css
│
├── public/                            # Static assets
│   ├── images/
│   └── icons/
│
├── .env.local.example                 # Environment template
├── .gitignore
├── .eslintrc.json                     # ESLint config
├── .prettierrc                        # Prettier config
├── next.config.mjs                    # Next.js config
├── tailwind.config.ts                 # Tailwind config
├── tsconfig.json                      # TypeScript config
├── package.json
└── README.md
```

---

## 🐹 Go Backend Guidelines

### **Critical Principles**

1. **Keep It Simple** - Start simple, add complexity only when needed
2. **Explicit Over Implicit** - Clear dependencies, no magic
3. **Interfaces at Boundaries** - Define interfaces where they're consumed, not where they're implemented
4. **Error Handling** - Always handle errors, never ignore them
5. **No Framework Lock-in** - Business logic independent of frameworks

### **Coding Standards**

#### **1. Project Layout**

- Use `internal/` for private code that shouldn't be imported by other projects
- Don't create packages just to organize files - create them for logical separation
- Avoid deep nesting - flat is better than nested
- No `pkg/`, `utils/`, `helpers/`, `models/` packages (antipattern)
- Group by feature/domain, not by layer

#### **2. Package Naming**

```go
// ✅ GOOD - lowercase, singular, descriptive
package employee
package document
package notification

// ❌ BAD - plural, mixed case, vague
package employees
package employeeModels
package utils
```

#### **3. Struct and Interface Naming**

```go
// ✅ GOOD - clear, concise
type Employee struct { ... }
type EmployeeRepository interface { ... }
type EmployeeService struct { ... }

// ❌ BAD - redundant prefixes
type EmployeeStruct struct { ... }
type IEmployeeRepository interface { ... }
```

#### **4. Error Handling**

```go
// ✅ GOOD - explicit error handling
employee, err := repo.FindByID(ctx, id)
if err != nil {
    return fmt.Errorf("failed to find employee: %w", err)
}

// ❌ BAD - ignoring errors
employee, _ := repo.FindByID(ctx, id)

// ✅ GOOD - early returns, less nesting
func GetEmployee(ctx context.Context, id string) (*Employee, error) {
    if id == "" {
        return nil, ErrInvalidID
    }
    
    employee, err := repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("repository error: %w", err)
    }
    
    return employee, nil
}
```

#### **5. Context Usage**

```go
// ✅ GOOD - context as first parameter
func (r *EmployeeRepository) FindByID(ctx context.Context, id string) (*Employee, error) {
    // Always pass context through
    return r.db.GetContext(ctx, &employee, query, id)
}

// ❌ BAD - no context or wrong position
func (r *EmployeeRepository) FindByID(id string, ctx context.Context) (*Employee, error) {}
```

#### **6. Repository Pattern**

```go
// ✅ GOOD - interface in consumer package
// internal/service/employee_service.go
package service

// Define interface where it's used
type EmployeeRepository interface {
    Create(ctx context.Context, employee *domain.Employee) error
    FindByID(ctx context.Context, id string) (*domain.Employee, error)
    Update(ctx context.Context, employee *domain.Employee) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filters *domain.EmployeeFilters) ([]*domain.Employee, error)
}

type EmployeeService struct {
    repo EmployeeRepository
    logger *zerolog.Logger
}

// ✅ Implementation in infrastructure layer
// internal/repository/employee_postgres.go
package repository

type PostgresEmployeeRepository struct {
    db *sqlx.DB
}

// Implements service.EmployeeRepository interface
func (r *PostgresEmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
    // Implementation
}
```

#### **7. Dependency Injection**

```go
// ✅ GOOD - constructor pattern
func NewEmployeeService(repo repository.EmployeeRepository, logger *zerolog.Logger) *EmployeeService {
    return &EmployeeService{
        repo:   repo,
        logger: logger,
    }
}

// Usage in main.go
func main() {
    db := setupDatabase()
    logger := setupLogger()
    
    // Inject dependencies
    empRepo := repository.NewPostgresEmployeeRepository(db)
    empService := service.NewEmployeeService(empRepo, logger)
    empHandler := handler.NewEmployeeHandler(empService)
    
    // Setup routes
    r := chi.NewRouter()
    empHandler.RegisterRoutes(r)
}
```

#### **8. Struct Tags and Validation**

```go
type Employee struct {
    ID          string    `json:"id" db:"id"`
    Name        string    `json:"name" db:"name" validate:"required,min=2,max=100"`
    Trade       string    `json:"trade" db:"trade" validate:"required"`
    CompanyID   string    `json:"company_id" db:"company_id" validate:"required,uuid"`
    Mobile      string    `json:"mobile" db:"mobile" validate:"required,e164"`
    JoiningDate time.Time `json:"joining_date" db:"joining_date" validate:"required"`
    PhotoURL    *string   `json:"photo_url,omitempty" db:"photo_url"` // Nullable
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

#### **9. Clean Architecture Layers**

**Domain Layer** (No dependencies)
```go
// internal/domain/employee.go
package domain

// Pure business entity - no external dependencies
type Employee struct {
    ID          string
    Name        string
    Trade       string
    CompanyID   string
    Mobile      string
    JoiningDate time.Time
    PhotoURL    *string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Business validation logic
func (e *Employee) Validate() error {
    if e.Name == "" {
        return ErrInvalidName
    }
    if e.Trade == "" {
        return ErrInvalidTrade
    }
    return nil
}
```

**Repository Layer** (Data access)
```go
// internal/repository/employee_postgres.go
package repository

type PostgresEmployeeRepository struct {
    db *sqlx.DB
}

func (r *PostgresEmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
    query := `
        INSERT INTO employees (id, name, trade, company_id, mobile, joining_date, photo_url)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    
    _, err := r.db.ExecContext(ctx, query,
        employee.ID,
        employee.Name,
        employee.Trade,
        employee.CompanyID,
        employee.Mobile,
        employee.JoiningDate,
        employee.PhotoURL,
    )
    
    if err != nil {
        return fmt.Errorf("failed to create employee: %w", err)
    }
    
    return nil
}
```

**Service Layer** (Business logic)
```go
// internal/service/employee_service.go
package service

type EmployeeService struct {
    repo   EmployeeRepository
    storage StorageService
    logger *zerolog.Logger
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, req *CreateEmployeeRequest) (*domain.Employee, error) {
    // Validation
    if err := req.Validate(); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    // Business logic
    employee := &domain.Employee{
        ID:          uuid.New().String(),
        Name:        req.Name,
        Trade:       req.Trade,
        CompanyID:   req.CompanyID,
        Mobile:      req.Mobile,
        JoiningDate: req.JoiningDate,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    // Upload photo if provided
    if req.Photo != nil {
        photoURL, err := s.storage.Upload(ctx, req.Photo)
        if err != nil {
            return nil, fmt.Errorf("failed to upload photo: %w", err)
        }
        employee.PhotoURL = &photoURL
    }
    
    // Persist
    if err := s.repo.Create(ctx, employee); err != nil {
        return nil, fmt.Errorf("repository error: %w", err)
    }
    
    s.logger.Info().Str("employee_id", employee.ID).Msg("employee created")
    
    return employee, nil
}
```

**Handler Layer** (HTTP delivery)
```go
// internal/handler/employee_handler.go
package handler

type EmployeeHandler struct {
    service *service.EmployeeService
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    var req CreateEmployeeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }
    
    employee, err := h.service.CreateEmployee(ctx, &req)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err.Error())
        return
    }
    
    respondJSON(w, http.StatusCreated, employee)
}

func (h *EmployeeHandler) RegisterRoutes(r chi.Router) {
    r.Route("/api/employees", func(r chi.Router) {
        r.Post("/", h.CreateEmployee)
        r.Get("/", h.ListEmployees)
        r.Get("/{id}", h.GetEmployee)
        r.Put("/{id}", h.UpdateEmployee)
        r.Delete("/{id}", h.DeleteEmployee)
    })
}
```

#### **10. Testing**

```go
// Unit test example
func TestEmployeeService_CreateEmployee(t *testing.T) {
    // Setup
    mockRepo := &MockEmployeeRepository{}
    mockStorage := &MockStorageService{}
    logger := zerolog.Nop()
    
    service := NewEmployeeService(mockRepo, mockStorage, &logger)
    
    // Test data
    req := &CreateEmployeeRequest{
        Name:        "John Doe",
        Trade:       "Engineer",
        CompanyID:   "company-123",
        Mobile:      "+1234567890",
        JoiningDate: time.Now(),
    }
    
    // Mock expectations
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    // Execute
    employee, err := service.CreateEmployee(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, employee)
    assert.Equal(t, req.Name, employee.Name)
    
    mockRepo.AssertExpectations(t)
}
```

### **Code Formatting**

Always run before committing:
```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .

# Run linter
golangci-lint run

# Run tests
go test ./... -v
```

---

## ⚛️ Next.js Frontend Guidelines

### **Critical Principles**

1. **TypeScript Everywhere** - No `any` types unless absolutely necessary
2. **Component Composition** - Small, focused components
3. **Server Components by Default** - Use Client Components only when needed
4. **Avoid Prop Drilling** - Use composition, not deep prop passing
5. **Responsive Mobile-First** - Design for mobile, enhance for desktop

### **Coding Standards**

#### **1. File and Component Naming**

```typescript
// ✅ GOOD - PascalCase for components
// components/employees/EmployeeCard.tsx
export function EmployeeCard({ employee }: { employee: Employee }) { ... }

// ✅ GOOD - kebab-case for non-component files
// lib/api-client.ts
// hooks/use-employees.ts

// ❌ BAD - inconsistent naming
// components/employeecard.tsx
// lib/apiClient.ts
```

#### **2. Component Structure**

```typescript
// ✅ GOOD - clear, organized structure
'use client'; // Only if client component

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import type { Employee } from '@/types/employee';

interface EmployeeCardProps {
  employee: Employee;
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
}

export function EmployeeCard({ employee, onEdit, onDelete }: EmployeeCardProps) {
  const [isLoading, setIsLoading] = useState(false);
  
  // Event handlers
  const handleEdit = () => {
    onEdit(employee.id);
  };
  
  // Render helpers (if complex)
  const renderStatus = () => {
    // ...
  };
  
  // Main render
  return (
    <div className="rounded-lg border p-4">
      <h3 className="text-lg font-semibold">{employee.name}</h3>
      <p className="text-sm text-muted-foreground">{employee.trade}</p>
      {renderStatus()}
      <Button onClick={handleEdit} disabled={isLoading}>
        Edit
      </Button>
    </div>
  );
}
```

#### **3. TypeScript Types**

```typescript
// ✅ GOOD - explicit, well-defined types
// types/employee.ts
export interface Employee {
  id: string;
  name: string;
  trade: string;
  companyId: string;
  mobile: string;
  joiningDate: string; // ISO date string
  photoUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateEmployeeRequest {
  name: string;
  trade: string;
  companyId: string;
  mobile: string;
  joiningDate: string;
  photo?: File;
}

export interface EmployeeFilters {
  companyId?: string;
  trade?: string;
  status?: 'all' | 'valid' | 'expiring' | 'expired';
  search?: string;
}

// ❌ BAD - any types
export interface Employee {
  id: any;
  data: any;
}
```

#### **4. API Client**

```typescript
// lib/api.ts
import { env } from '@/lib/env';

const API_BASE_URL = env.NEXT_PUBLIC_API_URL;

class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public data?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

async function fetcher<T>(
  endpoint: string,
  options?: RequestInit
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  });
  
  if (!response.ok) {
    const data = await response.json().catch(() => null);
    throw new ApiError(
      data?.message || 'An error occurred',
      response.status,
      data
    );
  }
  
  return response.json();
}

// API methods
export const api = {
  employees: {
    list: (filters?: EmployeeFilters) => 
      fetcher<Employee[]>('/api/employees', {
        method: 'GET',
        // Add query params if needed
      }),
    
    get: (id: string) =>
      fetcher<Employee>(`/api/employees/${id}`),
    
    create: (data: CreateEmployeeRequest) =>
      fetcher<Employee>('/api/employees', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    
    update: (id: string, data: Partial<Employee>) =>
      fetcher<Employee>(`/api/employees/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    
    delete: (id: string) =>
      fetcher<void>(`/api/employees/${id}`, {
        method: 'DELETE',
      }),
  },
  
  dashboard: {
    stats: () =>
      fetcher<DashboardStats>('/api/dashboard/stats'),
  },
};
```

#### **5. Custom Hooks**

```typescript
// hooks/use-employees.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { Employee, EmployeeFilters } from '@/types/employee';

export function useEmployees(filters?: EmployeeFilters) {
  return useQuery({
    queryKey: ['employees', filters],
    queryFn: () => api.employees.list(filters),
  });
}

export function useEmployee(id: string) {
  return useQuery({
    queryKey: ['employees', id],
    queryFn: () => api.employees.get(id),
    enabled: !!id,
  });
}

export function useCreateEmployee() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: api.employees.create,
    onSuccess: () => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ['employees'] });
    },
  });
}

export function useUpdateEmployee() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Employee> }) =>
      api.employees.update(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: ['employees', id] });
      queryClient.invalidateQueries({ queryKey: ['employees'] });
    },
  });
}
```

#### **6. Form Handling**

```typescript
// components/employees/EmployeeForm.tsx
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Form, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';

const employeeSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters').max(100),
  trade: z.string().min(1, 'Trade is required'),
  companyId: z.string().uuid('Invalid company ID'),
  mobile: z.string().regex(/^\+?[1-9]\d{1,14}$/, 'Invalid phone number'),
  joiningDate: z.string().min(1, 'Joining date is required'),
  photo: z.instanceof(File).optional(),
});

type EmployeeFormData = z.infer<typeof employeeSchema>;

interface EmployeeFormProps {
  initialData?: Partial<Employee>;
  onSubmit: (data: EmployeeFormData) => void;
  isLoading?: boolean;
}

export function EmployeeForm({ initialData, onSubmit, isLoading }: EmployeeFormProps) {
  const form = useForm<EmployeeFormData>({
    resolver: zodResolver(employeeSchema),
    defaultValues: {
      name: initialData?.name || '',
      trade: initialData?.trade || '',
      companyId: initialData?.companyId || '',
      mobile: initialData?.mobile || '',
      joiningDate: initialData?.joiningDate || '',
    },
  });
  
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <Input {...field} placeholder="Enter name" />
              <FormMessage />
            </FormItem>
          )}
        />
        
        {/* Other fields... */}
        
        <Button type="submit" disabled={isLoading}>
          {isLoading ? 'Saving...' : 'Save Employee'}
        </Button>
      </form>
    </Form>
  );
}
```

#### **7. Server Components vs Client Components**

```typescript
// ✅ GOOD - Server Component (default)
// app/employees/page.tsx
import { api } from '@/lib/api';
import { EmployeeTable } from '@/components/employees/EmployeeTable';

export default async function EmployeesPage() {
  const employees = await api.employees.list();
  
  return (
    <div>
      <h1>Employees</h1>
      <EmployeeTable employees={employees} />
    </div>
  );
}

// ✅ GOOD - Client Component (only when needed)
// components/employees/EmployeeTable.tsx
'use client';

import { useState } from 'react';
import { useEmployees } from '@/hooks/use-employees';

export function EmployeeTable({ employees: initialEmployees }: { employees: Employee[] }) {
  const [filters, setFilters] = useState<EmployeeFilters>({});
  const { data: employees = initialEmployees } = useEmployees(filters);
  
  return (
    // ... interactive table with filters
  );
}
```

#### **8. Tailwind CSS Best Practices**

```typescript
// ✅ GOOD - utility classes, no hardcoding
<div className="rounded-lg border bg-card p-4 shadow-sm">
  <h3 className="text-lg font-semibold text-card-foreground">
    {employee.name}
  </h3>
  <p className="mt-1 text-sm text-muted-foreground">
    {employee.trade}
  </p>
</div>

// ✅ GOOD - conditional classes using cn helper
import { cn } from '@/lib/utils';

<div className={cn(
  "rounded-lg border p-4",
  status === 'expired' && "border-destructive bg-destructive/10",
  status === 'expiring' && "border-warning bg-warning/10",
  status === 'valid' && "border-success bg-success/10"
)}>
  {/* ... */}
</div>

// ❌ BAD - inline styles, hardcoded colors
<div style={{ backgroundColor: '#ff0000', padding: '16px' }}>
  {/* ... */}
</div>
```

#### **9. Error Handling**

```typescript
// ✅ GOOD - comprehensive error handling
'use client';

import { useEffect } from 'react';
import { Button } from '@/components/ui/button';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    // Log to error reporting service
    console.error('Error:', error);
  }, [error]);
  
  return (
    <div className="flex min-h-[400px] flex-col items-center justify-center">
      <h2 className="text-xl font-semibold">Something went wrong!</h2>
      <p className="mt-2 text-sm text-muted-foreground">
        {error.message}
      </p>
      <Button onClick={reset} className="mt-4">
        Try again
      </Button>
    </div>
  );
}
```

### **Folder Organization Rules**

1. **Colocation** - Keep related files close
2. **Feature Folders** - Group by feature, not type
3. **Index Files** - Use barrel exports sparingly (only for clean APIs)
4. **Shared Components** - Generic components in `components/ui/`
5. **Feature Components** - Domain-specific in `components/[feature]/`

---

## ✅ Code Quality Standards

### **Go Checklist**

- [ ] All code formatted with `gofmt` and `goimports`
- [ ] No unused imports or variables
- [ ] All errors handled (no `_` for errors)
- [ ] Context passed as first parameter
- [ ] Interfaces defined where consumed
- [ ] Repository pattern for data access
- [ ] Dependency injection via constructors
- [ ] Unit tests for business logic (>70% coverage)
- [ ] Integration tests for critical paths
- [ ] Meaningful variable names (no single letters except `i`, `j` in loops)
- [ ] Comments on exported functions
- [ ] No global variables (except constants)

### **TypeScript/React Checklist**

- [ ] All files use TypeScript (no `.js`)
- [ ] No `any` types (use `unknown` if needed)
- [ ] All components typed (props, state)
- [ ] Server Components by default
- [ ] Client Components only when necessary
- [ ] Forms validated with Zod schemas
- [ ] API calls use TanStack Query
- [ ] Proper error boundaries
- [ ] Accessible components (ARIA labels)
- [ ] Responsive design (mobile-first)
- [ ] ESLint and Prettier configured
- [ ] No console.logs in production

### **Git Commit Convention**

```bash
# Format: <type>(<scope>): <subject>

# Examples:
feat(employee): add create employee endpoint
fix(document): correct expiry date calculation
refactor(repository): simplify query logic
docs(readme): update setup instructions
test(service): add employee service tests
chore(deps): update dependencies
```

---

## 🚀 Setup Instructions

### **Prerequisites**

```bash
# Go 1.22+
go version

# Node.js 20+
node --version

# PostgreSQL 15+
psql --version

# Docker (optional)
docker --version

# AWS CLI (for deployment)
aws --version
```

### **Backend Setup**

```bash
# 1. Clone repository
git clone <repo-url>
cd manpower-backend

# 2. Copy environment file
cp .env.example .env
# Edit .env with your values

# 3. Install dependencies
go mod download

# 4. Install tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 5. Setup database
createdb manpower_dev
make migrate-up

# 6. Run application
make run
# OR
go run cmd/api/main.go

# 7. Run tests
make test
```

### **Frontend Setup**

```bash
# 1. Navigate to frontend
cd manpower-frontend

# 2. Install dependencies
npm install

# 3. Copy environment file
cp .env.local.example .env.local
# Edit .env.local with your values

# 4. Run development server
npm run dev

# 5. Open browser
# http://localhost:3000

# 6. Build for production
npm run build
npm start
```

### **Docker Setup (Optional)**

```bash
# 1. Build and run with Docker Compose
docker-compose up -d

# 2. Stop services
docker-compose down

# 3. View logs
docker-compose logs -f
```

### **Environment Variables**

**Backend (.env)**
```env
# Database
DATABASE_URL=postgres://user:password@localhost:5432/manpower_dev?sslmode=disable

# Server
PORT=8080
ENV=development

# AWS
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-key
AWS_SECRET_ACCESS_KEY=your-secret
S3_BUCKET=manpower-files
SES_FROM_EMAIL=alerts@example.com

# Notification
NOTIFICATION_EMAIL_ENABLED=true
NOTIFICATION_SMS_ENABLED=false

# Cron
EXPIRY_CHECK_TIME=09:00 # 9 AM daily
```

**Frontend (.env.local)**
```env
# API
NEXT_PUBLIC_API_URL=http://localhost:8080

# Optional
NEXT_PUBLIC_APP_URL=http://localhost:3000
```

---

## 🤖 AI Editor Prompts

### **General Guidelines for AI**

When generating code for this project, follow these rules:

1. **Go Code:**
   - Use Clean/Hexagonal Architecture
   - Define interfaces where they're consumed
   - Handle all errors explicitly
   - Use sqlc for database queries
   - Follow repository pattern
   - Context as first parameter
   - Dependency injection via constructors

2. **TypeScript/React Code:**
   - Use TypeScript, no `any` types
   - Server Components by default
   - Zod schemas for validation
   - TanStack Query for data fetching
   - shadcn/ui for components
   - Tailwind CSS for styling
   - Mobile-first responsive design

3. **Testing:**
   - Unit tests for business logic
   - Integration tests for critical paths
   - Use table-driven tests in Go
   - Mock external dependencies

### **Example Prompts for AI Editor**

#### **Creating a New API Endpoint**

```
Create a new Go API endpoint for listing employees with filtering.

Requirements:
- Domain layer: Employee struct in internal/domain/
- Repository layer: EmployeeRepository interface and PostgresEmployeeRepository implementation
- Service layer: EmployeeService with ListEmployees method
- Handler layer: HTTP handler for GET /api/employees
- Support filters: company_id, trade, status, search
- Support pagination: page, limit
- Return JSON response
- Follow Clean Architecture
- Use Chi router
- Handle errors properly
- Add logging with zerolog
```

#### **Creating a Frontend Component**

```
Create a TypeScript React component for displaying employee cards.

Requirements:
- Component name: EmployeeCard
- Props: employee (Employee type), onEdit, onDelete callbacks
- Use shadcn/ui Card component
- Display: name, trade, company, joining date, status badge
- Status badge colors: green (valid), yellow (expiring), red (expired)
- Action buttons: Edit, Delete with confirmation
- Responsive design (mobile-first)
- Accessible (ARIA labels)
- TypeScript with full typing
- Use Tailwind CSS for styling
```

#### **Creating a Database Migration**

```
Create a PostgreSQL migration for the employees table.

Requirements:
- Table name: employees
- Columns: id (uuid primary key), name (varchar), trade (varchar), company_id (uuid foreign key), mobile (varchar), joining_date (date), photo_url (varchar nullable), created_at (timestamp), updated_at (timestamp)
- Indexes: company_id, name (for search)
- Foreign key to companies table
- Up and down migrations
- Use golang-migrate format
```

---

## 📚 Reference Documentation

### **Go Resources**

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture in Go](https://threedots.tech/post/introducing-clean-architecture/)

### **TypeScript/React Resources**

- [Next.js Documentation](https://nextjs.org/docs)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [React Documentation](https://react.dev)
- [TanStack Query](https://tanstack.com/query/latest)
- [shadcn/ui](https://ui.shadcn.com/)
- [Tailwind CSS](https://tailwindcss.com/docs)

### **Architecture Resources**

- [Clean Architecture (Uncle Bob)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Repository Pattern](https://threedots.tech/post/repository-pattern-in-go/)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)

---

## 🎯 Quick Reference

### **Go Commands**

```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .

# Run linter
golangci-lint run

# Run tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Generate sqlc
sqlc generate

# Run migrations
migrate -path migrations -database "postgres://..." up

# Build
go build -o bin/api cmd/api/main.go

# Run
./bin/api
```

### **Frontend Commands**

```bash
# Development
npm run dev

# Build
npm run build

# Lint
npm run lint

# Format
npm run format

# Type check
npm run type-check

# Test
npm run test
```

---

## ✨ Final Notes

**Key Principles:**
1. **Simplicity First** - Start simple, add complexity when needed
2. **Explicit Over Implicit** - Clear dependencies, no magic
3. **Type Safety** - Use Go's and TypeScript's type systems fully
4. **Testability** - Design for testing from the start
5. **Maintainability** - Code is read more than written
6. **Scalability** - Design for growth, but don't over-engineer

**Remember:**
- This is a learning project (Go backend)
- Keep frontend simple (AI-generated, focus on backend)
- Follow standards, but pragmatism over dogma
- Test critical paths, don't aim for 100% coverage
- Document decisions, especially deviations from standards

**When in Doubt:**
- Check official docs first
- Look at standard library examples
- Ask "Is this the simplest solution?"
- Write tests to clarify behavior
- Use AI editor to generate boilerplate, review critically

---

**Document Version:** 1.0  
**Last Updated:** February 2026  
**Maintainer:** Development Team
