-- 004_primary_document.sql
-- Makes expiry_date optional and adds a "primary" flag to documents.
-- The primary document is the one tracked for visa/permit expiry alerts.

-- 1. Allow documents without an expiry date (contracts, certificates, etc.)
ALTER TABLE documents ALTER COLUMN expiry_date DROP NOT NULL;

-- 2. Add primary flag — only one per employee, meaning "this is the
--    document we track for fines/expiry".
ALTER TABLE documents ADD COLUMN IF NOT EXISTS is_primary BOOLEAN DEFAULT FALSE;

-- 3. Partial unique index: guarantees at most ONE primary document per employee.
--    Application code toggles it, but the DB enforces the constraint.
CREATE UNIQUE INDEX IF NOT EXISTS idx_documents_primary_per_employee
  ON documents (employee_id) WHERE is_primary = TRUE;
