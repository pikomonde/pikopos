package entity

// RoleStatus is used to indicate whether a role active or not
type RoleStatus int

const (
	// RoleStatusActive is active
	RoleStatusActive RoleStatus = iota
	// RoleStatusInactive is inactive
	RoleStatusInactive
)

const (
	// RoleSuperAdmin is the highest role in the app
	RoleSuperAdmin = "Super Admin"
)

func (rs RoleStatus) String() string {
	return [...]string{"active", "inactive"}[rs]
}

// Role is used to contains role data
type Role struct {
	CompanyID int
	ID        int
	Name      string `max:"16"`
	Status    RoleStatus
}
