package models

import "time"

// Employee represents an employee record in the database.
type Employee struct {
	ID              string    `json:"id"`
	CompanyID       string    `json:"companyId"`
	Name            string    `json:"name"`
	Trade           string    `json:"trade"`
	Mobile          string    `json:"mobile"`
	JoiningDate     string    `json:"joiningDate"`
	PhotoURL        *string   `json:"photoUrl"`
	Gender          *string   `json:"gender,omitempty"`
	DateOfBirth     *string   `json:"dateOfBirth,omitempty"`
	Nationality     *string   `json:"nationality,omitempty"`
	PassportNumber  *string   `json:"passportNumber,omitempty"`
	NativeLocation  *string   `json:"nativeLocation,omitempty"`
	CurrentLocation *string   `json:"currentLocation,omitempty"`
	Salary          *float64  `json:"salary,omitempty"`
	Status          string    `json:"status"` // active, inactive, on_leave
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// EmployeeWithCompany includes the company name alongside employee data.
// ExpiryDaysLeft and PrimaryDocType are computed from the employee's primary document.
type EmployeeWithCompany struct {
	Employee
	CompanyName    string  `json:"companyName"`
	DocStatus      string  `json:"docStatus"`                // "valid", "expiring", "expired", "none"
	ExpiryDaysLeft *int    `json:"expiryDaysLeft,omitempty"` // days until primary doc expires (negative = overdue)
	PrimaryDocType *string `json:"primaryDocType,omitempty"` // e.g. "Work Permit"
}

// CreateEmployeeRequest holds the fields needed to create an employee.
type CreateEmployeeRequest struct {
	CompanyID       string   `json:"companyId"`
	Name            string   `json:"name"`
	Trade           string   `json:"trade"`
	Mobile          string   `json:"mobile"`
	JoiningDate     string   `json:"joiningDate"`
	PhotoURL        string   `json:"photoUrl,omitempty"`
	Gender          *string  `json:"gender,omitempty"`
	DateOfBirth     *string  `json:"dateOfBirth,omitempty"`
	Nationality     *string  `json:"nationality,omitempty"`
	PassportNumber  *string  `json:"passportNumber,omitempty"`
	NativeLocation  *string  `json:"nativeLocation,omitempty"`
	CurrentLocation *string  `json:"currentLocation,omitempty"`
	Salary          *float64 `json:"salary,omitempty"`
	Status          string   `json:"status,omitempty"`
}

// UpdateEmployeeRequest holds the fields that can be updated.
type UpdateEmployeeRequest struct {
	CompanyID       *string  `json:"companyId,omitempty"`
	Name            *string  `json:"name,omitempty"`
	Trade           *string  `json:"trade,omitempty"`
	Mobile          *string  `json:"mobile,omitempty"`
	JoiningDate     *string  `json:"joiningDate,omitempty"`
	PhotoURL        *string  `json:"photoUrl,omitempty"`
	Gender          *string  `json:"gender,omitempty"`
	DateOfBirth     *string  `json:"dateOfBirth,omitempty"`
	Nationality     *string  `json:"nationality,omitempty"`
	PassportNumber  *string  `json:"passportNumber,omitempty"`
	NativeLocation  *string  `json:"nativeLocation,omitempty"`
	CurrentLocation *string  `json:"currentLocation,omitempty"`
	Salary          *float64 `json:"salary,omitempty"`
	Status          *string  `json:"status,omitempty"`
}

// Validate checks if the create request contains valid data.
func (r *CreateEmployeeRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if len(r.Name) < 2 || len(r.Name) > 100 {
		errors["name"] = "Name must be between 2 and 100 characters"
	}
	if len(r.Trade) < 2 {
		errors["trade"] = "Trade is required (min 2 characters)"
	}
	if r.CompanyID == "" {
		errors["companyId"] = "Company is required"
	}
	if r.Mobile == "" {
		errors["mobile"] = "Mobile number is required"
	}
	if r.JoiningDate == "" {
		errors["joiningDate"] = "Joining date is required"
	}

	return errors
}
