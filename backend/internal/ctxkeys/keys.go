// Package ctxkeys defines typed context keys shared between middleware and handlers.
// This avoids import cycles: both middleware and handlers import this package,
// but neither imports the other for context key types.
package ctxkeys

// Key is a typed string used as context key to prevent collisions.
type Key string

const (
	UserID   Key = "userID"
	UserRole Key = "userRole"
)
