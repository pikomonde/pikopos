package auth

import (
	"github.com/pikomonde/pikopos/repository"
)

// ServiceAuth contains repositories and Auth use cases
type ServiceAuth struct {
	RepositoryCompany          repository.RepositoryCompany
	RepositoryEmployee         repository.RepositoryEmployee
	RepositoryEmployeeRegister repository.RepositoryEmployeeRegister
	RepositoryRole             repository.RepositoryRole
}
