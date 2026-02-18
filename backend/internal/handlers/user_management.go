package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"manpower-backend/internal/ctxkeys"
	"manpower-backend/internal/database"
	"manpower-backend/internal/models"
)

// UserManagementHandler provides admin-only user listing, role changes, and deletion.
type UserManagementHandler struct {
	db database.Service
}

func NewUserManagementHandler(db database.Service) *UserManagementHandler {
	return &UserManagementHandler{db: db}
}

// List returns all users (admin-only).
func (h *UserManagementHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	rows, err := pool.Query(ctx, `
		SELECT id, email, name, role, created_at::text, updated_at::text
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			log.Printf("Failed to scan user row: %v", err)
			continue
		}
		users = append(users, u)
	}

	if users == nil {
		users = []models.User{}
	}

	JSON(w, http.StatusOK, map[string]interface{}{"data": users})
}

// UpdateRole changes a user's role (admin-only). Cannot change your own role.
func (h *UserManagementHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	targetID := chi.URLParam(r, "id")
	currentUserID, _ := r.Context().Value(ctxkeys.UserID).(string)

	if targetID == currentUserID {
		JSONError(w, http.StatusBadRequest, "Cannot change your own role")
		return
	}

	var req models.UpdateRoleRequest
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

	var user models.User
	err := pool.QueryRow(ctx, `
		UPDATE users SET role = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, email, name, role, created_at::text, updated_at::text
	`, req.Role, targetID).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		JSONError(w, http.StatusNotFound, "User not found")
		return
	}

	go logActivity(pool, currentUserID, "updated_role", "user", targetID, map[string]interface{}{
		"newRole": req.Role,
		"email":   user.Email,
	})

	JSON(w, http.StatusOK, map[string]interface{}{
		"data":    user,
		"message": "Role updated successfully",
	})
}

// Delete removes a user (admin-only). Cannot delete yourself.
func (h *UserManagementHandler) Delete(w http.ResponseWriter, r *http.Request) {
	targetID := chi.URLParam(r, "id")
	currentUserID, _ := r.Context().Value(ctxkeys.UserID).(string)

	if targetID == currentUserID {
		JSONError(w, http.StatusBadRequest, "Cannot delete your own account")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pool := h.db.GetPool()

	// Get user info for audit log before deleting
	var email string
	err := pool.QueryRow(ctx, `SELECT email FROM users WHERE id = $1`, targetID).Scan(&email)
	if err != nil {
		JSONError(w, http.StatusNotFound, "User not found")
		return
	}

	tag, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, targetID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	if tag.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "User not found")
		return
	}

	go logActivity(pool, currentUserID, "deleted", "user", targetID, map[string]interface{}{
		"email": email,
	})

	JSON(w, http.StatusOK, map[string]interface{}{"message": "User deleted successfully"})
}
