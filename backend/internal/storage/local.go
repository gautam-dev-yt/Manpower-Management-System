package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStore saves files to the local filesystem.
// Implements the Store interface for development and single-server deployments.
type LocalStore struct {
	basePath string // Root directory for uploads (e.g., "./uploads")
	baseURL  string // URL prefix for serving files (e.g., "http://localhost:8080/api/files")
}

// NewLocalStore creates a LocalStore and ensures the upload directory exists.
func NewLocalStore(basePath, baseURL string) (*LocalStore, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("create upload directory: %w", err)
	}

	// Ensure base URL doesn't have a trailing slash
	baseURL = strings.TrimRight(baseURL, "/")

	return &LocalStore{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

// Save writes the file to disk under basePath/path.
// Creates subdirectories automatically (e.g., basePath/documents/employee-id/).
func (s *LocalStore) Save(ctx context.Context, path string, file io.Reader, contentType string) (*FileInfo, error) {
	fullPath := filepath.Join(s.basePath, filepath.Clean(path))

	// Create parent directories (e.g., uploads/documents/abc-123/)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, fmt.Errorf("create subdirectory: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		// Clean up partial file on failure
		os.Remove(fullPath)
		return nil, fmt.Errorf("write file: %w", err)
	}

	return &FileInfo{
		URL:      s.URL(path),
		FileName: filepath.Base(path),
		FileSize: written,
		FileType: contentType,
	}, nil
}

// Delete removes a file from disk. Returns nil if file doesn't exist.
func (s *LocalStore) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, filepath.Clean(path))

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}

// URL returns the HTTP URL to access the file via the API.
func (s *LocalStore) URL(path string) string {
	// Convert backslashes to forward slashes for URL compatibility
	return s.baseURL + "/" + strings.ReplaceAll(filepath.Clean(path), "\\", "/")
}
