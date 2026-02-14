package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"manpower-backend/internal/database"
	"manpower-backend/internal/models"
)

// EmployeeHandler handles employee-related HTTP requests.
type EmployeeHandler struct {
	db database.Service
}

// NewEmployeeHandler creates a new EmployeeHandler.
func NewEmployeeHandler(db database.Service) *EmployeeHandler {
	return &EmployeeHandler{db: db}
}

// ── Columns ────────────────────────────────────────────────────
// Central list so Create/GetByID/List all stay in sync.
const employeeCols = `e.id, e.company_id, e.name, e.trade, e.mobile,
	e.joining_date::text, e.photo_url,
	e.gender, e.date_of_birth::text, e.nationality, e.passport_number,
	e.native_location, e.current_location, e.salary, e.status,
	e.created_at, e.updated_at`

func scanEmployee(scanner interface {
	Scan(dest ...interface{}) error
}, emp *models.Employee) error {
	return scanner.Scan(
		&emp.ID, &emp.CompanyID, &emp.Name, &emp.Trade, &emp.Mobile,
		&emp.JoiningDate, &emp.PhotoURL,
		&emp.Gender, &emp.DateOfBirth, &emp.Nationality, &emp.PassportNumber,
		&emp.NativeLocation, &emp.CurrentLocation, &emp.Salary, &emp.Status,
		&emp.CreatedAt, &emp.UpdatedAt,
	)
}

func scanEmployeeWithCompany(scanner interface {
	Scan(dest ...interface{}) error
}, emp *models.EmployeeWithCompany) error {
	return scanner.Scan(
		&emp.ID, &emp.CompanyID, &emp.Name, &emp.Trade, &emp.Mobile,
		&emp.JoiningDate, &emp.PhotoURL,
		&emp.Gender, &emp.DateOfBirth, &emp.Nationality, &emp.PassportNumber,
		&emp.NativeLocation, &emp.CurrentLocation, &emp.Salary, &emp.Status,
		&emp.CreatedAt, &emp.UpdatedAt,
		&emp.CompanyName, &emp.DocStatus,
		&emp.ExpiryDaysLeft, &emp.PrimaryDocType,
	)
}

// ── Create ─────────────────────────────────────────────────────

// Create handles POST /api/employees
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if errs := req.Validate(); len(errs) > 0 {
		JSON(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"error":   "Validation failed",
			"details": errs,
		})
		return
	}

	// Default status to "active" if not provided
	if req.Status == "" {
		req.Status = "active"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	var employee models.Employee
	err := pool.QueryRow(ctx, `
		INSERT INTO employees AS e (
			company_id, name, trade, mobile, joining_date, photo_url,
			gender, date_of_birth, nationality, passport_number,
			native_location, current_location, salary, status
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING `+employeeCols,
		req.CompanyID, req.Name, req.Trade, req.Mobile, req.JoiningDate,
		nilIfEmpty(req.PhotoURL),
		req.Gender, req.DateOfBirth, req.Nationality, req.PassportNumber,
		req.NativeLocation, req.CurrentLocation, req.Salary, req.Status,
	).Scan(
		&employee.ID, &employee.CompanyID, &employee.Name,
		&employee.Trade, &employee.Mobile, &employee.JoiningDate,
		&employee.PhotoURL,
		&employee.Gender, &employee.DateOfBirth, &employee.Nationality, &employee.PassportNumber,
		&employee.NativeLocation, &employee.CurrentLocation, &employee.Salary, &employee.Status,
		&employee.CreatedAt, &employee.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error creating employee: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to create employee")
		return
	}

	JSON(w, http.StatusCreated, map[string]interface{}{
		"data":    employee,
		"message": "Employee created successfully",
	})
}

// ── List ───────────────────────────────────────────────────────

// List handles GET /api/employees
func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	companyID := q.Get("company_id")
	trade := q.Get("trade")
	search := q.Get("search")
	docStatus := q.Get("status")     // document status filter
	empStatus := q.Get("emp_status") // employee active/inactive filter
	nationality := q.Get("nationality")
	sortBy := q.Get("sort_by")
	sortOrder := q.Get("sort_order")

	// Whitelist allowed sort columns
	allowedSorts := map[string]string{
		"name":         "e.name",
		"joining_date": "e.joining_date",
		"created_at":   "e.created_at",
		"salary":       "e.salary",
	}
	sortCol, ok := allowedSorts[sortBy]
	if !ok {
		sortCol = "e.name"
	}
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	// Build dynamic WHERE clause
	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if companyID != "" {
		where += fmt.Sprintf(" AND e.company_id = $%d", argIdx)
		args = append(args, companyID)
		argIdx++
	}
	if trade != "" {
		where += fmt.Sprintf(" AND e.trade = $%d", argIdx)
		args = append(args, trade)
		argIdx++
	}
	if search != "" {
		where += fmt.Sprintf(" AND e.name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if empStatus != "" {
		where += fmt.Sprintf(" AND e.status = $%d", argIdx)
		args = append(args, empStatus)
		argIdx++
	}
	if nationality != "" {
		where += fmt.Sprintf(" AND e.nationality ILIKE $%d", argIdx)
		args = append(args, "%"+nationality+"%")
		argIdx++
	}

	// Doc status filter — now uses the primary document via LEFT JOIN
	statusFilter := ""
	if docStatus == "expiring" {
		statusFilter = " AND pd.expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days'"
	} else if docStatus == "expired" {
		statusFilter = " AND pd.expiry_date < CURRENT_DATE"
	} else if docStatus == "valid" || docStatus == "active" {
		statusFilter = " AND pd.expiry_date > CURRENT_DATE + INTERVAL '30 days'"
	}

	// Count total for pagination
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM employees e
		LEFT JOIN documents pd ON pd.employee_id = e.id AND pd.is_primary = TRUE
		%s %s
	`, where, statusFilter)
	var total int
	if err := pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		log.Printf("Error counting employees: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch employees")
		return
	}

	// Fetch employees — doc status derived from primary document only
	query := fmt.Sprintf(`
		SELECT 
			%s,
			c.name AS company_name,
			CASE
				WHEN pd.expiry_date IS NULL THEN 'none'
				WHEN pd.expiry_date < CURRENT_DATE THEN 'expired'
				WHEN pd.expiry_date <= CURRENT_DATE + INTERVAL '30 days' THEN 'expiring'
				ELSE 'valid'
			END AS doc_status,
			(pd.expiry_date - CURRENT_DATE) AS expiry_days_left,
			pd.document_type AS primary_doc_type
		FROM employees e
		JOIN companies c ON e.company_id = c.id
		LEFT JOIN documents pd ON pd.employee_id = e.id AND pd.is_primary = TRUE
		%s %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, employeeCols, where, statusFilter, sortCol, sortOrder, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Error querying employees: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch employees")
		return
	}
	defer rows.Close()

	employees := []models.EmployeeWithCompany{}
	for rows.Next() {
		var emp models.EmployeeWithCompany
		if err := scanEmployeeWithCompany(rows, &emp); err != nil {
			log.Printf("Error scanning employee: %v", err)
			continue
		}
		employees = append(employees, emp)
	}

	JSON(w, http.StatusOK, PaginatedResponse{
		Data: employees,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		},
	})
}

// ── GetByID ────────────────────────────────────────────────────

// GetByID handles GET /api/employees/{id}
func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Employee ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	var emp models.EmployeeWithCompany
	err := pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT 
			%s,
			c.name AS company_name,
			CASE
				WHEN pd.expiry_date IS NULL THEN 'none'
				WHEN pd.expiry_date < CURRENT_DATE THEN 'expired'
				WHEN pd.expiry_date <= CURRENT_DATE + INTERVAL '30 days' THEN 'expiring'
				ELSE 'valid'
			END AS doc_status,
			(pd.expiry_date - CURRENT_DATE) AS expiry_days_left,
			pd.document_type AS primary_doc_type
		FROM employees e
		JOIN companies c ON e.company_id = c.id
		LEFT JOIN documents pd ON pd.employee_id = e.id AND pd.is_primary = TRUE
		WHERE e.id = $1
	`, employeeCols), id).Scan(
		&emp.ID, &emp.CompanyID, &emp.Name, &emp.Trade, &emp.Mobile,
		&emp.JoiningDate, &emp.PhotoURL,
		&emp.Gender, &emp.DateOfBirth, &emp.Nationality, &emp.PassportNumber,
		&emp.NativeLocation, &emp.CurrentLocation, &emp.Salary, &emp.Status,
		&emp.CreatedAt, &emp.UpdatedAt,
		&emp.CompanyName, &emp.DocStatus,
		&emp.ExpiryDaysLeft, &emp.PrimaryDocType,
	)
	if err != nil {
		log.Printf("Error fetching employee %s: %v", id, err)
		JSONError(w, http.StatusNotFound, "Employee not found")
		return
	}

	// Fetch documents for this employee (with nullable expiry + primary flag)
	docRows, err := pool.Query(ctx, `
		SELECT d.id, d.employee_id, d.document_type,
			COALESCE(d.expiry_date::text, ''),
			d.is_primary,
			d.file_url, d.file_name, d.file_size, d.file_type,
			d.last_updated, d.created_at
		FROM documents d
		WHERE d.employee_id = $1
		ORDER BY d.is_primary DESC, d.expiry_date ASC NULLS LAST
	`, id)
	if err != nil {
		log.Printf("Error fetching documents for employee %s: %v", id, err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch employee documents")
		return
	}
	defer docRows.Close()

	documents := []models.Document{}
	for docRows.Next() {
		var doc models.Document
		var expiryRaw string
		if err := docRows.Scan(
			&doc.ID, &doc.EmployeeID, &doc.DocumentType,
			&expiryRaw,
			&doc.IsPrimary,
			&doc.FileURL, &doc.FileName, &doc.FileSize, &doc.FileType,
			&doc.LastUpdated, &doc.CreatedAt,
		); err != nil {
			log.Printf("Error scanning document: %v", err)
			continue
		}
		if expiryRaw != "" {
			doc.ExpiryDate = &expiryRaw
		}
		documents = append(documents, doc)
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data":      emp,
		"documents": documents,
	})
}

// ── Update ─────────────────────────────────────────────────────

// Update handles PUT /api/employees/{id}
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Employee ID is required")
		return
	}

	var req models.UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	// Build dynamic SET clause — only update provided fields
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	addField := func(col string, val interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, val)
		argIdx++
	}

	if req.Name != nil {
		addField("name", *req.Name)
	}
	if req.Trade != nil {
		addField("trade", *req.Trade)
	}
	if req.CompanyID != nil {
		addField("company_id", *req.CompanyID)
	}
	if req.Mobile != nil {
		addField("mobile", *req.Mobile)
	}
	if req.JoiningDate != nil {
		addField("joining_date", *req.JoiningDate)
	}
	if req.PhotoURL != nil {
		addField("photo_url", *req.PhotoURL)
	}
	if req.Gender != nil {
		addField("gender", *req.Gender)
	}
	if req.DateOfBirth != nil {
		addField("date_of_birth", *req.DateOfBirth)
	}
	if req.Nationality != nil {
		addField("nationality", *req.Nationality)
	}
	if req.PassportNumber != nil {
		addField("passport_number", *req.PassportNumber)
	}
	if req.NativeLocation != nil {
		addField("native_location", *req.NativeLocation)
	}
	if req.CurrentLocation != nil {
		addField("current_location", *req.CurrentLocation)
	}
	if req.Salary != nil {
		addField("salary", *req.Salary)
	}
	if req.Status != nil {
		addField("status", *req.Status)
	}

	if len(setClauses) == 0 {
		JSONError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	// Always update updated_at
	setClauses = append(setClauses, "updated_at = NOW()")

	query := fmt.Sprintf(`
		UPDATE employees AS e SET %s
		WHERE id = $%d
		RETURNING %s
	`, strings.Join(setClauses, ", "), argIdx, employeeCols)
	args = append(args, id)

	var employee models.Employee
	err := pool.QueryRow(ctx, query, args...).Scan(
		&employee.ID, &employee.CompanyID, &employee.Name,
		&employee.Trade, &employee.Mobile, &employee.JoiningDate,
		&employee.PhotoURL,
		&employee.Gender, &employee.DateOfBirth, &employee.Nationality, &employee.PassportNumber,
		&employee.NativeLocation, &employee.CurrentLocation, &employee.Salary, &employee.Status,
		&employee.CreatedAt, &employee.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error updating employee %s: %v", id, err)
		JSONError(w, http.StatusNotFound, "Employee not found")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data":    employee,
		"message": "Employee updated successfully",
	})
}

// ── Delete ─────────────────────────────────────────────────────

// Delete handles DELETE /api/employees/{id}
func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Employee ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	tag, err := pool.Exec(ctx, "DELETE FROM employees WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting employee %s: %v", id, err)
		JSONError(w, http.StatusInternalServerError, "Failed to delete employee")
		return
	}

	if tag.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Employee not found")
		return
	}

	JSON(w, http.StatusOK, map[string]string{
		"message": "Employee deleted successfully",
	})
}

// ── Export ──────────────────────────────────────────────────────

// Export handles GET /api/employees/export — returns CSV
func (h *EmployeeHandler) Export(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	rows, err := pool.Query(ctx, `
		SELECT e.name, e.trade, e.mobile, e.joining_date::text,
			COALESCE(e.gender,''), COALESCE(e.nationality,''),
			COALESCE(e.passport_number,''), COALESCE(e.native_location,''),
			COALESCE(e.current_location,''), COALESCE(e.salary::text,''),
			e.status, c.name
		FROM employees e
		JOIN companies c ON e.company_id = c.id
		ORDER BY e.name ASC
	`)
	if err != nil {
		log.Printf("Error exporting employees: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to export")
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=employees.csv")

	// Write CSV header
	fmt.Fprintln(w, "Name,Trade,Mobile,Joining Date,Gender,Nationality,Passport,Native Location,Current Location,Salary,Status,Company")

	for rows.Next() {
		var name, trade, mobile, joiningDate, gender, nationality, passport, nativeLoc, currentLoc, salary, status, company string
		if err := rows.Scan(&name, &trade, &mobile, &joiningDate, &gender, &nationality, &passport, &nativeLoc, &currentLoc, &salary, &status, &company); err != nil {
			continue
		}
		fmt.Fprintf(w, "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			csvEscape(name), csvEscape(trade), csvEscape(mobile), joiningDate,
			gender, nationality, passport,
			csvEscape(nativeLoc), csvEscape(currentLoc), salary, status, csvEscape(company))
	}
}

// ── Helpers ────────────────────────────────────────────────────

// nilIfEmpty returns nil if the string is empty, otherwise returns a pointer to it.
func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// csvEscape wraps a value in quotes if it contains commas.
func csvEscape(s string) string {
	if strings.Contains(s, ",") || strings.Contains(s, "\"") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}
