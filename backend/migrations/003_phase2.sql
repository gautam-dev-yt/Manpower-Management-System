-- Phase 2: Extended employee profile, salary tracker, notifications
-- Safe to run multiple times (IF NOT EXISTS / ADD COLUMN IF NOT EXISTS)

-- ── Employee profile extensions ──────────────────────────────────
ALTER TABLE employees ADD COLUMN IF NOT EXISTS gender VARCHAR(10);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS date_of_birth DATE;
ALTER TABLE employees ADD COLUMN IF NOT EXISTS nationality VARCHAR(60);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS passport_number VARCHAR(30);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS native_location VARCHAR(120);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS current_location VARCHAR(120);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS salary NUMERIC(12,2);
ALTER TABLE employees ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'active';

CREATE INDEX IF NOT EXISTS idx_employees_status ON employees(status);
CREATE INDEX IF NOT EXISTS idx_employees_nationality ON employees(nationality);

-- ── Salary records ───────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS salary_records (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    month       INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    year        INT NOT NULL CHECK (year >= 2020),
    amount      NUMERIC(12,2) NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    paid_date   DATE,
    notes       TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(employee_id, month, year)
);

CREATE INDEX IF NOT EXISTS idx_salary_employee ON salary_records(employee_id);
CREATE INDEX IF NOT EXISTS idx_salary_month_year ON salary_records(year, month);
CREATE INDEX IF NOT EXISTS idx_salary_status ON salary_records(status);

-- ── Notifications ────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS notifications (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(200) NOT NULL,
    message     TEXT NOT NULL,
    type        VARCHAR(30) NOT NULL,
    read        BOOLEAN DEFAULT FALSE,
    entity_type VARCHAR(30),
    entity_id   UUID,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id, read, created_at DESC);
