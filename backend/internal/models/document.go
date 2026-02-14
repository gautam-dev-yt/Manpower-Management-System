package models

import "time"

// Document represents a document record in the database.
// ExpiryDate is nullable — not all documents need an expiry (e.g. contracts).
// IsPrimary marks the single tracked document per employee for visa/permit alerts.
type Document struct {
	ID           string    `json:"id"`
	EmployeeID   string    `json:"employeeId"`
	DocumentType string    `json:"documentType"`
	ExpiryDate   *string   `json:"expiryDate"` // nullable — nil means no expiry
	IsPrimary    bool      `json:"isPrimary"`  // only one per employee
	FileURL      string    `json:"fileUrl"`
	FileName     string    `json:"fileName"`
	FileSize     int64     `json:"fileSize"`
	FileType     string    `json:"fileType"`
	LastUpdated  time.Time `json:"lastUpdated"`
	CreatedAt    time.Time `json:"createdAt"`
}

// DocumentWithEmployee includes the employee and company name.
type DocumentWithEmployee struct {
	Document
	EmployeeName string `json:"employeeName"`
	CompanyName  string `json:"companyName"`
}

// CreateDocumentRequest holds the fields needed to create a document.
// ExpiryDate is optional — documents without dates are supported.
type CreateDocumentRequest struct {
	DocumentType string  `json:"documentType"`
	ExpiryDate   *string `json:"expiryDate,omitempty"` // optional
	FileURL      string  `json:"fileUrl"`
	FileName     string  `json:"fileName"`
	FileSize     int64   `json:"fileSize"`
	FileType     string  `json:"fileType"`
}

// UpdateDocumentRequest holds the fields that can be updated.
type UpdateDocumentRequest struct {
	DocumentType *string `json:"documentType,omitempty"`
	ExpiryDate   *string `json:"expiryDate,omitempty"`
	FileURL      *string `json:"fileUrl,omitempty"`
	FileName     *string `json:"fileName,omitempty"`
	FileSize     *int64  `json:"fileSize,omitempty"`
	FileType     *string `json:"fileType,omitempty"`
	IsPrimary    *bool   `json:"isPrimary,omitempty"`
}

// Validate checks if the create request contains valid data.
// ExpiryDate is deliberately not required — many document types don't expire.
func (r *CreateDocumentRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if len(r.DocumentType) < 2 {
		errors["documentType"] = "Document type is required (min 2 characters)"
	}
	if r.FileURL == "" {
		errors["fileUrl"] = "File URL is required"
	}
	if r.FileName == "" {
		errors["fileName"] = "File name is required"
	}

	return errors
}
