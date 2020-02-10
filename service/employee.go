package service

// EmployeeOutput is used as response for employee
type EmployeeOutput struct {
	CompanyID   int    `json:"company_id"`
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	RoleID      int    `json:"role_id"`
	RoleName    string `json:"role_name"`
	Status      string `json:"status"`
}
