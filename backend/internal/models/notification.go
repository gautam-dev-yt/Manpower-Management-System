package models

// Notification represents an in-app notification for a user.
type Notification struct {
	ID         string  `json:"id"`
	UserID     string  `json:"userId"`
	Title      string  `json:"title"`
	Message    string  `json:"message"`
	Type       string  `json:"type"` // document_expiry, salary_due, system
	Read       bool    `json:"read"`
	EntityType *string `json:"entityType,omitempty"` // employee, document, salary
	EntityID   *string `json:"entityId,omitempty"`
	CreatedAt  string  `json:"createdAt"`
}
