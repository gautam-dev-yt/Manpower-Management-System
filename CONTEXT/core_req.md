# Employee Document Compliance Management System

## Requirements Specification

---

## Why This App Exists

The client operates multiple companies in Dubai, UAE. Every employee working legally in Dubai requires 7 mandatory government documents. Each document has a different validity period, and when any document expires, the employer faces daily fines — not the employee.

The client is currently losing money because documents expire without anyone noticing until an inspection happens or a renewal gets blocked. Manual tracking with spreadsheets doesn't scale across multiple companies and dozens of employees.

The fines are real and steep:
- Visa overstay: AED 50/day (after 30-day grace)
- Emirates ID late renewal: AED 20/day (after 30-day grace, capped at AED 1,000)
- Work permit lapse: AED 500 flat, or AED 50,000–100,000 per worker if completely missing
- Health insurance non-compliance: AED 500–150,000/month per uncovered employee
- ILOE insurance missing: AED 400 per employee (blocks visa renewal)
- General labour violations: AED 100,000–1,000,000 under 2024 amendments

In 2024, MOHRE conducted 700,000 inspections and found 29,000 violations. This is not optional compliance — it's actively enforced.

**The client's core need:** One screen to see "am I in trouble?" and automated alerts so nothing ever expires unnoticed again.

---

## Core Requirements

### 1. Multi-Company Support

The client manages multiple companies. Each company is a separate legal entity with its own trade license, employees, and compliance obligations.

- Each company should have: name, currency, trade license number, establishment card number, MOHRE category (1/2/3), regulatory authority (MOHRE, JAFZA, DMCC, DAFZA, DIFC, etc.)
- Dashboard metrics and fine exposure must be viewable per-company and aggregated across all companies
- Employees belong to one company (their visa sponsor)

### 2. Seven Mandatory Documents Per Employee

Every employee in the UAE private sector must have these 7 documents. This is not optional — it's law.

| Document | Typical Validity | Grace Period | Fine After Grace |
|---|---|---|---|
| Passport | Varies (must have 6+ months remaining) | None | Blocks all other renewals |
| Residence Visa | 2 years | 30 days | AED 50/day |
| Emirates ID | Linked to visa (2 years) | 30 days | AED 20/day (max AED 1,000) |
| Work Permit / Labour Card | 2 years | 50–60 days | AED 500 flat; AED 50,000+ if fully missing |
| Health Insurance | 1 year | 30 days (Dubai) | AED 500–150,000/month |
| ILOE Insurance | 1 year | None | AED 400 one-time (blocks visa renewal) |
| Medical Fitness Certificate | Per visa cycle | None | Blocks visa issuance |

Note: "Work Permit" and "Labour Card" are the same thing in the current UAE system. The labour card is now digital, issued by MoHRE. Treat them as one document type displayed as "Work Permit / Labour Card".

### 3. Auto-Created Document Slots

When an employee is added, the system must automatically create 7 document tracking slots — one for each mandatory type. The user does NOT manually "add" mandatory documents one by one.

Each slot starts in an "Incomplete" state. The user fills in document details (number, dates, file scan) when they have them. Whatever they don't fill remains visibly incomplete.

This is important because:
- Employers often don't have all documents ready at once (new hire's visa may be processing)
- But nothing should be invisible or forgotten
- The system treats an incomplete slot as an active compliance gap

The "+ Add Document" action should only be used for additional non-mandatory documents (trade license copies, NOCs, educational certificates, etc.).

### 4. Document Fields

Every document needs more than just a type and expiry date. Required fields per document:

**Common to all types:**
- Document number (unique identifier on the document itself)
- Issue date
- Expiry date
- Grace period in days (auto-filled based on document type, but editable)
- File upload (scan/photo of the physical document)

**Type-specific additional fields:**
- **Visa:** Visa UID, visa type (Employment/Residence/Mission/Green/Golden), sponsor company
- **Emirates ID:** 15-digit EID number
- **Work Permit:** Permit number, MoHRE file number, job title
- **Health Insurance:** Policy number, insurer name, plan type, coverage amount
- **ILOE Insurance:** Category (A for salary ≤ AED 16,000, B for above), subscription status
- **Medical Fitness:** Certificate number, test result (Fit/Unfit)
- **Passport:** Issuing country (nationality can auto-fill from employee profile)

### 5. Document Lifecycle & Status

Each document must show its current state clearly. Status is computed automatically based on dates — the user never manually sets it.

**The five states:**

| State | When | What the user sees |
|---|---|---|
| Incomplete | Required fields are missing | Gray card, "Not yet tracked", prompt to complete |
| Valid | Fully filled, expiry more than 30 days away | Green indicator, "X days remaining" |
| Expiring Soon | Expiry within 30 days | Amber indicator, countdown, "Initiate renewal" prompt |
| Expired — In Grace | Past expiry but within grace period | Orange indicator, "X grace days remaining", "No fine yet" |
| Expired — Fine Active | Past expiry + past grace period | Red indicator, "AED X/day accumulating", total estimated fine |

### 6. Fine Calculation & Exposure

The system must calculate and display estimated fine exposure. This is the feature that makes the app worth using.

**Per document:** Based on document type's fine rate, fine type (daily/monthly/one-time), and days past grace period. Some fines have caps (Emirates ID caps at AED 1,000).

**Per employee:** Sum of all their document fines.

**Per company:** Sum of all employee fines within that company.

**Overall:** Sum across all companies.

The dashboard must show:
- Total accumulated fines right now (what the client would owe if inspected today)
- Daily burn rate (how much fines are growing per day if nothing is done)
- Company-wise breakdown

Fine rates and grace periods should be stored per-document so they can adapt if UAE regulations change. The defaults come from current law but an admin should be able to adjust them.

### 7. Document Dependency Alerts

UAE documents are interdependent. If one blocks another's renewal, the system must warn the user.

Key dependencies:
- **Passport → Visa:** Passport must have 6+ months validity to renew a visa. If passport expires before visa renewal window, flag it.
- **Health Insurance → Work Permit:** MOHRE will not renew a work permit without valid health insurance.
- **ILOE → Visa:** Unpaid ILOE fine (AED 400) blocks visa renewal processing.
- **Visa → Emirates ID:** Valid visa required to renew Emirates ID.
- **Medical Fitness → Visa:** Medical fitness certificate required for visa issuance/renewal.

Display these as warning banners on the employee detail page. Example: "Passport expires on Feb 20 before Visa renewal window — renew passport first."

### 8. Notification System

Automated multi-tier alerts. The bell icon already exists in the header but has no UI behind it.

**Alert schedule per document:**
- 90 days before expiry — planning reminder
- 60 days before expiry — initiate renewal (labour card renewal starts at 60 days per MOHRE)
- 30 days before expiry — urgent, most grace periods start here if missed
- 15 days before expiry — critical alert
- 7 days before expiry — final warning
- Daily after expiry — during grace: "X grace days left." After grace: "AED X/day accumulating, total: AED Y"

**Delivery:**
- In-app: bell icon with unread badge count, dropdown or page listing notifications
- Each notification links to the relevant employee's detail page
- "Mark as read" and "Mark all as read" actions
- Optional: email notifications (can be a later phase)

### 9. Dashboard

The dashboard is the primary screen. The client should be able to open the app and in 5 seconds know: "Am I in trouble?"

**Must show:**
- Fine exposure card: total accumulated fines + daily burn rate
- Document health counts: how many valid, expiring soon, in grace, penalty active, incomplete
- Company-wise compliance breakdown: each company's employee count, penalty count, fine exposure
- Completion rate: what percentage of employees have all 7 documents fully tracked
- Critical alerts: list of employees/documents in grace or penalty, sorted by urgency, showing fine amounts

### 10. Employee Detail Page

The employee profile page should show everything about that employee in one place.

**Profile section:** All fields visible — name, trade, company, mobile, joining date, gender, DOB, nationality, passport number, native location, current location, salary (with correct currency from company), employment status. Currently most of these are missing from the detail view despite being settable in the edit form.

**Document section:** 7 mandatory document cards in a grid layout. Each card shows its lifecycle state with appropriate visual treatment. A completion bar at the top ("5/7 mandatory documents tracked"). Dependency alerts if applicable.

**Additional documents section:** Below the mandatory cards. For trade licenses, NOCs, educational certificates, and any other documents the employer wants to track. The "+ Add Document" button lives here.

---

## Additional Requirements

### 11. Employee Onboarding Flow

Adding a new employee should feel natural and not overwhelming:

**Step 1 — Profile:** Fill in employee personal and employment details. Save.

**Step 2 — Documents:** The system shows all 7 mandatory document slots. The user fills what they have. They can skip any and come back later. Skipped slots remain as "Incomplete" on the profile.

This two-step flow acknowledges reality: a new hire's visa might still be processing, medical test might not be done yet. But the system never lets you forget what's pending.

### 12. Employee List

Should support filtering by:
- Company
- Document status (including the new states: incomplete, in_grace, penalty_active)
- Trade / job role
- Nationality
- Employee status (active, inactive, on_leave)

Sorting by name, joining date, company, trade.

Search by name.

Export to CSV with all profile and document status fields.

### 13. Reporting

- Compliance report: all employees across all companies with their document statuses
- Fine exposure report: company-wise and employee-wise breakdown of current and projected fines
- Expiry calendar: upcoming renewals across all companies in a timeline view
- All reports exportable as CSV

### 14. User Roles

- **Admin / Owner:** Full access to everything — all companies, all employees, settings, reports
- **HR / PRO:** Add/edit employees and documents for assigned companies, view dashboard and reports, cannot change system settings
- **Viewer:** Read-only access to dashboard and employee data

### 15. Company Details

Company cards and the company management page should show more than just a name:
- Currency (default AED, selectable)
- Trade license number
- Establishment card number
- MOHRE category
- Regulatory authority (MOHRE vs free zone authority)

This matters because free zone companies may have slightly different grace period rules.

---

## What This App Is NOT

- It is NOT a payroll system (salary module exists but is secondary)
- It is NOT an HR management platform (no leave tracking, performance reviews, etc.)
- It is NOT a government portal integration (no API connection to MOHRE/ICP — all data is manually entered)
- It does NOT replace legal advice (fine amounts are estimates based on published regulations)

The app is a **compliance tracking and alerting tool**. Its job is to make sure no document ever expires unnoticed and to show the financial cost of inaction.

---

## Success Criteria

The app is successful when:

1. The client opens the dashboard and instantly sees their compliance posture across all companies
2. No document expiry goes unnoticed — the notification system alerts well before any grace period ends
3. The client can see exactly how much money they're losing (or would lose) from expired documents
4. Adding a new employee automatically sets up all 7 document tracking slots with correct grace periods and fine rates
5. Document dependencies are flagged — no more "couldn't renew visa because passport expired and nobody noticed"
6. The completion percentage motivates the client's team to fill in all document details, not just some