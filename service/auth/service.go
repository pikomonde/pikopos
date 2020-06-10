package auth

import (
	"github.com/pikomonde/pikopos/repository"
)

// ServiceAuth contains repositories and Auth use cases
type ServiceAuth struct {
	RepositoryAuth     repository.RepositoryAuth
	RepositoryCompany  repository.RepositoryCompany
	RepositoryEmployee repository.RepositoryEmployee
	RepositoryRole     repository.RepositoryRole
}
