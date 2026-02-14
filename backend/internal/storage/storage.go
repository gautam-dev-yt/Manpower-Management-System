// Package storage defines a file storage abstraction layer.
//
// This package uses the Strategy pattern — handlers depend on the Store interface,
// not a concrete implementation. To switch from local filesystem to S3:
//  1. Implement the Store interface in a new file (e.g., s3.go)
//  2. Change one line in main.go to inject the new implementation
//  3. No handler code changes needed
package storage

import (
	"context"
	"io"
)

// FileInfo holds metadata returned after a successful upload.
type FileInfo struct {
	URL      string `json:"url"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	FileType string `json:"fileType"`
}

// Store defines the contract for file storage operations.
// All implementations must be safe for concurrent use.
type Store interface {
	// Save persists a file at the given path and returns its metadata.
	// The path should be relative (e.g., "documents/employee-id/visa.pdf").
	Save(ctx context.Context, path string, file io.Reader, contentType string) (*FileInfo, error)

	// Delete removes a file at the given path.
	// Returns nil if the file doesn't exist (idempotent).
	Delete(ctx context.Context, path string) error

	// URL returns the publicly accessible URL for a stored file.
	URL(path string) string
}
