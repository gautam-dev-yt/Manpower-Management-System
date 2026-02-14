package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"manpower-backend/internal/storage"
)

// Allowed file types and size limit for uploads.
const maxUploadSize = 10 << 20 // 10 MB

var allowedTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/jpg":       true,
	"image/png":       true,
}

// UploadHandler handles file upload requests.
// It depends on the storage.Store interface, not a specific implementation.
type UploadHandler struct {
	store storage.Store
}

// NewUploadHandler creates an UploadHandler with the given storage backend.
func NewUploadHandler(store storage.Store) *UploadHandler {
	return &UploadHandler{store: store}
}

// Upload handles multipart file uploads.
// Accepts: POST with multipart/form-data containing a "file" field.
// Returns: file metadata (url, name, size, type) as JSON.
func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Enforce size limit before reading body
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		JSONError(w, http.StatusBadRequest, "File too large. Maximum size is 10MB.")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Missing 'file' field in form data.")
		return
	}
	defer file.Close()

	// Validate file type by reading the first 512 bytes (MIME sniffing)
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		JSONError(w, http.StatusBadRequest, "Could not read file.")
		return
	}
	contentType := http.DetectContentType(buffer[:n])

	if !allowedTypes[contentType] {
		JSONError(w, http.StatusBadRequest, fmt.Sprintf(
			"File type '%s' not allowed. Accepted: PDF, JPG, PNG.", contentType,
		))
		return
	}

	// Reset file reader to beginning after MIME sniffing
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to process file.")
		return
	}

	// Build storage path: category/timestamp_filename
	// Optional "category" param allows organizing (e.g., "documents", "photos")
	category := r.FormValue("category")
	if category == "" {
		category = "general"
	}

	// Sanitize filename and add timestamp to prevent collisions
	safeName := sanitizeFilename(header.Filename)
	storagePath := fmt.Sprintf("%s/%d_%s", category, time.Now().Unix(), safeName)

	// Persist via the storage interface
	info, err := h.store.Save(r.Context(), storagePath, file, contentType)
	if err != nil {
		log.Printf("Upload failed: %v", err)
		JSONError(w, http.StatusInternalServerError, "Failed to save file.")
		return
	}

	JSON(w, http.StatusOK, info)
}

// ServeFile serves uploaded files from local storage.
// This endpoint is only needed for local storage — S3 uses pre-signed URLs instead.
func (h *UploadHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	// Extract the full file path from the URL (everything after /api/files/)
	filePath := strings.TrimPrefix(r.URL.Path, "/api/files/")
	if filePath == "" {
		JSONError(w, http.StatusBadRequest, "File path required.")
		return
	}

	// Serve using Go's built-in file server (handles Content-Type, caching, etc.)
	http.ServeFile(w, r, filepath.Join("uploads", filepath.Clean(filePath)))
}

// sanitizeFilename removes path separators and unsafe characters.
func sanitizeFilename(name string) string {
	// Keep only the base name (no directory components)
	name = filepath.Base(name)
	// Replace spaces with underscores for URL safety
	name = strings.ReplaceAll(name, " ", "_")
	return name
}
