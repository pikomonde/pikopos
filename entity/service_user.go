package entity

type ServiceUserSession struct {
	CompanyID       int
	CompanyUsername string
	CompanyName     string
	ID              int
	FullName        string
	Email           string
	PhoneNumber     string
	Privileges      []string
}
