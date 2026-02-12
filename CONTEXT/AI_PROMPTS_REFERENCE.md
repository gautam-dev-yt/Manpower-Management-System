# 🤖 AI Editor Prompts - Quick Reference

**For Cursor, Windsurf, GitHub Copilot**

Use these prompts to generate high-quality code that follows project standards.

---

## 📝 General Instructions for AI

Before generating code, tell your AI editor:

```
This is a Clean Architecture project using:
- Backend: Go 1.22+ with Chi router, PostgreSQL, sqlc
- Frontend: Next.js 15 (App Router), TypeScript, Tailwind CSS, shadcn/ui
- Follow the guidelines in DEVELOPMENT_GUIDELINES.md
- Use Hexagonal Architecture pattern
- All errors must be handled
- TypeScript with no 'any' types
```

---

## 🐹 Go Backend Prompts

### **Create Complete Feature (CRUD)**

```
Create a complete CRUD feature for [entity name] in Go following Clean Architecture.

Entity: [Employee/Document/Company]
Fields: [list fields with types]

Requirements:
1. Domain layer: Pure entity struct in internal/domain/[entity].go
2. Repository interface in internal/service/[entity]_service.go
3. Repository implementation in internal/repository/[entity]_postgres.go
4. Service layer in internal/service/[entity]_service.go
5. HTTP handlers in internal/handler/[entity]_handler.go
6. Chi router registration
7. Include validation using go-playground/validator
8. Context as first parameter in all methods
9. Proper error handling with wrapped errors
10. Logging with zerolog
11. Unit tests for service layer

Example entity structure:
type Employee struct {
    ID          string    `json:"id" db:"id" validate:"required,uuid"`
    Name        string    `json:"name" db:"name" validate:"required,min=2"`
    // ... other fields
}
```

### **Create Repository Layer**

```
Create PostgreSQL repository implementation for [Entity].

Requirements:
- File: internal/repository/[entity]_postgres.go
- Struct: Postgres[Entity]Repository with *sqlx.DB
- Methods: Create, FindByID, Update, Delete, List
- Use context.Context as first parameter
- Use prepared statements or sqlx.Named
- Wrap errors with fmt.Errorf and %w
- Add proper SQL queries with placeholders
- Include transaction support for complex operations
- Return domain entities, not database models
```

### **Create Service Layer**

```
Create service layer for [Entity] with business logic.

Requirements:
- File: internal/service/[entity]_service.go
- Define repository interface in this file
- Struct: [Entity]Service with dependencies injected
- Constructor: New[Entity]Service with dependency injection
- Methods: implement business logic
- Validation before calling repository
- Proper error handling with context
- Logging for important operations
- Return DTOs or domain entities
```

### **Create HTTP Handlers**

```
Create HTTP handlers for [Entity] using Chi router.

Requirements:
- File: internal/handler/[entity]_handler.go
- Struct: [Entity]Handler with service dependency
- Methods: Create, Get, Update, Delete, List
- JSON request/response
- Proper HTTP status codes
- Error response format: {"error": "message"}
- Request validation
- Register routes method using Chi router
- Include middleware: logging, recovery
```

### **Create Middleware**

```
Create HTTP middleware for [purpose].

Requirements:
- File: internal/middleware/[name].go
- Function signature: func(next http.Handler) http.Handler
- Chain-compatible with Chi router
- Include logging
- Handle errors gracefully
- Set appropriate headers
```

---

## ⚛️ Next.js Frontend Prompts

### **Create Page Component**

```
Create a Next.js 15 page component for [feature].

Requirements:
- File: src/app/[route]/page.tsx
- Server Component by default (no 'use client' unless needed)
- TypeScript with full typing
- Use TanStack Query for data fetching
- shadcn/ui components
- Tailwind CSS for styling
- Responsive mobile-first design
- Loading state with loading.tsx
- Error handling with error.tsx
- Include metadata export
```

### **Create Client Component**

```
Create a client-side React component [ComponentName].

Requirements:
- 'use client' directive at top
- File: src/components/[feature]/[name].tsx
- TypeScript interface for props
- No 'any' types
- Use React hooks if needed (useState, useEffect)
- shadcn/ui components
- Tailwind CSS (no inline styles)
- Responsive design
- Accessible (ARIA labels)
- Loading and error states
```

### **Create Form Component**

```
Create a form component for [purpose].

Requirements:
- Use React Hook Form
- Zod schema for validation
- shadcn/ui Form components
- TypeScript interfaces
- Handle submission with loading state
- Display validation errors
- Optimistic UI updates
- Toast notifications on success/error
- Accessible form fields
```

### **Create API Hook**

```
Create custom React hook for [entity] API operations.

Requirements:
- File: src/hooks/use-[entity].ts
- Use TanStack Query (useQuery, useMutation)
- TypeScript types imported from types/[entity].ts
- Query key management
- Cache invalidation on mutations
- Optimistic updates where appropriate
- Error handling
- Loading states
```

### **Create Type Definitions**

```
Create TypeScript type definitions for [entity].

Requirements:
- File: src/types/[entity].ts
- Interface for entity
- Interface for create/update requests
- Interface for API responses
- Enum types if needed
- Export all types
- Match backend API structure
- Use ISO date strings for dates
```

---

## 🗄️ Database Prompts

### **Create Migration**

```
Create PostgreSQL migration for [table_name].

Requirements:
- Up migration: migrations/[timestamp]_[name].up.sql
- Down migration: migrations/[timestamp]_[name].down.sql
- Table: [table_name]
- Columns: [list with types]
- Constraints: primary key, foreign keys, unique, not null
- Indexes: [specify which columns]
- Use UUID for IDs
- Include created_at, updated_at timestamps
```

### **Create SQL Queries (sqlc)**

```
Create SQL queries for [entity] using sqlc.

Requirements:
- File: queries/[entity].sql
- Queries: InsertEmployee, GetEmployee, UpdateEmployee, DeleteEmployee, ListEmployees
- Use proper PostgreSQL syntax
- Use $1, $2 placeholders
- Include comments for sqlc generation
- Return proper types
- Handle NULLable columns
```

---

## 🎨 UI Component Prompts

### **Create shadcn/ui Based Component**

```
Create a [component type] using shadcn/ui primitives.

Requirements:
- Use existing shadcn/ui components: Button, Card, Input, etc.
- TypeScript with proper props typing
- Tailwind CSS classes
- Variants using cva (if needed)
- Forward refs where appropriate
- Accessible (ARIA)
- Dark mode support (via Tailwind)
- Responsive design
```

### **Create Data Table**

```
Create a data table component for [entity] using TanStack Table.

Requirements:
- Use TanStack Table v8
- Column definitions with TypeScript
- Sorting, filtering, pagination
- Row selection (if needed)
- shadcn/ui Table components
- Responsive design (mobile cards)
- Loading skeleton
- Empty state
- Action column with Edit/Delete
```

---

## 🧪 Testing Prompts

### **Create Unit Test (Go)**

```
Create unit tests for [service/handler/repository].

Requirements:
- File: [original_file]_test.go
- Use testify/assert and testify/mock
- Table-driven tests
- Test success cases
- Test error cases
- Mock dependencies
- Test edge cases
- Proper test naming: Test[FunctionName]_[Scenario]
- Setup and teardown if needed
```

### **Create Component Test (React)**

```
Create tests for React component [ComponentName].

Requirements:
- Use React Testing Library
- Test rendering
- Test user interactions
- Test different states (loading, error, success)
- Test accessibility
- Mock API calls
- Snapshot tests (if appropriate)
```

---

## 📊 Common Patterns

### **Error Handling Pattern (Go)**

```
Implement proper error handling for [function].

Pattern:
- Use fmt.Errorf with %w to wrap errors
- Add context to error messages
- Log errors with appropriate level
- Return custom error types where needed
- Example: return fmt.Errorf("failed to create employee: %w", err)
```

### **Response Pattern (HTTP)**

```
Create standard JSON response helpers.

Requirements:
- respondJSON(w, status, data)
- respondError(w, status, message)
- respondCreated(w, data)
- respondNoContent(w)
- Set Content-Type: application/json
- Handle marshaling errors
```

### **Validation Pattern**

```
Add validation to [struct/handler].

Go:
- Use struct tags: validate:"required,min=2,max=100"
- Call validate.Struct(data)
- Return validation errors with proper messages

TypeScript:
- Use Zod schemas
- Define schema: z.object({ ... })
- Validate with schema.parse() or schema.safeParse()
```

---

## 🔍 Code Review Checklist

After AI generates code, verify:

**Go:**
- [ ] All errors handled (no `_` for errors)
- [ ] Context as first parameter
- [ ] Interfaces defined where consumed
- [ ] Dependency injection used
- [ ] Proper logging added
- [ ] Code formatted with `gofmt`

**TypeScript:**
- [ ] No `any` types
- [ ] Props interface defined
- [ ] 'use client' only when needed
- [ ] Responsive design implemented
- [ ] Accessible (ARIA labels)
- [ ] Error handling included

---

## 💡 Pro Tips

1. **Be Specific:** More details = better code
2. **Reference Files:** "Following the pattern in employee_service.go, create..."
3. **Include Context:** "This is for a manpower management system where..."
4. **Iterate:** Start with basic structure, then refine
5. **Review Output:** AI is helpful but not perfect - always review generated code
6. **Use Examples:** "Like the EmployeeCard component, but for documents..."

---

## 🎯 Example Complete Feature Prompt

```
Create a complete feature for managing Documents in the manpower system.

Context:
- Documents belong to employees (foreign key relationship)
- Each document has: type (flexible string), expiry_date, file_url, last_updated
- Need CRUD operations + file upload to S3
- Need to check expiry status (valid, expiring, expired)

Backend Requirements:
1. Domain entity in internal/domain/document.go
2. Repository interface and implementation
3. Service layer with:
   - CreateDocument (with S3 upload)
   - UpdateDocument
   - DeleteDocument
   - GetDocument
   - ListDocumentsByEmployee
   - CheckExpiryStatus method
4. HTTP handlers with Chi routes
5. File upload handling (multipart/form-data)

Frontend Requirements:
1. Types in src/types/document.ts
2. API hooks in src/hooks/use-documents.ts
3. DocumentCard component
4. DocumentList component
5. DocumentForm with file upload
6. DocumentUpload component with drag-and-drop

Follow Clean Architecture, handle all errors, include logging, use TypeScript with full typing.
```

---

**Remember:** These prompts are starting points. Adjust based on your specific needs and always review generated code for quality and correctness.
