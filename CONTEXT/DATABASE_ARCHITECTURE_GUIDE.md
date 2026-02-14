# Database & Architecture Design — Why We Chose What We Chose

> A guide explaining the thinking behind every schema and architecture decision in this project.

---

## 1. How to Design a Database Schema

### Step 1: Identify Your "Nouns"

Read the requirement: *"Track employees across companies, manage document expiries"*

The nouns become your **tables**:
- **Company** — an organization where employees work
- **Employee** — a person working at a company
- **Document** — a file (visa, passport) with an expiry date
- **User** — someone who logs in to manage data

### Step 2: Identify Relationships

Ask: "How do these relate?"

| Relationship | Type | Example |
|-------------|------|---------|
| User → Company | One-to-Many | One user manages multiple companies |
| Company → Employee | One-to-Many | One company has many employees |
| Employee → Document | One-to-Many | One employee has many documents |

**One-to-Many** means: the "many" side gets a foreign key. So:
- `companies` has `user_id` (points to its owner)
- `employees` has `company_id` (points to its company)
- `documents` has `employee_id` (points to its employee)

### Step 3: Choose Primary Keys

| Option | Pros | Cons | Our Choice |
|--------|------|------|------------|
| Auto-increment integer (`1, 2, 3`) | Simple, small, fast | Predictable (security risk), hard to merge databases | ❌ |
| UUID (`550e8400-e29b-...`) | Globally unique, safe to expose in URLs, easy to merge | Slightly larger (16 bytes vs 4) | ✅ |

We chose **UUID** because:
- URLs like `/employees/3` leak how many employees exist
- If you ever merge two databases, UUIDs won't collide
- PostgreSQL has `gen_random_uuid()` built-in — zero extra work

### Step 4: Decide Data Types

```sql
name VARCHAR(255)     -- not TEXT, because we want a max length
mobile VARCHAR(20)    -- phone numbers include +country code
expiry_date DATE      -- only need date, not time
file_size BIGINT      -- files can be > 2GB (INT max is ~2.1GB)
details JSONB         -- flexible key-value data (audit trail)
role VARCHAR(20)      -- enum-like: 'admin' | 'viewer'
```

**Why `VARCHAR(N)` over `TEXT`?**
- `TEXT` has no max — a bug could store entire files in a name field
- `VARCHAR(255)` is a safety net, not a performance choice

**Why `JSONB` for audit details?**
- Each action logs different fields changed
- Rigid columns would mean 50+ nullable columns
- `JSONB` is queryable: `WHERE details->>'field' = 'name'`

---

## 2. Why These Tables? Could We Do It Differently?

### Alternative: Single "People" table (rejected)

```sql
-- BAD: mixing users and employees in one table
CREATE TABLE people (
    id UUID, name TEXT, type TEXT, -- 'user' or 'employee'
    email TEXT, password TEXT,     -- null for employees
    company_id UUID, trade TEXT,  -- null for users
);
```

**Why bad:**
- Half the columns are null for each row (wasteful)
- Can't have different validation rules
- Queries get ugly: `WHERE type = 'employee' AND ...`
- Violates **Single Responsibility** — one table, two purposes

### Alternative: Documents as JSON inside Employee (rejected)

```sql
-- BAD: embedding documents as JSON
CREATE TABLE employees (
    id UUID, name TEXT,
    documents JSONB  -- [{"type": "Visa", "expiry": "2026-03-15"}, ...]
);
```

**Why bad:**
- Can't efficiently query "all documents expiring in 30 days" across all employees
- No foreign keys — can't enforce data integrity
- No indexes on nested JSON fields (slow queries)
- Can't update one document without rewriting the entire array

**Rule of thumb:** If you need to **query, filter, or sort** by something, it should be its own table.

### Alternative: Separate table per document type (rejected)

```sql
-- BAD: one table per type
CREATE TABLE visas (employee_id UUID, expiry DATE, ...);
CREATE TABLE passports (employee_id UUID, expiry DATE, ...);
CREATE TABLE emirates_ids (employee_id UUID, expiry DATE, ...);
```

**Why bad:**
- Adding a new document type means creating a new table
- Dashboard queries need `UNION ALL` across all tables
- Code duplication — each table needs its own handler

**Our approach:** One `documents` table with a `document_type` column. Flexible, queryable, DRY.

---

## 3. Normalization — When to Split, When to Keep Together

### What is it?
**Normalization** = eliminate duplicate data by splitting into related tables.

### Our example:

**Before (denormalized — bad):**
```
| employee_name | company_name        | visa_expiry | passport_expiry |
|---------------|---------------------|-------------|-----------------|
| John Doe      | ABC Construction    | 2026-03-15  | 2027-01-20      |
| Jane Smith    | ABC Construction    | 2026-02-19  | NULL            |
```

Problems: "ABC Construction" is stored twice. If renamed, you must update every row.

**After (normalized — good):**
```
companies: { id: 1, name: "ABC Construction" }
employees: { id: 1, name: "John Doe", company_id: 1 }
documents: { id: 1, employee_id: 1, type: "Visa", expiry: "2026-03-15" }
documents: { id: 2, employee_id: 1, type: "Passport", expiry: "2027-01-20" }
```

Company name stored **once**. Change it in one place, reflected everywhere.

### When NOT to normalize

**Performance-critical reads.** If you always need employee + company name together, a JOIN every time can be slow at scale. Solution: **denormalize strategically** or use **database views**.

We're fine with JOINs at our scale (~1000 employees).

---

## 4. Indexes — Making Queries Fast

An index is like a **book's index page** — instead of reading every page to find "Visa", you look up the index.

```sql
CREATE INDEX idx_documents_expiry_date ON documents(expiry_date);
```

### Where we added indexes and why:

| Index | Query it speeds up |
|-------|-------------------|
| `idx_employees_company_id` | "Get all employees of company X" |
| `idx_employees_name` | "Search employees by name" |
| `idx_documents_employee_id` | "Get all documents of employee X" |
| `idx_documents_expiry_date` | "Find expiring documents" (dashboard) |
| `idx_companies_user_id` | "Get all companies of user X" |

### When NOT to index
- Columns you rarely query by
- Tables with < 1000 rows (full scan is fast enough)
- Columns with low cardinality (e.g., boolean `is_active` — only 2 values)

---

## 5. CASCADE — What Happens When You Delete

```sql
company_id UUID REFERENCES companies(id) ON DELETE CASCADE
```

This means: **if you delete a company, all its employees are automatically deleted.**

| Option | Behavior | When to use |
|--------|----------|-------------|
| `CASCADE` | Delete children too | Parent owns children completely |
| `SET NULL` | Set FK to NULL | Children can exist without parent |
| `RESTRICT` | Block delete if children exist | Safety — prevent accidental data loss |

**Our choices:**
- Delete company → cascade delete employees → cascade delete documents ✅
- Delete employee → cascade delete their documents ✅
- Delete user → **should NOT cascade** (what about their companies?) — we'll use `RESTRICT`

---

## 6. Architecture Decisions

### Why Go + PostgreSQL? (Not Node.js + MongoDB)

| Factor | Go + PostgreSQL | Node.js + MongoDB |
|--------|----------------|-------------------|
| Type safety | ✅ Compile-time errors | ❌ Runtime errors |
| Data integrity | ✅ Foreign keys, constraints | ❌ No built-in relations |
| Complex queries | ✅ SQL JOINs, aggregations | ❌ Aggregation pipeline is complex |
| Performance | ✅ Compiled binary, goroutines | ⚠️ Single-threaded event loop |
| Learning | ⚠️ Steeper initially | ✅ Familiar if you know JS |
| Schema evolution | ✅ Migrations | ✅ Schemaless (but that's also a risk) |

**Why PostgreSQL over MySQL?**
- `UUID` generation built-in (`gen_random_uuid()`)
- `JSONB` type for flexible data (audit trail)
- Better `INTERVAL` support (`CURRENT_DATE + INTERVAL '30 days'`)
- More powerful indexing (partial indexes, GIN for JSONB)

### Why REST API? (Not GraphQL)

| Factor | REST | GraphQL |
|--------|------|---------|
| Simplicity | ✅ URL = resource | ⚠️ Schema + resolvers |
| Caching | ✅ HTTP caching built-in | ❌ All POST requests |
| Learning curve | ✅ Easy | ⚠️ Extra concepts |
| Over-fetching | ⚠️ Fixed response shape | ✅ Client picks fields |
| For our scale | ✅ Perfect | ❌ Overkill |

REST is perfect for ~15 endpoints. GraphQL shines at 100+ endpoints with complex client needs.

### Why Chi Router? (Not Gin, Fiber, Echo)

| Router | Style | Our pick |
|--------|-------|----------|
| **Chi** | Stdlib-compatible, minimal, composable middleware | ✅ |
| Gin | Framework-like, custom context | Heavier than needed |
| Fiber | Express-like, uses fasthttp (not stdlib) | Non-standard |
| Echo | Similar to Gin | Also heavier |

Chi uses Go's standard `http.Handler` interface — no vendor lock-in. You can use any stdlib-compatible middleware.

### Why JWT for Auth? (Not Sessions)

| Factor | JWT (Token) | Sessions (Cookie) |
|--------|------------|-------------------|
| Storage | Client (localStorage) | Server (DB/Redis) |
| Scalability | ✅ Stateless — no server storage | ⚠️ Need shared session store |
| Mobile friendly | ✅ Works with any client | ❌ Cookies are browser-only |
| Revocation | ⚠️ Can't revoke until expiry | ✅ Delete from server |
| Our choice | ✅ Simple, fits our needs | — |

### Storage Interface Pattern (Dependency Injection)

```
┌──────────────────┐
│   Upload Handler  │ ← depends on interface, not implementation
│                    │
│  storage.Upload() │
└────────┬─────────┘
         │ uses interface
┌────────▼─────────┐
│   FileStorage     │ ← Go interface
│   interface       │
└────────┬─────────┘
         │ implemented by
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼────┐
│ Local  │ │  S3   │  ← swap by changing 1 line in main.go
│Storage │ │Storage│
└────────┘ └───────┘
```

**This is the Dependency Inversion Principle (the D in SOLID):**
- High-level code (handlers) depends on abstractions (interface)
- Low-level code (local/S3) implements the abstraction
- You can swap implementations without touching handlers

---

## 7. Summary: Our Schema Philosophy

1. **One table per real-world entity** — don't mix concerns
2. **Foreign keys for relationships** — let the DB enforce integrity
3. **UUIDs for IDs** — safe, globally unique, merge-friendly
4. **Index what you query** — but not everything
5. **CASCADE with care** — own children fully, restrict parents
6. **Interfaces for services** — swap implementations, not logic
7. **Normalize, but be practical** — JOINs are fine at our scale
