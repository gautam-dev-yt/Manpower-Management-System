# Current Status — Manpower Management System

> **Purpose:** Append-only project state so any new chat or AI IDE has current context. When you change migrations, deploy config, or ship features, append a 1–3 line entry below. Do not remove old entries; trim only very old ones if the file grows large (e.g. keep last 30).

---

**Last updated:** 2026-03-01 (b)

## Deploy state

- **Database:** Neon PostgreSQL — live; migrations applied through **010**.
- **Backend:** Render (Go/Chi API).
- **Frontend:** Vercel (Next.js 16).
- **Storage:** Cloudflare R2 (documents, employee photos); public CDN URLs.

## Latest migration

- `010_drop_deprecated_doc_cols.sql` — dropped deprecated document columns (grace/fine/mandatory now from compliance_rules/document_types).

## Recent changes (append here)

- 2026-02: Production live (Vercel, Render, Neon, R2); role system (009), document rework (008), admin settings (007); compliance and settings in place.
- 2026-02-28: Added CURRENT_STATUS.md and Cursor rules/skills for AI/LLM context; PROJECT_ANALYSIS.md is canonical reference.
- 2026-03-01: Updated production_setup.sql to match schema through migration 010 (no deprecated doc columns, document_types field config, user_companies); suitable for fresh local Postgres or Neon.
- 2026-03-01: Dashboard/employees filter fix: "Active Documents" list now shows employees with at least one valid doc (EXISTS), not only full compliance; employees page normalizes URL `?status=active` → doc-status "Valid" so dropdown is not blank when redirected from dashboard.
- 2026-03-01: Dashboard donut chart now includes "In Grace" segment (orange). Company Compliance table replaced "Incomplete" column with "In Grace" and "Expiring" per company (backend + frontend). Company detail employee list now shows document compliance status (priority-based: penalty > grace > expiring > valid) with urgent doc name instead of bare "active" HR status.
- 2026-03-01: Companies list: compliance summary (penalty/grace/expiring counts) on each card; employee count excludes exited. Employee detail: top-level compliance badge with urgent doc name. Excluded exited employees from company counts (list, detail, dashboard bar chart) and employee list default view (default filter "Active", backend excludes exit_type IS NULL unless emp_status=all).
