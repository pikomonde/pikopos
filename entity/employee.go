package entity

import (
	"time"
)

// EmployeeStatus is used to indicate whether a role unvalidate, active,
// or inactive
type EmployeeStatus int

const (
	// EmployeeStatusUnverified is registered, but not yet verified
	EmployeeStatusUnverified EmployeeStatus = iota
	// EmployeeStatusActive is active
	EmployeeStatusActive
	// EmployeeStatusInactive is inactive
	EmployeeStatusInactive
)

func (es EmployeeStatus) String() string {
	return [...]string{"unverified", "active", "inactive"}[es]
}

// Employee is used to contains employee data
type Employee struct {
	CompanyID   int
	ID          int
	FullName    string `max:"32"`
	Email       string `max:"48"`
	PhoneNumber string `max:"16"`
	RoleID      int
	Status      EmployeeStatus
	// Password    string  `max:"64"`
	// PIN         *string `max:"64"`
}

// EmployeeRegister is used to contains employee registration data
type EmployeeRegister struct {
	ID          int
	EmployeeID  int
	Email       string `max:"48"`
	PhoneNumber string `max:"16"`
	OTPCode     string `max:"64"`
	ExpiredAt   time.Time
}
