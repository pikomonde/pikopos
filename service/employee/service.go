package employee

import (
	"github.com/pikomonde/pikopos/repository"
)

// ServiceEmployee contains repositories and Employee use cases
type ServiceEmployee struct {
	RepositoryEmployee repository.RepositoryEmployee
	RepositoryRole     repository.RepositoryRole
}
