-- ============================================
-- Migration 005: Currency & Settings
-- ============================================
-- Adds a per-company currency preference.
-- Safe to run multiple times (IF NOT EXISTS / ADD COLUMN IF NOT EXISTS).

-- Company-level currency (default AED for existing data).
ALTER TABLE companies ADD COLUMN IF NOT EXISTS currency VARCHAR(10) DEFAULT 'AED';
