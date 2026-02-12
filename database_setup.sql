-- ============================================
-- MANPOWER MANAGEMENT SYSTEM - DATABASE SETUP
-- ============================================
-- Copy and paste this entire file into pgAdmin Query Tool

-- Step 1: Create Companies Table
CREATE TABLE IF NOT EXISTS companies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL UNIQUE,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Step 2: Create Employees Table
CREATE TABLE IF NOT EXISTS employees (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id      UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    trade           VARCHAR(100) NOT NULL,
    mobile          VARCHAR(20) NOT NULL,
    joining_date    DATE NOT NULL,
    photo_url       TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Step 3: Create Documents Table
CREATE TABLE IF NOT EXISTS documents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id     UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    document_type   VARCHAR(100) NOT NULL,
    expiry_date     DATE NOT NULL,
    file_url        TEXT NOT NULL,
    file_name       VARCHAR(255) NOT NULL,
    file_size       BIGINT NOT NULL,
    file_type       VARCHAR(50) NOT NULL,
    last_updated    TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Step 4: Create Indexes for Better Performance
CREATE INDEX IF NOT EXISTS idx_employees_company_id ON employees(company_id);
CREATE INDEX IF NOT EXISTS idx_employees_name ON employees(name);
CREATE INDEX IF NOT EXISTS idx_documents_employee_id ON documents(employee_id);
CREATE INDEX IF NOT EXISTS idx_documents_expiry_date ON documents(expiry_date);

-- Step 5: Insert Sample Companies
INSERT INTO companies (id, name) VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'ABC Construction Co.'),
('660e8400-e29b-41d4-a716-446655440001', 'XYZ Engineering')
ON CONFLICT (name) DO NOTHING;

-- Step 6: Insert Sample Employees
INSERT INTO employees (id, company_id, name, trade, mobile, joining_date) VALUES 
('770e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', 'John Doe', 'Welder', '+971501234567', '2024-01-15'),
('770e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440000', 'Jane Smith', 'Electrician', '+971501234568', '2024-02-20'),
('770e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440001', 'Mike Johnson', 'Engineer', '+971501234569', '2024-03-10'),
('770e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440000', 'Sarah Williams', 'Safety Officer', '+971501234570', '2023-11-05'),
('770e8400-e29b-41d4-a716-446655440005', '660e8400-e29b-41d4-a716-446655440001', 'Ahmed Ali', 'Plumber', '+971501234571', '2024-04-12')
ON CONFLICT (id) DO NOTHING;

-- Step 7: Insert Sample Documents (with varying expiry dates)
INSERT INTO documents (employee_id, document_type, expiry_date, file_url, file_name, file_size, file_type) VALUES 
-- Expiring in 30 days (Warning)
('770e8400-e29b-41d4-a716-446655440001', 'Visa', CURRENT_DATE + INTERVAL '30 days', 's3://bucket/visa1.pdf', 'visa_john_doe.pdf', 1048576, 'application/pdf'),
-- Expiring in 7 days (Urgent)
('770e8400-e29b-41d4-a716-446655440002', 'Passport', CURRENT_DATE + INTERVAL '7 days', 's3://bucket/passport1.pdf', 'passport_jane_smith.pdf', 2097152, 'application/pdf'),
-- Already expired (Expired)
('770e8400-e29b-41d4-a716-446655440003', 'Emirates ID', CURRENT_DATE - INTERVAL '2 days', 's3://bucket/eid1.pdf', 'eid_mike_johnson.pdf', 1500000, 'application/pdf'),
-- Valid (more than 30 days)
('770e8400-e29b-41d4-a716-446655440001', 'Labor Card', CURRENT_DATE + INTERVAL '90 days', 's3://bucket/labor1.pdf', 'labor_john_doe.pdf', 800000, 'application/pdf'),
('770e8400-e29b-41d4-a716-446655440004', 'Visa', CURRENT_DATE + INTERVAL '120 days', 's3://bucket/visa2.pdf', 'visa_sarah_williams.pdf', 950000, 'application/pdf'),
('770e8400-e29b-41d4-a716-446655440005', 'Passport', CURRENT_DATE + INTERVAL '180 days', 's3://bucket/passport2.pdf', 'passport_ahmed_ali.pdf', 1200000, 'application/pdf')
ON CONFLICT (id) DO NOTHING;

-- Verification: Check if everything was created successfully
SELECT 'Companies Count:' as info, COUNT(*) as count FROM companies
UNION ALL
SELECT 'Employees Count:', COUNT(*) FROM employees
UNION ALL
SELECT 'Documents Count:', COUNT(*) FROM documents;

-- THIS IS WHAT YOU SHOULD SEE:
-- Companies Count: 2
-- Employees Count: 5
-- Documents Count: 6
