package models

// DashboardMetrics holds the main dashboard statistics.
type DashboardMetrics struct {
	TotalEmployees  int `json:"totalEmployees"`
	ActiveDocuments int `json:"activeDocuments"`
	ExpiringSoon    int `json:"expiringSoon"`
	Expired         int `json:"expired"`
}

// Company represents a company record.
type Company struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Currency  string `json:"currency"` // e.g. "AED", "USD"
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CompanySummary includes employee count per company.
type CompanySummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Currency      string `json:"currency"`
	EmployeeCount int    `json:"employeeCount"`
}

// ExpiryAlert represents a document nearing/past expiry.
type ExpiryAlert struct {
	DocumentID   string `json:"documentId"`
	EmployeeID   string `json:"employeeId"`
	EmployeeName string `json:"employeeName"`
	CompanyName  string `json:"companyName"`
	DocumentType string `json:"documentType"`
	ExpiryDate   string `json:"expiryDate"`
	DaysLeft     int    `json:"daysLeft"`
	Status       string `json:"status"` // "expired", "urgent", "warning"
}
