package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"manpower-backend/internal/database"
	"manpower-backend/internal/models"
)

// DashboardHandler handles dashboard-related HTTP requests.
type DashboardHandler struct {
	db database.Service
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(db database.Service) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetMetrics handles GET /api/dashboard/metrics
func (h *DashboardHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	metrics := models.DashboardMetrics{}

	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM employees").Scan(&metrics.TotalEmployees)
	if err != nil {
		log.Printf("Error querying total employees: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch metrics")
		return
	}

	err = pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM documents 
		WHERE is_primary = TRUE AND expiry_date IS NOT NULL
		  AND expiry_date > CURRENT_DATE + INTERVAL '30 days'
	`).Scan(&metrics.ActiveDocuments)
	if err != nil {
		log.Printf("Error querying active documents: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch metrics")
		return
	}

	err = pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM documents 
		WHERE is_primary = TRUE AND expiry_date IS NOT NULL
		  AND expiry_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '30 days'
	`).Scan(&metrics.ExpiringSoon)
	if err != nil {
		log.Printf("Error querying expiring soon: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch metrics")
		return
	}

	err = pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM documents 
		WHERE is_primary = TRUE AND expiry_date IS NOT NULL
		  AND expiry_date < CURRENT_DATE
	`).Scan(&metrics.Expired)
	if err != nil {
		log.Printf("Error querying expired: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch metrics")
		return
	}

	JSON(w, http.StatusOK, metrics)
}

// GetExpiryAlerts handles GET /api/dashboard/expiring
func (h *DashboardHandler) GetExpiryAlerts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	rows, err := pool.Query(ctx, `
		SELECT 
			d.id, e.id, e.name, c.name, d.document_type,
			d.expiry_date::text,
			(d.expiry_date - CURRENT_DATE) AS days_left,
			CASE
				WHEN d.expiry_date < CURRENT_DATE THEN 'expired'
				WHEN d.expiry_date <= CURRENT_DATE + INTERVAL '7 days' THEN 'urgent'
				ELSE 'warning'
			END AS status
		FROM documents d
		JOIN employees e ON d.employee_id = e.id
		JOIN companies c ON e.company_id = c.id
		WHERE d.is_primary = TRUE AND d.expiry_date IS NOT NULL
		  AND d.expiry_date <= CURRENT_DATE + INTERVAL '30 days'
		ORDER BY d.expiry_date ASC
	`)
	if err != nil {
		log.Printf("Error fetching expiry alerts: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch alerts")
		return
	}
	defer rows.Close()

	alerts := []models.ExpiryAlert{}
	for rows.Next() {
		var a models.ExpiryAlert
		if err := rows.Scan(
			&a.DocumentID, &a.EmployeeID, &a.EmployeeName,
			&a.CompanyName, &a.DocumentType, &a.ExpiryDate,
			&a.DaysLeft, &a.Status,
		); err != nil {
			log.Printf("Error scanning alert: %v", err)
			continue
		}
		alerts = append(alerts, a)
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data":  alerts,
		"total": len(alerts),
	})
}

// GetCompanySummary handles GET /api/dashboard/company-summary
func (h *DashboardHandler) GetCompanySummary(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	rows, err := pool.Query(ctx, `
		SELECT c.id, c.name, COUNT(e.id) AS employee_count
		FROM companies c
		LEFT JOIN employees e ON e.company_id = c.id
		GROUP BY c.id, c.name
		ORDER BY employee_count DESC
	`)
	if err != nil {
		log.Printf("Error fetching company summary: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch company summary")
		return
	}
	defer rows.Close()

	companies := []models.CompanySummary{}
	for rows.Next() {
		var cs models.CompanySummary
		if err := rows.Scan(&cs.ID, &cs.Name, &cs.EmployeeCount); err != nil {
			log.Printf("Error scanning company summary: %v", err)
			continue
		}
		companies = append(companies, cs)
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": companies,
	})
}
