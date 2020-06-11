package auth

import (
	"github.com/pikomonde/pikopos/repository"
)

// ServiceAuth contains repositories and Auth use cases
type ServiceAuth struct {
	repositoryAuth     repository.RepositoryAuth
	repositoryCompany  repository.RepositoryCompany
	repositoryEmployee repository.RepositoryEmployee
	repositoryRole     repository.RepositoryRole
}

// New returns the ServiceAuth service
func New(
	rAuth repository.RepositoryAuth,
	rCompany repository.RepositoryCompany,
	rEmployee repository.RepositoryEmployee,
	rRole repository.RepositoryRole,
) *ServiceAuth {
	return &ServiceAuth{
		repositoryAuth:     rAuth,
		repositoryCompany:  rCompany,
		repositoryEmployee: rEmployee,
		repositoryRole:     rRole,
	}
}
