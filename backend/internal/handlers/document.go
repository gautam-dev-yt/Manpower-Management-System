package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"manpower-backend/internal/ctxkeys"
	"manpower-backend/internal/database"
	"manpower-backend/internal/models"
)

// DocumentHandler handles document-related HTTP requests.
type DocumentHandler struct {
	db database.Service
}

// NewDocumentHandler creates a new DocumentHandler.
func NewDocumentHandler(db database.Service) *DocumentHandler {
	return &DocumentHandler{db: db}
}

// ── Column list & scan helper ────────────────────────────────────
// Keeps every query in sync with the Document struct.
const docCols = `d.id, d.employee_id, d.document_type,
	COALESCE(d.expiry_date::text, ''),
	d.is_primary,
	d.file_url, d.file_name, d.file_size, d.file_type,
	d.last_updated, d.created_at`

func scanDocument(scanner interface {
	Scan(dest ...interface{}) error
}, doc *models.Document) error {
	var expiryRaw string
	err := scanner.Scan(
		&doc.ID, &doc.EmployeeID, &doc.DocumentType,
		&expiryRaw,
		&doc.IsPrimary,
		&doc.FileURL, &doc.FileName, &doc.FileSize, &doc.FileType,
		&doc.LastUpdated, &doc.CreatedAt,
	)
	if err != nil {
		return err
	}
	if expiryRaw != "" {
		doc.ExpiryDate = &expiryRaw
	}
	return nil
}

// Create handles POST /api/employees/{employeeId}/documents
func (h *DocumentHandler) Create(w http.ResponseWriter, r *http.Request) {
	employeeID := chi.URLParam(r, "employeeId")
	if employeeID == "" {
		JSONError(w, http.StatusBadRequest, "Employee ID is required")
		return
	}

	var req models.CreateDocumentRequest
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	// Verify employee exists
	var exists bool
	if err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM employees WHERE id = $1)", employeeID).Scan(&exists); err != nil || !exists {
		JSONError(w, http.StatusNotFound, "Employee not found")
		return
	}

	var doc models.Document
	err := pool.QueryRow(ctx, fmt.Sprintf(`
		INSERT INTO documents (employee_id, document_type, expiry_date, file_url, file_name, file_size, file_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING %s
	`, docCols),
		employeeID, req.DocumentType, req.ExpiryDate, req.FileURL,
		req.FileName, req.FileSize, req.FileType,
	)
	if err2 := scanDocument(err, &doc); err2 != nil {
		log.Printf("Error creating document: %v", err2)
		JSONError(w, http.StatusInternalServerError, "Failed to create document")
		return
	}

	// Audit trail
	userID, _ := r.Context().Value(ctxkeys.UserID).(string)
	logActivity(pool, userID, "created", "document", doc.ID, map[string]interface{}{
		"type": doc.DocumentType, "employeeId": employeeID,
	})

	JSON(w, http.StatusCreated, map[string]interface{}{
		"data":    doc,
		"message": "Document created successfully",
	})
}

// ListByEmployee handles GET /api/employees/{id}/documents
func (h *DocumentHandler) ListByEmployee(w http.ResponseWriter, r *http.Request) {
	// Support both {id} (when nested under /employees/{id}) and {employeeId} (legacy)
	employeeID := chi.URLParam(r, "id")
	if employeeID == "" {
		employeeID = chi.URLParam(r, "employeeId")
	}
	if employeeID == "" {
		JSONError(w, http.StatusBadRequest, "Employee ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	rows, err := pool.Query(ctx, fmt.Sprintf(`
		SELECT %s
		FROM documents d
		WHERE d.employee_id = $1
		ORDER BY d.is_primary DESC, d.expiry_date ASC NULLS LAST
	`, docCols), employeeID)
	if err != nil {
		log.Printf("Error fetching documents: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch documents")
		return
	}
	defer rows.Close()

	documents := []models.Document{}
	for rows.Next() {
		var doc models.Document
		if err := scanDocument(rows, &doc); err != nil {
			log.Printf("Error scanning document: %v", err)
			continue
		}
		documents = append(documents, doc)
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": documents,
	})
}

// GetByID handles GET /api/documents/{id}
func (h *DocumentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	row := pool.QueryRow(ctx, `
		SELECT d.id, d.employee_id, d.document_type,
			COALESCE(d.expiry_date::text, ''),
			d.is_primary,
			d.file_url, d.file_name, d.file_size, d.file_type,
			d.last_updated, d.created_at,
			e.name AS employee_name, c.name AS company_name
		FROM documents d
		JOIN employees e ON d.employee_id = e.id
		JOIN companies c ON e.company_id = c.id
		WHERE d.id = $1
	`, id)

	var doc models.DocumentWithEmployee
	var expiryRaw string
	err := row.Scan(
		&doc.ID, &doc.EmployeeID, &doc.DocumentType,
		&expiryRaw,
		&doc.IsPrimary,
		&doc.FileURL, &doc.FileName, &doc.FileSize, &doc.FileType,
		&doc.LastUpdated, &doc.CreatedAt,
		&doc.EmployeeName, &doc.CompanyName,
	)
	if err != nil {
		log.Printf("Error fetching document %s: %v", id, err)
		JSONError(w, http.StatusNotFound, "Document not found")
		return
	}
	if expiryRaw != "" {
		doc.ExpiryDate = &expiryRaw
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": doc,
	})
}

// Update handles PUT /api/documents/{id}
func (h *DocumentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	var req models.UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	// Build dynamic SET clause
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.DocumentType != nil {
		setClauses = append(setClauses, fmt.Sprintf("document_type = $%d", argIdx))
		args = append(args, *req.DocumentType)
		argIdx++
	}
	if req.ExpiryDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("expiry_date = $%d", argIdx))
		args = append(args, *req.ExpiryDate)
		argIdx++
	}
	if req.FileURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("file_url = $%d", argIdx))
		args = append(args, *req.FileURL)
		argIdx++
	}
	if req.FileName != nil {
		setClauses = append(setClauses, fmt.Sprintf("file_name = $%d", argIdx))
		args = append(args, *req.FileName)
		argIdx++
	}
	if req.FileSize != nil {
		setClauses = append(setClauses, fmt.Sprintf("file_size = $%d", argIdx))
		args = append(args, *req.FileSize)
		argIdx++
	}
	if req.FileType != nil {
		setClauses = append(setClauses, fmt.Sprintf("file_type = $%d", argIdx))
		args = append(args, *req.FileType)
		argIdx++
	}

	if len(setClauses) == 0 {
		JSONError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	setClauses = append(setClauses, "last_updated = NOW()")

	setStr := ""
	for i, clause := range setClauses {
		if i > 0 {
			setStr += ", "
		}
		setStr += clause
	}

	query := fmt.Sprintf(`
		UPDATE documents d SET %s
		WHERE d.id = $%d
		RETURNING %s
	`, setStr, argIdx, docCols)
	args = append(args, id)

	var doc models.Document
	if err := scanDocument(pool.QueryRow(ctx, query, args...), &doc); err != nil {
		log.Printf("Error updating document %s: %v", id, err)
		JSONError(w, http.StatusNotFound, "Document not found")
		return
	}

	// Audit trail
	userID, _ := r.Context().Value(ctxkeys.UserID).(string)
	logActivity(pool, userID, "updated", "document", doc.ID, map[string]interface{}{
		"type": doc.DocumentType,
	})

	JSON(w, http.StatusOK, map[string]interface{}{
		"data":    doc,
		"message": "Document updated successfully",
	})
}

// TogglePrimary handles PATCH /api/documents/{id}/primary
// Sets this document as the primary tracked document for the employee.
// Only allowed if the document has an expiry date set.
// Uses a transaction to atomically unset the previous primary and set the new one.
func (h *DocumentHandler) TogglePrimary(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	tx, err := pool.Begin(ctx)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback(ctx)

	// 1. Fetch the target document — check it exists and has an expiry date
	var employeeID string
	var expiryExists bool
	var currentlyPrimary bool
	err = tx.QueryRow(ctx, `
		SELECT employee_id, (expiry_date IS NOT NULL), is_primary
		FROM documents WHERE id = $1
	`, id).Scan(&employeeID, &expiryExists, &currentlyPrimary)
	if err != nil {
		JSONError(w, http.StatusNotFound, "Document not found")
		return
	}

	if !expiryExists {
		JSONError(w, http.StatusBadRequest, "Cannot set as primary: document has no expiry date")
		return
	}

	if currentlyPrimary {
		// Un-primary (toggle off)
		_, err = tx.Exec(ctx, `UPDATE documents SET is_primary = FALSE WHERE id = $1`, id)
	} else {
		// Unset any existing primary for this employee, then set the new one
		_, err = tx.Exec(ctx, `UPDATE documents SET is_primary = FALSE WHERE employee_id = $1 AND is_primary = TRUE`, employeeID)
		if err == nil {
			_, err = tx.Exec(ctx, `UPDATE documents SET is_primary = TRUE WHERE id = $1`, id)
		}
	}

	if err != nil {
		log.Printf("Error toggling primary document %s: %v", id, err)
		JSONError(w, http.StatusInternalServerError, "Failed to toggle primary")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to commit")
		return
	}

	// Audit trail
	userID, _ := r.Context().Value(ctxkeys.UserID).(string)
	action := "set_primary"
	if currentlyPrimary {
		action = "unset_primary"
	}
	logActivity(pool, userID, action, "document", id, map[string]interface{}{
		"employeeId": employeeID,
	})

	JSON(w, http.StatusOK, map[string]string{
		"message": "Primary document updated successfully",
	})
}

// Delete handles DELETE /api/documents/{id}
func (h *DocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	tag, err := pool.Exec(ctx, "DELETE FROM documents WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting document %s: %v", id, err)
		JSONError(w, http.StatusInternalServerError, "Failed to delete document")
		return
	}

	if tag.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Document not found")
		return
	}

	// Audit trail
	userID, _ := r.Context().Value(ctxkeys.UserID).(string)
	logActivity(pool, userID, "deleted", "document", id, nil)

	JSON(w, http.StatusOK, map[string]string{
		"message": "Document deleted successfully",
	})
}

// ── Renew ──────────────────────────────────────────────────────

// Renew handles POST /api/documents/{id}/renew
// Creates a new document with the same type and employee, sets it as primary,
// and archives the old one by unsetting its primary flag.
// The request body should contain the new expiry date and optionally new file info.
func (h *DocumentHandler) Renew(w http.ResponseWriter, r *http.Request) {
	oldID := chi.URLParam(r, "id")
	if oldID == "" {
		JSONError(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	var req struct {
		ExpiryDate string `json:"expiryDate"`
		FileURL    string `json:"fileUrl,omitempty"`
		FileName   string `json:"fileName,omitempty"`
		FileSize   int64  `json:"fileSize,omitempty"`
		FileType   string `json:"fileType,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if req.ExpiryDate == "" {
		JSONError(w, http.StatusUnprocessableEntity, "New expiry date is required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	// Fetch the old document to copy its type, employee, and file info
	var employeeID, docType, fileURL, fileName, fileType string
	var fileSize int64
	err := pool.QueryRow(ctx, `
		SELECT employee_id, document_type, file_url, file_name, file_size, file_type
		FROM documents WHERE id = $1
	`, oldID).Scan(&employeeID, &docType, &fileURL, &fileName, &fileSize, &fileType)
	if err != nil {
		JSONError(w, http.StatusNotFound, "Original document not found")
		return
	}

	// Use new file info if provided, otherwise keep the old file
	if req.FileURL != "" {
		fileURL = req.FileURL
	}
	if req.FileName != "" {
		fileName = req.FileName
	}
	if req.FileSize > 0 {
		fileSize = req.FileSize
	}
	if req.FileType != "" {
		fileType = req.FileType
	}

	// Transaction: unset old primary → insert new doc as primary
	tx, err := pool.Begin(ctx)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback(ctx)

	// Unset old doc's primary flag
	_, err = tx.Exec(ctx, `UPDATE documents SET is_primary = FALSE WHERE employee_id = $1 AND is_primary = TRUE`, employeeID)
	if err != nil {
		log.Printf("Error unsetting primary for employee %s: %v", employeeID, err)
		JSONError(w, http.StatusInternalServerError, "Failed to archive old document")
		return
	}

	// Insert renewed document as the new primary
	var newDoc models.Document
	row := tx.QueryRow(ctx, fmt.Sprintf(`
		INSERT INTO documents (employee_id, document_type, expiry_date, is_primary, file_url, file_name, file_size, file_type)
		VALUES ($1, $2, $3, TRUE, $4, $5, $6, $7)
		RETURNING %s
	`, docCols), employeeID, docType, req.ExpiryDate, fileURL, fileName, fileSize, fileType)

	if err := scanDocument(row, &newDoc); err != nil {
		log.Printf("Error inserting renewed document: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to create renewed document")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to commit renewal")
		return
	}

	// Audit trail
	userID, _ := r.Context().Value(ctxkeys.UserID).(string)
	logActivity(pool, userID, "renewed", "document", newDoc.ID, map[string]interface{}{
		"previousDocId": oldID, "type": docType, "newExpiry": req.ExpiryDate,
	})

	JSON(w, http.StatusCreated, map[string]interface{}{
		"data":    newDoc,
		"message": "Document renewed successfully",
	})
}
