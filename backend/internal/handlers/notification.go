package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"manpower-backend/internal/ctxkeys"
	"manpower-backend/internal/database"
	"manpower-backend/internal/models"
)

// NotificationHandler handles notification HTTP requests.
type NotificationHandler struct {
	db database.Service
}

// NewNotificationHandler creates a new NotificationHandler.
func NewNotificationHandler(db database.Service) *NotificationHandler {
	return &NotificationHandler{db: db}
}

// List handles GET /api/notifications
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(ctxkeys.UserID).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	rows, err := pool.Query(ctx, `
		SELECT id, user_id, title, message, type, read,
			entity_type, entity_id, created_at::text
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 50
	`, userID)
	if err != nil {
		log.Printf("Error fetching notifications: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}
	defer rows.Close()

	notifications := []models.Notification{}
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(
			&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.Read,
			&n.EntityType, &n.EntityID, &n.CreatedAt,
		); err != nil {
			log.Printf("Error scanning notification: %v", err)
			continue
		}
		notifications = append(notifications, n)
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"data": notifications,
	})
}

// UnreadCount handles GET /api/notifications/count
func (h *NotificationHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(ctxkeys.UserID).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	var count int
	err := pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM notifications
		WHERE user_id = $1 AND read = false
	`, userID).Scan(&count)
	if err != nil {
		log.Printf("Error counting notifications: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to count notifications")
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"count": count,
	})
}

// MarkRead handles PATCH /api/notifications/{id}/read
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := r.Context().Value(ctxkeys.UserID).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	tag, err := pool.Exec(ctx, `
		UPDATE notifications SET read = true
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		log.Printf("Error marking notification read: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	if tag.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Notification not found")
		return
	}

	JSON(w, http.StatusOK, map[string]string{"message": "Marked as read"})
}

// MarkAllRead handles PATCH /api/notifications/read-all
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(ctxkeys.UserID).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	_, err := pool.Exec(ctx, `
		UPDATE notifications SET read = true
		WHERE user_id = $1 AND read = false
	`, userID)
	if err != nil {
		log.Printf("Error marking all notifications read: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to mark all as read")
		return
	}

	JSON(w, http.StatusOK, map[string]string{"message": "All marked as read"})
}
