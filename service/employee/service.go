package employee

import (
	"github.com/pikomonde/pikopos/repository"
)

// ServiceEmployee contains repositories and Employee use cases
type ServiceEmployee struct {
	repositoryEmployee repository.RepositoryEmployee
	repositoryRole     repository.RepositoryRole
}

// New returns the ServiceEmployee service
func New(
	rEmployee repository.RepositoryEmployee,
	rRole repository.RepositoryRole,
) *ServiceEmployee {
	return &ServiceEmployee{
		repositoryEmployee: rEmployee,
		repositoryRole:     rRole,
	}
}
