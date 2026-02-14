package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"manpower-backend/internal/database"
	"manpower-backend/internal/models"
)

// ActivityHandler handles activity log HTTP requests.
type ActivityHandler struct {
	db database.Service
}

// NewActivityHandler creates a new ActivityHandler.
func NewActivityHandler(db database.Service) *ActivityHandler {
	return &ActivityHandler{db: db}
}

// List handles GET /api/activity?limit=N
func (h *ActivityHandler) List(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 100 {
		limit = 20
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	rows, err := pool.Query(ctx, `
		SELECT a.id, a.user_id, COALESCE(u.name, 'System') AS user_name,
			a.action, a.entity_type, a.entity_id, a.details,
			a.created_at::text
		FROM activity_log a
		LEFT JOIN users u ON a.user_id::uuid = u.id
		ORDER BY a.created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		log.Printf("Error fetching activity log: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch activity log")
		return
	}
	defer rows.Close()

	activities := []models.ActivityLog{}
	for rows.Next() {
		var a models.ActivityLog
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.UserName,
			&a.Action, &a.EntityType, &a.EntityID, &a.Details,
			&a.CreatedAt,
		); err != nil {
			log.Printf("Error scanning activity log: %v", err)
			continue
		}
		activities = append(activities, a)
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data":  activities,
		"total": len(activities),
	})
}
