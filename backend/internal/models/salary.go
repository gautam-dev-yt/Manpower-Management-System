package models

// SalaryRecord represents a single month's salary entry for an employee.
type SalaryRecord struct {
	ID         string  `json:"id"`
	EmployeeID string  `json:"employeeId"`
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"` // pending, paid, partial
	PaidDate   *string `json:"paidDate,omitempty"`
	Notes      *string `json:"notes,omitempty"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

// SalaryRecordWithEmployee extends SalaryRecord with display fields.
type SalaryRecordWithEmployee struct {
	SalaryRecord
	EmployeeName string `json:"employeeName"`
	CompanyName  string `json:"companyName"`
	Currency     string `json:"currency"`
}

// SalarySummary provides aggregated salary data for a given month/year.
type SalarySummary struct {
	TotalAmount  float64 `json:"totalAmount"`
	PaidAmount   float64 `json:"paidAmount"`
	PendingCount int     `json:"pendingCount"`
	PaidCount    int     `json:"paidCount"`
	PartialCount int     `json:"partialCount"`
	TotalCount   int     `json:"totalCount"`
	Currency     string  `json:"currency"` // "AED", "USD", or "Mixed"
}

// GenerateSalaryRequest triggers salary record creation for a month.
type GenerateSalaryRequest struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

// UpdateSalaryStatusRequest is for quick status toggle.
type UpdateSalaryStatusRequest struct {
	Status string `json:"status"`
}

// BulkUpdateSalaryRequest is for marking multiple records at once.
type BulkUpdateSalaryRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}
