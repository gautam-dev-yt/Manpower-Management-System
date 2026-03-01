-- Drop the is_primary column and its partial index from documents.
-- This flag was never consumed by any query or business logic.

DROP INDEX IF EXISTS idx_documents_primary;
ALTER TABLE documents DROP COLUMN IF EXISTS is_primary;
