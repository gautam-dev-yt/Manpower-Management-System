package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// logActivity records a user action in the activity_log table.
// This is called after every successful mutation (create, update, delete)
// to provide a full audit trail visible on the Activity page.
//
// Parameters:
//   - pool:       database connection pool
//   - userID:     ID of the user who performed the action (from JWT context)
//   - action:     short verb describing what happened ("created", "updated", "deleted", "renewed", "toggled_primary")
//   - entityType: the kind of entity affected ("employee", "document", "salary")
//   - entityID:   UUID of the affected entity
//   - details:    optional key-value pairs with additional context (nil is fine)
func logActivity(pool *pgxpool.Pool, userID, action, entityType, entityID string, details map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var detailsJSON []byte
	if details != nil {
		var err error
		detailsJSON, err = json.Marshal(details)
		if err != nil {
			log.Printf("audit: failed to marshal details: %v", err)
			detailsJSON = nil
		}
	}

	_, err := pool.Exec(ctx, `
		INSERT INTO activity_log (user_id, action, entity_type, entity_id, details)
		VALUES ($1::uuid, $2, $3, $4::uuid, $5::jsonb)
	`, nilIfEmptyStr(userID), action, entityType, entityID, nilIfEmptyStr(string(detailsJSON)))

	if err != nil {
		// Log but don't fail the request — audit is best-effort.
		log.Printf("audit: failed to write activity log: %v", err)
	}
}
