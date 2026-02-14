package models

// ActivityLog records who changed what and when.
// The Details field uses JSONB to flexibly store change specifics
// (e.g., which fields were modified, old vs new values).
type ActivityLog struct {
	ID         string      `json:"id"`
	UserID     string      `json:"userId"`
	UserName   string      `json:"userName,omitempty"`
	Action     string      `json:"action"`     // "create", "update", "delete"
	EntityType string      `json:"entityType"` // "employee", "document", "company"
	EntityID   string      `json:"entityId"`
	Details    interface{} `json:"details,omitempty"` // Flexible JSON  data
	CreatedAt  string      `json:"createdAt"`
}
