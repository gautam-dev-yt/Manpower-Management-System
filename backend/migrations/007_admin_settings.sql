-- Migration 007: Admin Settings
-- Adds document_types and compliance_rules tables for admin-configurable compliance.
-- All changes are additive. No existing tables are modified.
-- Seed data matches current hardcoded values exactly — day-1 behavior is identical.

-- ── 1. Document Types Table ──────────────────────────────────────

CREATE TABLE IF NOT EXISTS document_types (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    doc_type           VARCHAR(100) NOT NULL UNIQUE,
    display_name       VARCHAR(200) NOT NULL,
    is_mandatory       BOOLEAN NOT NULL DEFAULT FALSE,
    has_expiry         BOOLEAN NOT NULL DEFAULT TRUE,
    number_label       VARCHAR(100) NOT NULL DEFAULT 'Document Number',
    number_placeholder VARCHAR(200) NOT NULL DEFAULT '',
    expiry_label       VARCHAR(100) NOT NULL DEFAULT 'Expiry Date',
    sort_order         INT NOT NULL DEFAULT 100,
    metadata_fields    JSONB NOT NULL DEFAULT '[]',
    is_system          BOOLEAN NOT NULL DEFAULT FALSE,
    is_active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ── 2. Compliance Rules Table ────────────────────────────────────

CREATE TABLE IF NOT EXISTS compliance_rules (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id        UUID REFERENCES companies(id) ON DELETE CASCADE,
    doc_type          VARCHAR(100) NOT NULL,
    grace_period_days INT NOT NULL DEFAULT 0,
    fine_per_day      DECIMAL(10,2) NOT NULL DEFAULT 0,
    fine_type         VARCHAR(20) NOT NULL DEFAULT 'daily',
    fine_cap          DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_mandatory      BOOLEAN DEFAULT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, doc_type)
);

-- Partial unique index for global defaults (company_id IS NULL)
CREATE UNIQUE INDEX IF NOT EXISTS idx_compliance_rules_global_unique
    ON compliance_rules(doc_type) WHERE company_id IS NULL;

-- Index for quick lookup during employee creation
CREATE INDEX IF NOT EXISTS idx_compliance_rules_company
    ON compliance_rules(company_id) WHERE company_id IS NOT NULL;

-- ── 3. Seed Document Types (only if table is empty) ─────────────

INSERT INTO document_types (doc_type, display_name, is_mandatory, has_expiry, number_label, number_placeholder, expiry_label, sort_order, metadata_fields, is_system)
SELECT * FROM (VALUES
    ('passport',         'Passport',                    TRUE,  TRUE, 'Passport Number',             'e.g. A12345678',          'Expiry Date',   10, '[{"key":"nationality","label":"Nationality","type":"text","placeholder":"e.g. Indian"},{"key":"issuing_country","label":"Issuing Country","type":"text","placeholder":"e.g. India"}]'::jsonb, TRUE),
    ('visa',             'Residence Visa',              TRUE,  TRUE, 'Visa Number',                 'e.g. 201/2024/1234567',   'Expiry Date',   20, '[{"key":"visa_type","label":"Visa Type","type":"select","options":[{"value":"employment","label":"Employment"},{"value":"residence","label":"Residence"},{"value":"mission","label":"Mission"},{"value":"green","label":"Green Visa"},{"value":"golden","label":"Golden Visa"}]},{"key":"sponsor","label":"Sponsor / Company","type":"text","placeholder":"Sponsoring company"},{"key":"linked_passport","label":"Linked Passport Number","type":"text","placeholder":"Passport number"}]'::jsonb, TRUE),
    ('emirates_id',      'Emirates ID',                 TRUE,  TRUE, 'Emirates ID Number',          'e.g. 784-1990-1234567-1', 'Expiry Date',   30, '[{"key":"linked_visa","label":"Linked Visa Number","type":"text","placeholder":"Visa number"}]'::jsonb, TRUE),
    ('work_permit',      'Work Permit / Labour Card',   TRUE,  TRUE, 'Permit / Labour Card Number', 'e.g. 1234567',            'Expiry Date',   40, '[{"key":"mohre_file_number","label":"MoHRE File Number","type":"text","placeholder":"e.g. 12345"},{"key":"job_title","label":"Job Title (on permit)","type":"text","placeholder":"e.g. Electrician"}]'::jsonb, TRUE),
    ('health_insurance', 'Health Insurance',             TRUE,  TRUE, 'Policy Number',               'e.g. POL-2024-12345',     'Expiry Date',   50, '[{"key":"insurer_name","label":"Insurance Provider","type":"text","placeholder":"e.g. Daman, Oman Insurance"},{"key":"coverage_amount","label":"Coverage Amount (AED)","type":"number","placeholder":"e.g. 250000"}]'::jsonb, TRUE),
    ('iloe_insurance',   'ILOE Insurance',              TRUE,  TRUE, 'Subscription ID',             'e.g. ILOE-2024-12345',    'Renewal Date',  60, '[{"key":"category","label":"Category","type":"select","options":[{"value":"A","label":"Category A (≤ AED 16,000 salary)"},{"value":"B","label":"Category B (> AED 16,000 salary)"}]},{"key":"subscription_status","label":"Subscription Status","type":"select","options":[{"value":"active","label":"Active"},{"value":"lapsed","label":"Lapsed"}]}]'::jsonb, TRUE),
    ('medical_fitness',  'Medical Fitness Certificate',  TRUE,  TRUE, 'Certificate Number',          'e.g. MED-2024-12345',     'Valid Until',   70, '[{"key":"test_date","label":"Test Date","type":"date"},{"key":"result","label":"Result","type":"select","options":[{"value":"fit","label":"Fit"},{"value":"unfit","label":"Unfit"}]}]'::jsonb, TRUE),
    ('trade_license',    'Trade License',               FALSE, TRUE, 'License Number',              'e.g. TL-12345',           'Expiry Date',   80, '[]'::jsonb, TRUE),
    ('other',            'Other',                       FALSE, TRUE, 'Document Number',             'e.g. DOC-12345',          'Expiry Date',  999, '[{"key":"custom_name","label":"Document Name","type":"text","placeholder":"e.g. Certificate of Good Conduct","required":true}]'::jsonb, TRUE)
) AS seed(a, b, c, d, e, f, g, h, i, j)
WHERE NOT EXISTS (SELECT 1 FROM document_types LIMIT 1);

-- ── 4. Seed Global Compliance Rules (only if table is empty) ────

INSERT INTO compliance_rules (company_id, doc_type, grace_period_days, fine_per_day, fine_type, fine_cap)
SELECT NULL, a, b, c, d, e FROM (VALUES
    ('passport',         0,   0.00, 'daily',    0.00),
    ('visa',             0,  50.00, 'daily',    0.00),
    ('emirates_id',     30,  20.00, 'daily', 1000.00),
    ('work_permit',     50, 500.00, 'one_time', 500.00),
    ('health_insurance', 0, 500.00, 'monthly', 150000.00),
    ('iloe_insurance',   0, 400.00, 'one_time',  400.00),
    ('medical_fitness',  0,   0.00, 'daily',    0.00)
) AS seed(a, b, c, d, e)
WHERE NOT EXISTS (SELECT 1 FROM compliance_rules LIMIT 1);
