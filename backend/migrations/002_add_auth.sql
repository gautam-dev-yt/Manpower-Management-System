-- ============================================
-- Migration 002: Authentication & Audit Trail
-- ============================================
-- Run this in pgAdmin Query Tool or via psql
-- This adds user accounts, company ownership, and activity tracking

-- Step 1: Create users table
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(100) NOT NULL,
    role          VARCHAR(20) NOT NULL DEFAULT 'admin',
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Step 2: Link companies to their owner
-- Nullable first — existing companies don't have a user yet
ALTER TABLE companies ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id);

-- Step 3: Create activity log for audit trail
CREATE TABLE IF NOT EXISTS activity_log (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID REFERENCES users(id),
    action      VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id   UUID NOT NULL,
    details     JSONB,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Step 4: Add indexes for common queries
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_companies_user_id ON companies(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_user_id ON activity_log(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_entity ON activity_log(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_created ON activity_log(created_at DESC);

-- ============================================
-- IMPORTANT: After creating your first user via the Register API,
-- run this to link existing companies to your account:
--
--   UPDATE companies SET user_id = 'your-user-uuid-here' WHERE user_id IS NULL;
--
-- Replace 'your-user-uuid-here' with the UUID from the registration response.
-- ============================================
